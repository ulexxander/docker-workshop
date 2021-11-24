package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type client struct {
	mux *http.ServeMux
}

func (c *client) Request(t *testing.T, method, path string, reqBody io.Reader, resBody interface{}) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(method, path, reqBody)
	if err != nil {
		t.Fatalf("error initializing request: %s", err)
	}

	c.mux.ServeHTTP(rec, req)

	if err := json.NewDecoder(rec.Body).Decode(&resBody); err != nil {
		t.Fatalf("error decoding body: %s", err)
	}
	return rec
}

func (c *client) Get(t *testing.T, path string, reqBody io.Reader, resBody interface{}) *httptest.ResponseRecorder {
	return c.Request(t, "GET", path, reqBody, resBody)
}

func TestNotesEndpoints(t *testing.T) {
	mux := http.NewServeMux()
	endpoints := Endpoints{}
	endpoints.Register(mux)

	c := client{mux}

	t.Run("no notes yet", func(t *testing.T) {
		var notes []Note
		rec := c.Get(t, "/notes/all", nil, &notes)
		if rec.Result().StatusCode != 200 {
			t.Fatalf("expected status code 200, got: %d", rec.Result().StatusCode)
		}
		if len(notes) != 0 {
			t.Fatalf("expected 0 notes, got: %d", len(notes))
		}
	})
}
