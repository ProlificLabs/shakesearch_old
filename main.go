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
)

var (
	endOfSentence = []string{
		".",
		"?",
		"!",
		"\"",
	}
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

// Searcher allows for keyword search in the complete works of Shakespare.
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

// Load pulls data from a file and puts the contents into a suffixarray.SuffixArray for quick keyword search.
func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}

	// This will be in a browser, store with linebreaks to make it easier to read.
	body := strings.Replace(string(dat), "\r\n", "<br>", -1)
	s.CompleteWorks = body

	// Store the data as lowercase in the suffix array.
	dataSource := strings.ToLower(body)
	s.SuffixArray = suffixarray.New([]byte(dataSource))
	return nil
}

// Search looks up a string in our data store in memory.
// Searches are case insensitive.
func (s *Searcher) Search(query string) []string {
	idxs := s.SuffixArray.Lookup([]byte(strings.ToLower(query)), -1)
	results := []string{}
	for _, idx := range idxs {
		result := removeSentenceFragments(s.CompleteWorks[idx-250 : idx+250])
		// Don't put empty results in the search
		if result == "" {
			continue
		}
		results = append(results, result)
	}
	return results
}

// min finds the lowest number in a list of numbers.
// it will return 0 given an empty list
func min(numbers []int) int {
	if len(numbers) == 0 {
		return 0
	}

	if len(numbers) == 1 {
		return numbers[0]
	}

	smallest := numbers[0]
	for _, number := range numbers[1:] {
		if smallest < number {
			smallest = number
		}
	}
	return smallest
}

func max(numbers []int) int {
	if len(numbers) == 0 {
		return 0
	}

	if len(numbers) == 1 {
		return numbers[0]
	}

	largest := numbers[0]
	for _, number := range numbers[1:] {
		if largest > number {
			largest = number
		}
	}
	return largest
}

// removeSentenceFragments removes sentence fragments
func removeSentenceFragments(text string) string {
	if text == "" {
		return text
	}

	var firstSentenceEnds []int
	for _, character := range endOfSentence {
		firstSentenceIndex := strings.Index(text, character)
		// Prevent not found from messing with the min.
		if firstSentenceIndex == -1 {
			firstSentenceIndex = 0
		}
		firstSentenceEnds = append(firstSentenceEnds, firstSentenceIndex)
	}

	// The first sentence ends right after its punctuation.
	firstSentenceEnd := min(firstSentenceEnds) + 1
	text = text[firstSentenceEnd:]
	return text
}
