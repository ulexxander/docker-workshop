package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	log.Println("starting service")

	ctx := context.Background()
	mongoc, err := setupMongo(ctx)
	if err != nil {
		log.Fatalf("error setting up mongo: %s", err)
	}
	defer mongoc.Disconnect(ctx)

	notes := NewNotesStoreMongo(mongoc.Database("docker-workshop").Collection("notes"))

	mux := http.NewServeMux()
	endpoints := Endpoints{Notes: notes}
	endpoints.Register(mux)

	addr, ok := os.LookupEnv("HTTP_SERVER_ADDR")
	if !ok {
		addr = ":80"
	}

	log.Println("listening http", addr)
	err = http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatalf("error listening http: %s", err)
	}
}

func setupMongo(ctx context.Context) (*mongo.Client, error) {
	mongoURI, ok := os.LookupEnv("MONGO_URI")
	if !ok {
		mongoURI = "mongodb://docker-workshop:123123@localhost:27017"
	}

	log.Println("connecting to mongodb", mongoURI)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, errors.Wrap(err, "connecting to mongo")
	}

	return client, nil
}

type Note struct {
	ID        int
	Text      string
	CreatedAt time.Time
}

type NoteCreateParams struct {
	Text string
}

type NotesStore interface {
	AllNotes() ([]Note, error)
	CreateNote(p NoteCreateParams) (*Note, error)
}

type NotesStoreMemory struct {
	notes []Note
}

func NewNotesStoreMemory() *NotesStoreMemory {
	return &NotesStoreMemory{
		notes: []Note{},
	}
}

func (ns *NotesStoreMemory) AllNotes() ([]Note, error) {
	return ns.notes, nil
}

func (ns *NotesStoreMemory) CreateNote(p NoteCreateParams) (*Note, error) {
	note := Note{
		ID:        len(ns.notes),
		Text:      p.Text,
		CreatedAt: time.Now(),
	}
	ns.notes = append(ns.notes, note)
	return &note, nil
}

type NotesStoreMongo struct {
	count      int
	collection *mongo.Collection
}

func NewNotesStoreMongo(collection *mongo.Collection) *NotesStoreMongo {
	return &NotesStoreMongo{
		collection: collection,
	}
}

func (ns *NotesStoreMongo) AllNotes() ([]Note, error) {
	c, err := ns.collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "finding notes")
	}
	var notes []Note
	if err := c.All(context.TODO(), &notes); err != nil {
		return nil, errors.Wrap(err, "decoding notes")
	}
	if notes == nil {
		notes = []Note{}
	}
	return notes, nil
}

func (ns *NotesStoreMongo) CreateNote(p NoteCreateParams) (*Note, error) {
	note := Note{
		ID:        ns.count,
		Text:      p.Text,
		CreatedAt: time.Now(),
	}
	ns.count++
	if _, err := ns.collection.InsertOne(context.TODO(), note); err != nil {
		return nil, errors.Wrap(err, "inserting item")
	}
	return &note, nil
}

func (ns *NotesStoreMongo) Reset() error {
	_, err := ns.collection.DeleteMany(context.TODO(), bson.D{})
	return err
}

type Endpoints struct {
	Notes NotesStore
}

func (e *Endpoints) Register(m *http.ServeMux) {
	m.HandleFunc("/notes/all", func(w http.ResponseWriter, r *http.Request) {
		notes, err := e.Notes.AllNotes()
		if err != nil {
			responseError(w, err)
			return
		}
		responseData(w, notes)
	})

	m.HandleFunc("/notes/create", func(w http.ResponseWriter, r *http.Request) {
		var params NoteCreateParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			responseError(w, err)
			return
		}
		note, err := e.Notes.CreateNote(params)
		if err != nil {
			responseError(w, err)
			return
		}
		responseData(w, note)
	})
}

type APIResponse struct {
	Data interface{}
}

type APIError struct {
	Error string
}

func response(w http.ResponseWriter, res interface{}) {
	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func responseData(w http.ResponseWriter, data interface{}) {
	res := APIResponse{Data: data}
	response(w, res)
}

func responseError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	res := APIError{Error: err.Error()}
	response(w, res)
}
