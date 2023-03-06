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

	fs := http.FileServer(http.Dir("./client/build"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handleSearch(searcher))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("Listening on port %s...", port)
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
		w.Header().Set("Access-Control-Allow-Origin", "*")

		queryLengthLimit := 100
		currentPageNum := 1
		resultsLimit := 2000
		resultsPerPageNum := resultsLimit


		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}

		if len(query) > queryLengthLimit {
			query = query[:queryLengthLimit]
		}	

		if pageNum := r.URL.Query().Get("pageNum"); pageNum != "" {
			convertedPageNum, err := strconv.Atoi(pageNum)

			if err == nil {
				currentPageNum = convertedPageNum
			} else {
				w.Write([]byte("pageNum type error"))
				return
			}
		} 

		if resultsPerPage := r.URL.Query().Get("resultsPerPage"); resultsPerPage != "" {
			convertedResultsPerPage, err := strconv.Atoi(resultsPerPage)
			if err == nil {
				resultsPerPageNum = convertedResultsPerPage
			} else {
				w.Write([]byte("resultsPerPage type error"))
				return
			}
		} 

		results := searcher.Search(query[0], currentPageNum, resultsPerPageNum)
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
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	var datString = string(dat)
	s.CompleteWorks = datString
	s.SuffixArray = suffixarray.New([]byte(strings.ToLower(datString)))
	return nil
}

func (s *Searcher) Search(query string, pageNum int, resultsPerPage int) []string {
	startIdx := (pageNum -1)*resultsPerPage
	endIdx := startIdx + resultsPerPage

	idxs := s.SuffixArray.Lookup([]byte(strings.ToLower(query)), -1)

	results := []string{}

	for i, idx := range idxs {
		if i>= startIdx && i < endIdx {
			textStart := idx - 250
			textEnd := idx + 250
			
			if textStart < 0 {
				textStart = 0
			}

			if textEnd > len(s.CompleteWorks) {
				textEnd = len(s.CompleteWorks)
			}

			results = append(results, s.CompleteWorks[textStart:textEnd])

			if len(results) == endIdx {
				break
			}
		}
	
	}
	return results
}
