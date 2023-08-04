package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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

		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			page = 1
		}

		pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
		if err != nil {
			pageSize = 20
		}

		results := searcher.Search(query[0], page, pageSize)
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err = enc.Encode(results)
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
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	lowerCaseDat := strings.ToLower(string(dat))
	s.CompleteWorks = string(dat)
	s.SuffixArray = suffixarray.New([]byte(lowerCaseDat))
	return nil
}

func (s *Searcher) Search(query string, page int, pageSize int) []string {
	lowerCaseQuery := strings.ToLower(query)
	idxs := s.SuffixArray.Lookup([]byte(lowerCaseQuery), -1)
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(idxs) {
		return []string{}
	}
	if end > len(idxs) {
		end = len(idxs)
	}
	results := []string{}
	for _, idx := range idxs[start:end] {
		results = append(results, s.CompleteWorks[idx-250:idx+250])
	}
	return results
}
