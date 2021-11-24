package main

import (
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
}
