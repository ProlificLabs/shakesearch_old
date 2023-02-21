package main

import (
	"log"
	"strings"
	"testing"
)

func TestSearcher_findPlayName(t *testing.T) {

	tests := []struct {
		name  string
		query string
		want  string
	}{
		{
			name:  "test1",
			query: "a she-lamb of a twelvemonth to crooked-pated",
			want:  "AS YOU LIKE IT",
		},
		{
			name:  "test2",
			query: "thy youngest daughter does not love thee least",
			want:  "THE TRAGEDY OF KING LEAR",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Searcher{}
			err := s.Load("completeworks.txt")
			if err != nil {
				log.Fatal(err)
			}
			indexes := s.LowerSuffixArray.Lookup([]byte(strings.ToLower(tt.query)), 1)
			t.Logf("indexes: %v", indexes)
			if len(indexes) == 0 {
				t.Fatalf("no indexes found for query: %s", tt.query)
			}
			if got := s.findPlayName(indexes[0]); strings.ToLower(got) != strings.ToLower(tt.want) {
				t.Errorf("findPlayName() = %v, want %v", got, tt.want)
			}
		})
	}
}
