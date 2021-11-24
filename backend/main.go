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

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello\n")
	})

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
