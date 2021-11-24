package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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

type Endpoints struct{}

func (e *Endpoints) Register(m *http.ServeMux) {
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello\n")
	})
}
