package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
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

const (
	storeMemory = "memory"
	storeMongo  = "mongo"
)

var flagStore = flag.String("store", storeMemory, "implementation of NotesStore - "+storeMemory+" or "+storeMongo)

func TestNotesEndpoints(t *testing.T) {
	ctx := context.Background()
	mongoc, err := setupMongo(ctx)
	if err != nil {
		t.Fatalf("error setting up mongo: %s", err)
	}
	defer mongoc.Disconnect(ctx)

	var notes NotesStore
	switch *flagStore {
	case storeMemory:
		notes = NewNotesStoreMemory()
	case storeMongo:
		store := NewNotesStoreMongo(mongoc.Database("docker-workshop-test").Collection("notes"))
		if err := store.Reset(); err != nil {
			t.Fatalf("error resetting notes: %s", err)
		}
		notes = store
	default:
		panic(fmt.Sprintf("unknown store: %s", *flagStore))
	}

	mux := http.NewServeMux()
	endpoints := Endpoints{Notes: notes}
	endpoints.Register(mux)

	c := client{mux}

	t.Run("no notes yet", func(t *testing.T) {
		var resBody struct {
			Data []Note
			APIError
		}
		rec := c.Get(t, "/notes/all", nil, &resBody)
		if rec.Result().StatusCode != 200 {
			t.Errorf("expected status code 200, got: %d", rec.Result().StatusCode)
		}
		if resBody.Error != "" {
			t.Fatalf("unexpected api error: %s", resBody.Error)
		}
		if len(resBody.Data) != 0 {
			t.Fatalf("expected 0 notes, got: %d", len(resBody.Data))
		}
		if resBody.Data == nil {
			t.Errorf("empty slice must not be nil")
		}
	})

	text := "my note"

	t.Run("adding note", func(t *testing.T) {
		var resBody struct {
			Data Note
			APIError
		}
		rec := c.Get(t, "/notes/create", NoteCreateParams{
			Text: text,
		}, &resBody)
		if rec.Result().StatusCode != 200 {
			t.Errorf("expected status code 200, got: %d", rec.Result().StatusCode)
		}
		if resBody.Error != "" {
			t.Fatalf("unexpected api error: %s", resBody.Error)
		}
		if resBody.Data.ID != 0 {
			t.Errorf("ID of created note must be 0, got: %d", resBody.Data.ID)
		}
		if resBody.Data.Text != text {
			t.Errorf("Text of created note must be %s, got: %s", text, resBody.Data.Text)
		}
		if resBody.Data.CreatedAt.IsZero() {
			t.Errorf("CreatedAt must not be zero")
		}
	})

	t.Run("created note listed", func(t *testing.T) {
		var resBody struct {
			Data []Note
			APIError
		}
		rec := c.Get(t, "/notes/all", nil, &resBody)
		if rec.Result().StatusCode != 200 {
			t.Errorf("expected status code 200, got: %d", rec.Result().StatusCode)
		}
		if resBody.Error != "" {
			t.Fatalf("unexpected api error: %s", resBody.Error)
		}
		if len(resBody.Data) != 1 {
			t.Fatalf("expected 1 notes, got: %d", len(resBody.Data))
		}
	})
}
