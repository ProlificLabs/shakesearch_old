package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"index/suffixarray"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handleSearch(searcher))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("shakesearch available at http://localhost:%s...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

type Searcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}
		results := searcher.Search(query[0])
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
}

func (s *Searcher) Load(filename string) error {
    dat, err := os.ReadFile(filename)
    if err != nil {
        return fmt.Errorf("Load: %w", err)
    }
    lowerData := strings.ToLower(string(dat))
    s.CompleteWorks = lowerData
    s.SuffixArray = suffixarray.New([]byte(lowerData))
    return nil
}

func (s *Searcher) Search(query string) []string {
	query = strings.ToLower(query)
	idxs := s.SuffixArray.Lookup([]byte(query), -1)
	results := []string{}
	for _, idx := range idxs {
		start := idx - 250
		end := idx + 250

		if start < 0 {
			start = 0
		}
		if end > len(s.CompleteWorks) {
			end = len(s.CompleteWorks)
		}
		results = append(results, s.CompleteWorks[start:end])
	}
	return results
}
