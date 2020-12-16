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
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.CompleteWorks = string(dat)
	s.SuffixArray = suffixarray.New([]byte(strings.ToLower(string(dat))))
	return nil
}

func (s *Searcher) Search(query string) []*returnStruct {
	idxs := s.SuffixArray.Lookup([]byte(strings.ToLower(query)), -1)
	results := []*returnStruct{}

	// Index to compact result
	lowerBounds := make(map[int]map[int]bool)

	for _, idx := range idxs {
		// Divide into blocks of 500 characters

		lb := (idx / 250) * 250
		if _, ok := lowerBounds[lb]; !ok {
			lowerBounds[lb] = make(map[int]bool)
		}
		lowerBounds[lb][idx-lb] = true
	}

	mut := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	sem := semaphore.NewWeighted(10) // don't spawn too many threads

	for lb, positionsMap := range lowerBounds {
		wg.Add(1)
		sem.Acquire(context.Background(), 1)

		go func(lb int, positionsMap map[int]bool) {

			positions := []int{}
			for k := range positionsMap {
				positions = append(positions, k)
			}
			sort.SliceStable(positions, func(i, j int) bool { return positions[i] < positions[j] })

			mut.Lock()
			results = append(results, &returnStruct{Text: s.CompleteWorks[lb : lb+500], Positions: positions})
			mut.Unlock()
			wg.Done()
			sem.Release(1)

		}(lb, positionsMap)

	}

	wg.Wait()

	return results
}

type returnStruct struct {
	Positions []int  `json:"positions"`
	Text      string `json:"text"`
}
