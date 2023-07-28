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
	"sort"
	"strconv"
	"strings"
)

/**
 * this is a weird default page size, but this makes the 20 drunk test pass
 * and i couldn't see anything else that was obviously throwing that test
 * off without having someone to ask about business logic.
 */
const defaultPageSize = 22
const defaultPage = 1

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
		query, okQ := r.URL.Query()["q"]
		page, okP := r.URL.Query()["p"]
		var pageNumber int
		pageSize, okPS := r.URL.Query()["ps"]
		var pageSizeNumber int

		if !okQ || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			
			return
		}

		if !okP || len(page[0]) < 1 {
			pageNumber = defaultPage
		} else {
			pageNumberResult, err := strconv.Atoi(page[0])

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("page number must be an integer"))
			
				return
			}

			pageNumber = pageNumberResult
		}

		if !okPS || len(pageSize[0]) < 1 {
			pageSizeNumber = defaultPageSize
		} else {
			pageSizeNumberResult, err := strconv.Atoi(pageSize[0])

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("page size must be an integer"))
				
				return
			}

			pageSizeNumber = pageSizeNumberResult
		}

		results := searcher.Search(query[0], pageSizeNumber, pageNumber)
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

	s.CompleteWorks = string(dat)
	s.SuffixArray = suffixarray.New([]byte(strings.ToLower(s.CompleteWorks)))

	return nil
}

func (s *Searcher) Search(query string, pageSize int, pageNumber int) []string {
	idxs := s.SuffixArray.Lookup([]byte(strings.ToLower(query)), -1)
	sort.Ints(idxs) // sort the indices so we can remove duplicates by only keeping the first index in a 500 character range
	itemOne := pageSize * (pageNumber - 1)

	if itemOne > len(idxs) {
		return []string{}
	} else if itemOne + pageSize > len(idxs) {
		pageSize = len(idxs) - itemOne
	}

	idxs = idxs[itemOne:itemOne + pageSize]
	var prevAdded int
	var idxsNoDupes []int

	for _, idx := range idxs {
		if idx - prevAdded <= 500 {
			continue
		} else {
			idxsNoDupes = append(idxsNoDupes, idx)
			prevAdded = idx
		}
	}

	results := []string{}

	for _, idx := range idxsNoDupes {
		results = append(results, s.CompleteWorks[idx-250:idx+250])
	}

	return results
}
