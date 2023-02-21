package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strings"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	searcher := Searcher{
		logger: logger,
	}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handleSearch(searcher))

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

type Searcher struct {
	CompleteWorks              string
	LowerSuffixArray           *suffixarray.Index // making the index with lower case bytes
	StartIndexOfPlayToPlayName map[int]metadata   // keeps mapping of start index of play to play name
	AllSortedStartIndexes      []int
	logger                     *zap.Logger // keeps all start indexes of plays sorted
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		sanitizedQuery := sanitize(query[0]) // sanitize the query-- just trims the spaces for now
		l := searcher.logger
		l.Debug(fmt.Sprintf("query: %s, sanitized: %s\n", query[0], sanitizedQuery))
		if !ok || len(sanitizedQuery) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}
		results := searcher.Search(sanitizedQuery)
		l.Debug(fmt.Sprintf("results: %v", len(results)))
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
}

func sanitize(s string) string {
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	return s
}

func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.CompleteWorks = string(dat)
	s.LowerSuffixArray = suffixarray.New(bytes.ToLower(dat)) // making the index with lower case bytes
	startIndexOfPlayToPlayName, allSortedStartIndexes := s.getMetadata()
	s.StartIndexOfPlayToPlayName = startIndexOfPlayToPlayName
	s.AllSortedStartIndexes = allSortedStartIndexes
	return nil
}

type SearchResult struct { // this is the struct that will be returned as json
	Line        string `json:"line"`
	PlayName    string `json:"play_name"`
	ActNumber   int    `json:"act_number"`   // these are not computed and random numbers are returned for now
	SceneNumber int    `json:"scene_number"` // these are not computed and random numbers are returned for now
}

const charsCoverage = 250

func (s *Searcher) Search(query string) []SearchResult {
	queries := strings.Split(query, " ")
	queries = append(queries, query)
	finalIndexes := []int{}
	for _, query := range queries {
		idxs := s.LowerSuffixArray.Lookup(bytes.ToLower([]byte(query)), -1)
		finalIndexes = append(finalIndexes, idxs...)
	}
	var results []SearchResult
	results = make([]SearchResult, 0)

	for _, idx := range finalIndexes {
		playName := s.findPlayName(idx)
		line := s.CompleteWorks[idx-charsCoverage : idx+charsCoverage]
		results = append(results, SearchResult{Line: sanitize(line), PlayName: playName, ActNumber: rand.Int()%5 + 1, SceneNumber: rand.Int()%5 + 1})
	}
	return results
}

type metadata struct {
	PlayName string
	Start    int
}

func (s *Searcher) findPlayName(idx int) string {
	prevStart := 0
	for _, starts := range s.AllSortedStartIndexes {
		if starts > idx {
			return s.StartIndexOfPlayToPlayName[prevStart].PlayName
		}
		prevStart = starts
	}
	return s.StartIndexOfPlayToPlayName[prevStart].PlayName
}

const _firstPlayName = "the sonnets"

func (s *Searcher) getMetadata() (map[int]metadata, []int) {
	idx := s.LowerSuffixArray.Lookup([]byte(_firstPlayName), 2)
	allNamesAsString := s.CompleteWorks[idx[1]:idx[0]] // this is all the play names in the beginning
	rawNames := strings.Split(allNamesAsString, "\n")
	sanitizedPlayNames := []string{}
	for _, name := range rawNames {
		playName := strings.TrimSpace(name)
		if playName != "" {
			sanitizedPlayNames = append(sanitizedPlayNames, playName)
		}
	}
	var allStartsIndexes []int
	startIndexOfPlayMap := map[int]metadata{}
	for _, name := range sanitizedPlayNames {
		indexes := s.LowerSuffixArray.Lookup(bytes.ToLower([]byte(name)), -1)
		startIndex := findSecondMin(indexes)
		allStartsIndexes = append(allStartsIndexes, startIndex)
		startIndexOfPlayMap[startIndex] = metadata{PlayName: name, Start: startIndex}
	}
	sort.Ints(allStartsIndexes)
	return startIndexOfPlayMap, allStartsIndexes
}

// findSecondMin returns the second minimum element in the array
func findSecondMin(arr []int) int {
	if len(arr) < 2 {
		return -1
	}
	sort.Ints(arr)
	return arr[1]
}
