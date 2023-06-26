package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		query          string
		expectedResult bool
	}{
		{
			name:           "Passing Test",
			query:          "Hamlet",
			expectedResult: true,
		},
		{
			name:           "Failing Test - Case Sensitive",
			query:          "juliet",
			expectedResult: true,
		},
		{
			name:           "Failing Test - Punctuation",
			query:          "to be or not to be",
			expectedResult: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/search?q="+test.query, nil)
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
				if strings.Contains(strings.ToLower(result), strings.ToLower(test.query)) {
					found = true
					break
				}
			}

			if found != test.expectedResult {
				t.Errorf("expected result for `%s`: got %v want %v", test.query, found, test.expectedResult)
			}
		})
	}
}
