package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type client struct {
	mux *http.ServeMux
}

func (c *client) Request(t *testing.T, method, path string, reqBody, resBody interface{}) *httptest.ResponseRecorder {
	var body io.Reader
	if reqBody != nil {
		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			t.Fatalf("error encoding body: %s", err)
		}
		body = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, path, body)
	if err != nil {
		t.Fatalf("error initializing request: %s", err)
	}

	rec := httptest.NewRecorder()
	c.mux.ServeHTTP(rec, req)

	if err := json.NewDecoder(rec.Body).Decode(&resBody); err != nil {
		t.Fatalf("error decoding body: %s", err)
	}
	return rec
}

func (c *client) Get(t *testing.T, path string, reqBody, resBody interface{}) *httptest.ResponseRecorder {
	return c.Request(t, http.MethodGet, path, reqBody, resBody)
}

func (c *client) Post(t *testing.T, path string, reqBody, resBody interface{}) *httptest.ResponseRecorder {
	return c.Request(t, http.MethodPost, path, reqBody, resBody)
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

	text := "my note"

	t.Run("adding note", func(t *testing.T) {
		var note Note
		rec := c.Get(t, "/notes/create", NoteCreateParams{
			Text: text,
		}, &note)
		if rec.Result().StatusCode != 200 {
			t.Fatalf("expected status code 200, got: %d", rec.Result().StatusCode)
		}
		if note.Text != text {
			t.Errorf("Text of created note must be %s, got: %s", text, note.Text)
		}
		if note.CreatedAt.IsZero() {
			t.Errorf("CreatedAt not be zero")
		}
	})
}
