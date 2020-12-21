package searcher_test

import (
	"reflect"
	"strings"
	"testing"

	"pulley.com/shakesearch/pkg/searcher"
)

//TODO Add more test cases
func TestSuffixArraySearch_Search(t *testing.T) {
	data := "The Pulley ShakeSearch Take-home Challenge"

	s := searcher.NewSuffixArraySearcher([]byte(data))

	word := "Pulley"

	req := searcher.Request{
		Query:            word,
		CaseSensitive:    false,
		ExactMatch:       false,
		CharBeforeQuery:  2,
		CharAfterQuery:   2,
		HighlightPreTag:  "<em>",
		HighlightPostTag: "</em>",
	}

	res, err := s.Search(req)
	if err != nil {
		t.Fatalf("unexptected error: %s", err)
	}

	if res.Query != req.Query {
		t.Errorf("query: want %s, got %s", req.Query, res.Query)
	}

	want := 1
	got := len(res.Hits)
	if got != want {
		t.Errorf("hits: want %d, got %d", want, got)
	}

	if len(res.Hits) == 0 {
		t.Errorf("hits: want 1 hit but is missing")
	} else if len(res.Hits) > 1 {
		t.Errorf("hits: want 1 hit but got %d", len(res.Hits))
	} else {
		gotHitQuery := res.Hits[0].Query
		if req.CaseSensitive {
			if gotHitQuery != req.Query {
				t.Errorf("hit: sensitive case: got %s, want %s", gotHitQuery, req.Query)
			}
		} else {
			if !strings.EqualFold(gotHitQuery, req.Query) {
				t.Errorf("hit: insensitive case: got %s and %s are not insensitive equal", gotHitQuery, req.Query)
			}
		}
	}

	if len(res.Highlights) == 0 {
		t.Errorf("highlights: want 1 hit but is missing")
	} else if len(res.Highlights) > 1 {
		t.Errorf("highlights: want 1 hit but got %d", len(res.Hits))
	} else {
		gotHighlightsText := res.Highlights[0].Text
		wantHiglightsText := "e <em>Pulley</em> S"
		if gotHighlightsText != wantHiglightsText {
			t.Errorf("highlights: got %q, want %q", gotHighlightsText, wantHiglightsText)
		}

		wantMatchedWords := []string{"Pulley"}
		if !reflect.DeepEqual(wantMatchedWords, res.Highlights[0].MatchedWords) {
			t.Errorf("highlights:matched_words: want %v, got %v", wantMatchedWords, res.Highlights[0].MatchedWords)
		}

	}
}
