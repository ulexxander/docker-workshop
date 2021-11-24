package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	log.Println("starting service")

	mux := http.NewServeMux()

	endpoints := Endpoints{
		Notes: &NotesStoreMemory{},
	}
	endpoints.Register(mux)

	addr, ok := os.LookupEnv("HTTP_SERVER_ADDR")
	if !ok {
		addr = ":80"
	}

	log.Println("listening http", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatalf("error listening http: %s", err)
	}
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
	CreateNote(p NoteCreateParams) (Note, error)
}

type NotesStoreMemory struct {
	notes []Note
}

func (ns *NotesStoreMemory) AllNotes() ([]Note, error) {
	return ns.notes, nil
}

func (ns *NotesStoreMemory) CreateNote(p NoteCreateParams) (Note, error) {
	note := Note{
		ID:        len(ns.notes),
		Text:      p.Text,
		CreatedAt: time.Now(),
	}
	ns.notes = append(ns.notes, note)
	return note, nil
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
