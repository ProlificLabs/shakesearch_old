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
	"strings"
    "strconv"
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
		pageNum, ok := r.URL.Query()["p"]
		var finalPageNum int
		if(pageNum == nil) {
			finalPageNum = 0
		} else {
			convertedPageNum, err := strconv.Atoi(pageNum[0])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("invalid number for page in URL params"))
				return
			}
			finalPageNum = convertedPageNum
			
		}
		results := searcher.Search(strings.ToLower(query[0]), finalPageNum)
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
	s.SuffixArray = suffixarray.New([]byte(strings.ToLower(string(dat))))
	return nil
}

func (s *Searcher) Search(query string, pageNum int) []string {
	allIdxs := s.SuffixArray.Lookup([]byte(query), -1)
	
	if(len(allIdxs) <= (pageNum)* 20){
		return []string{}
	}
	var nextPage = (pageNum + 1) * 20
	if (nextPage > len(allIdxs)){
		nextPage = len(allIdxs)
	}
	filteredIdxs := allIdxs[(pageNum * 20):nextPage]
	

	results := []string{}
	for _, idx := range filteredIdxs {
		results = append(results, s.CompleteWorks[idx-250:idx+250])
	}
	return results
}
