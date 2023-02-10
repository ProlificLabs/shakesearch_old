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

	//	"strconv"
	"strings"
	// "text/template"
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
var cache = make(map[string][]string)
var stopWords = []string{"a", "an", "and", "are", "as", "at", "be", "but", "by", "for", "if", "in", "into", "is", "it", "no", "not", "of", "on", "or", "such", "that", "the", "their", "then", "there", "these", "they", "this", "to", "was", "will", "with"}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}
		normalizedQuery := normalizeQuery(query[0])
		//Adding cache map to imrove subsequent search queries
		var results []string
		if cachedResults, ok := cache[normalizedQuery]; ok {
		results = cachedResults
		} else {
		results = searcher.Search(normalizedQuery)
		cache[normalizedQuery] = results
		}
		//results := searcher.Search(normalizedQuery)
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



// 1.Handles Query Normalization and removes the stop words
func normalizeQuery(query string) string {
	// Convert the query to lowercase
	query = strings.ToLower(query)
	// Trim any leading or trailing whitespace
	query = strings.TrimSpace(query)
	// Split the query into individual words
	words := strings.Split(query, " ")
	// Create a new slice to store the filtered words
	filteredWords := []string{}
	// Loop through each word in the query
	for _, word := range words {
		// Check if the word is not a stop word
		if !isStopWord(word) {
			filteredWords = append(filteredWords, word)
		}
	}
	// Join the filtered words back into a single string
	return strings.Join(filteredWords, " ")
}

func isStopWord(word string) bool {
	// Loop through each stop word
	for _, stopWord := range stopWords {
		// If the word matches a stop word, return true
		if word == stopWord {
			return true
		}
	}
	// If the word does not match any stop words, return false
	return false
}

func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.CompleteWorks = string(dat)
	s.SuffixArray = suffixarray.New(dat)
	return nil
}

func (s *Searcher) Search(query string) []string {
	// Normalize the query string by converting to lowercase
	query = strings.ToLower(query)

	idxs := s.SuffixArray.Lookup([]byte(query), -1)
	results := []string{}
	for _, idx := range idxs {
		match := s.CompleteWorks[idx-250 : idx+250]
		// 2. Bold Italicisizing the match
		match = strings.ReplaceAll(match, query, "<h3 style='color:#2996e8;font-weight:bold,font-size:20px'><i>"+query+"</i></h3>")
		results = append(results, match)
	}
	return results
}
