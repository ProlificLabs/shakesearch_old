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

type ConfigValues struct {
	Minimum int
	Maximum int
	Default int
}

type PaginationConfig struct {
	Offset ConfigValues
	Limit  ConfigValues
}

type Pagination struct {
	Offset int
	Limit  int
}

func (p *Pagination) isOffsetValid() bool {
	return p.Offset >= paginationConfig.Offset.Minimum
}

func (p *Pagination) isLimitValid() bool {
	return p.Limit >= paginationConfig.Limit.Minimum && p.Limit <= paginationConfig.Limit.Maximum
}

func (p *Pagination) validatePagination() {
	if !p.isOffsetValid() {
		p.Offset = paginationConfig.Offset.Default
	}
	if !p.isLimitValid() {
		p.Limit = paginationConfig.Limit.Default
	}
}

// TODO: load pagination values from config file
var paginationConfig = PaginationConfig{
	Offset: ConfigValues{
		Minimum: 0,
		Default: 0,
	},
	Limit: ConfigValues{
		Minimum: 1,
		Maximum: 20,
		Default: 20,
	},
}

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

func getPaginationParams(r *http.Request) (Pagination, error) {
	var err error
	var limit, offset int
	offsetParam, offsetOK := r.URL.Query()["o"]
	if offsetOK {
		offset, err = strconv.Atoi(offsetParam[0])
		if err != nil {
			return Pagination{}, fmt.Errorf("invalid offset number in URL params: %w", err)
		}
	}
	limitParam, limitOK := r.URL.Query()["l"]
	if limitOK {
		limit, err = strconv.Atoi(limitParam[0])
		if err != nil {
			return Pagination{}, fmt.Errorf("invalid limit number in URL params: %w", err)
		}
	}
	pagination := Pagination{
		Offset: offset,
		Limit:  limit,
	}
	pagination.validatePagination()
	return pagination, nil
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}
		pagination, err := getPaginationParams(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		results := searcher.SearchPaginated(query[0], pagination)
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
	lowerCaseDat := []byte(strings.ToLower(string(dat)))
	s.CompleteWorks = string(dat)
	s.SuffixArray = suffixarray.New(lowerCaseDat)
	return nil
}

func (s *Searcher) Search(query string) []string {
	pagination := Pagination{
		Offset: paginationConfig.Offset.Default,
		Limit:  paginationConfig.Limit.Default,
	}
	return s.SearchPaginated(query, pagination)
}

func (s *Searcher) SearchPaginated(query string, pagination Pagination) []string {
	if query != "" {
		query = strings.ToLower(query)
	}
	idxs := s.SuffixArray.Lookup([]byte(query), -1)
	results := []string{}
	if idxs == nil {
		return results
	}
	startIdx := pagination.Offset
	if startIdx > len(idxs) {
		return results
	}

	endIdx := pagination.Offset + pagination.Limit
	if endIdx > len(idxs) {
		endIdx = len(idxs)
	}
	idxRange := idxs[startIdx:endIdx]
	if idxRange == nil {
		return results
	}
	for _, idx := range idxRange {
		results = append(results, s.CompleteWorks[idx-250:idx+250])
	}
	return results
}
