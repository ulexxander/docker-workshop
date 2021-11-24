package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEndpoints(t *testing.T) {
	mux := http.NewServeMux()
	endpoints := Endpoints{}
	endpoints.Register(mux)

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("error initializing request: %s", err)
	}

	mux.ServeHTTP(rec, req)

	resBody := rec.Body.String()
	expected := "hello\n"
	if resBody != expected {
		t.Fatalf("expected response %s, got: %s", expected, resBody)
	}
}
