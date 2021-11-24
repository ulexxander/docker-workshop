package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	log.Println("starting service")

	mux := http.NewServeMux()

	endpoints := Endpoints{}
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
	Text      string
	CreatedAt time.Time
}

type NoteCreateParams struct {
	Text string
}

type Endpoints struct{}

func (e *Endpoints) Register(m *http.ServeMux) {
	m.HandleFunc("/notes/all", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "[]")
	})

	m.HandleFunc("/notes/create", func(w http.ResponseWriter, r *http.Request) {
		var params NoteCreateParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			responseError(w, err)
			return
		}
		note := Note{
			Text:      params.Text,
			CreatedAt: time.Now(),
		}
		response(w, note)
	})
}

type APIError struct {
	Error string
}

func response(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func responseError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)

	res := APIError{Error: err.Error()}
	response(w, res)
}
