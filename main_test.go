package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSearchHamlet(t *testing.T) {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		t.Fatal(err)
	}

	query := "Hamlet"
	req, err := http.NewRequest("GET", "/search?q="+query, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleSearch(searcher))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var results []string
	err = json.Unmarshal(rr.Body.Bytes(), &results)
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, result := range results {
		if strings.Contains(strings.ToLower(result), strings.ToLower(query)) {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected result not found for query: %s", query)
	}
}

func TestSearchCaseSensitive(t *testing.T) {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		t.Fatal(err)
	}

	query := "hAmLeT"
	req, err := http.NewRequest("GET", "/search?q="+query, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleSearch(searcher))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var results []string
	err = json.Unmarshal(rr.Body.Bytes(), &results)
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, result := range results {
		if strings.Contains(strings.ToLower(result), strings.ToLower(query)) {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected result not found for query: %s", query)
	}
}

func TestSearchDrunk(t *testing.T) {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		t.Fatal(err)
	}

	query := "drunk"
	req, err := http.NewRequest("GET", "/search?q="+query, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleSearch(searcher))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var results []string
	err = json.Unmarshal(rr.Body.Bytes(), &results)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 20 {
		t.Errorf("expected 20 results for query: %s, got %d", query, len(results))
	}
}
