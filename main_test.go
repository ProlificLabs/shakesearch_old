package main

import (
	"index/suffixarray"
	"reflect"
	"strings"
	"testing"
)

func TestSearch(t *testing.T){
	searcher := Searcher{}
	searcher.CompleteWorks = "Hamlet.Hamlet.Hamlet."

	sarray := suffixarray.New([]byte(strings.ToLower(searcher.CompleteWorks)))
  searcher.SuffixArray = sarray

	query := "Hamlet"
	expectedResults := []string{
		"Hamlet.Hamlet.Hamlet.",
		"Hamlet.Hamlet.Hamlet.",
		"Hamlet.Hamlet.Hamlet.",
	}

	resultsFull := searcher.Search(query, 1, 100)

	if !reflect.DeepEqual(expectedResults, resultsFull) {
		t.Errorf("expected %v, but got %v", expectedResults, resultsFull)
	}

	resultsPartail := searcher.Search(query, 2, 1)
	expectedPartailResults := []string{
		"Hamlet.Hamlet.Hamlet.",
	}

	if !reflect.DeepEqual(expectedPartailResults, resultsPartail) {
		t.Errorf("expected %v, but got %v", expectedPartailResults, resultsPartail)
	}
}
