package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"pulley.com/shakesearch/pkg/searcher"
)

type mockSearcher struct {
	res searcher.Response
}

func (s *mockSearcher) Search(r searcher.Request) (*searcher.Response, error) {
	return &s.res, nil
}

//TODO add more test cases
func TestHandler_Search(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.URL.RawQuery = "q=test"

	rr := httptest.NewRecorder()

	s := mockSearcher{
		res: searcher.Response{
			Query: "test",
			Hits: []*searcher.Hit{
				{},
			},
			Highlights: []*searcher.Highlight{
				{},
			},
		},
	}

	wantResBuf := bytes.NewBuffer(nil)
	err = json.NewEncoder(wantResBuf).Encode(s.res)
	if err != nil {
		t.Fatal(err)
	}

	handler := NewHandler(&s)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Search)
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if !bytes.Equal(rr.Body.Bytes(), wantResBuf.Bytes()) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), string(wantResBuf.Bytes()))
	}
}
