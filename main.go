package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"

	"golang.org/x/sync/semaphore"
)

const (
	maxNumberOfThreads      = 10
	paragraphSize           = 500
	subseqFragmentationSize = 3
)

type returnStruct struct {
	Positions []int  `json:"positions"`
	Text      string `json:"text"`
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

	fmt.Printf("Listening on port %s...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

type Searcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
	Occurrences   map[string]map[int]bool
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}
		results := searcher.Search(query[0], false, 2)
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
	asLowerCase := strings.ToLower(string(dat))

	s.CompleteWorks = string(dat)
	s.SuffixArray = suffixarray.New([]byte(asLowerCase))

	s.Occurrences = make(map[string]map[int]bool)

	for i := 0; i < len(asLowerCase)-subseqFragmentationSize; i++ {
		if _, ok := s.Occurrences[asLowerCase[i:i+subseqFragmentationSize]]; !ok {
			s.Occurrences[asLowerCase[i:i+subseqFragmentationSize]] = make(map[int]bool)
		}
		s.Occurrences[asLowerCase[i:i+subseqFragmentationSize]][i] = true
	}

	return nil
}

// Fault-tolerant merge of continuous interval into one lower bound position
func mergePositions(positions []int, maxMismatches, minLength int) []int {
	sort.SliceStable(positions, func(i, j int) bool { return positions[i] < positions[j] })

	result := []int{}

	continuousPositions := []int{}

	lastPos := 0
	for _, pos := range positions {
		if (len(continuousPositions) > 0) && (pos-lastPos >= maxMismatches) {
			// split found
			if continuousPositions[len(continuousPositions)-1]+subseqFragmentationSize-continuousPositions[0] >= minLength { // ensure minimum length is met
				result = append(result, continuousPositions[0])
			}
			continuousPositions = []int{}
		}
		continuousPositions = append(continuousPositions, pos)
		lastPos = pos
	}

	if len(continuousPositions) > 0 {
		result = append(result, continuousPositions[0])
	}

	return result

}

// Search(query string, exact bool)
// query: word to match for
// exact: match exact word or similar words
// maxMismatches: number of mismatches allowed (requires exact = false)
// Description: Fault-tolerant substring search.
// Author: Lemmer EL ASSAL
func (s *Searcher) Search(query string, exact bool, maxMismatches int) []*returnStruct {

	query = strings.ToLower(query)

	idxs := []int{}

	if exact {
		idxs = s.SuffixArray.Lookup([]byte(strings.ToLower(query)), -1)
	} else {

		uniqueSubsequences := make(map[string]bool)
		for i := 0; i < len(query)-subseqFragmentationSize+1; i++ {
			uniqueSubsequences[query[i:i+subseqFragmentationSize]] = true
		}

		for sub := range uniqueSubsequences {
			for pos := range s.Occurrences[sub] {
				idxs = append(idxs, pos)
			}
		}

		idxs = mergePositions(idxs, maxMismatches, len(query)-maxMismatches)
	}

	results := []*returnStruct{}

	// Index to compact result
	lowerBounds := make(map[int]map[int]bool)

	for _, idx := range idxs {
		// Divide into blocks of paragraphSize characters

		lb := (idx / paragraphSize) * paragraphSize
		if _, ok := lowerBounds[lb]; !ok {
			lowerBounds[lb] = make(map[int]bool)
		}
		lowerBounds[lb][idx-lb] = true
	}

	mut := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	sem := semaphore.NewWeighted(maxNumberOfThreads) // limit number of threads spawned

	for lb, positionsMap := range lowerBounds {
		wg.Add(1)
		sem.Acquire(context.Background(), 1)

		go func(lb int, positionsMap map[int]bool) {

			positions := []int{}
			for k := range positionsMap {
				positions = append(positions, k)
			}
			sort.SliceStable(positions, func(i, j int) bool { return positions[i] < positions[j] })

			ub := lb + paragraphSize
			if ub > len(s.CompleteWorks) {
				ub = len(s.CompleteWorks)
			}

			mut.Lock()
			results = append(results, &returnStruct{Text: s.CompleteWorks[lb:ub], Positions: positions})
			mut.Unlock()
			sem.Release(1) // Alternatively: parallel threads with data- and "done" channel - Lemmer

			wg.Done()

		}(lb, positionsMap)

	}

	wg.Wait()

	return results
}
