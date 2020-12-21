package searcher

import (
	"bytes"
	"index/suffixarray"
	"sort"
	"strings"
	"unicode/utf8"
)

type SuffixArraySearcher struct {
	data         string
	index        *suffixarray.Index
	indextoLower *suffixarray.Index
}

func NewSuffixArraySearcher(data []byte) *SuffixArraySearcher {
	return &SuffixArraySearcher{
		index:        suffixarray.New(data),
		indextoLower: suffixarray.New(bytes.ToLower(data)),
		data:         string(data),
	}
}

func (s *SuffixArraySearcher) lookup(query string, n int, sensitive bool) []int {
	if sensitive {
		return s.index.Lookup([]byte(query), n)
	}

	return s.indextoLower.Lookup([]byte(strings.ToLower(query)), n)
}

func (s *SuffixArraySearcher) Search(req Request) (*Response, error) {
	query := req.Query
	nbResultMax := -1

	tokens := []string{
		query,
	}

	if !req.ExactMatch {
		tokens = analyze(query)
	}

	var hits []*Hit
	for _, t := range tokens {
		idxs := s.lookup(t, nbResultMax, req.CaseSensitive)
		newHits := s.search(t, idxs)
		hits = append(hits, newHits...)
	}

	sort.Sort(Hits(hits))

	highlights := highlightText(s.data, hits, req.CharBeforeQuery, req.CharAfterQuery)

	return &Response{
		Query:      req.Query,
		Hits:       hits,
		Highlights: highlights,
	}, nil
}

func (s *SuffixArraySearcher) search(text string, idxs []int) []*Hit {
	var hits []*Hit

	textLen := utf8.RuneCountInString(text)

	for id, idx := range idxs {
		sourceText := s.data[idx : idx+textLen]

		hit := Hit{
			Query: sourceText,
			Order: id + 1,
			Start: idx,
			End:   idx + textLen,
		}

		hits = append(hits, &hit)
	}

	return hits
}

func highlightText(data string, hits []*Hit, before int, after int) []*Highlight {
	var size = len(data)
	var start int
	var end int

	var highlights []*Highlight
	var matchedWords []string

	for i, h := range hits {
		if i == 0 {
			matchedWords = append(matchedWords, data[h.Start:h.End])
			start, end = calculRangeText(h, size, before, after)
			if len(hits) == 1 {
				highlights = append(highlights, &Highlight{
					Text:         data[start:end],
					MatchedWords: uniqueString(matchedWords),
				})
			}
			continue
		}

		if h.Start-before > end {
			highlights = append(highlights, &Highlight{
				Text:         data[start:end],
				MatchedWords: uniqueString(matchedWords),
			})
			start, end = calculRangeText(h, size, before, after)
			matchedWords = []string{
				data[h.Start:h.End],
			}
			continue
		}

		matchedWords = append(matchedWords, data[h.Start:h.End])
		end = h.End + after
		if end >= size {
			end = size
		}

	}

	for _, h := range highlights {
		for _, word := range h.MatchedWords {
			h.Text = strings.ReplaceAll(h.Text, word, "<em>"+word+"</em>")
			h.Text = strings.ReplaceAll(h.Text, "\n", "<br />")
		}
	}

	return highlights
}

func uniqueString(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

func calculRangeText(h *Hit, size int, before int, after int) (start int, end int) {
	start = h.Start - before
	if start < 0 {
		start = 0
	}

	end = h.End + after
	if end >= size {
		end = size
	}

	return
}

func tokenize(text string) []string {
	return strings.Fields(text)
}

func lowercaseFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = strings.ToLower(token)
	}
	return r
}

var stopwords = map[string]struct{}{
	"a": {}, "and": {}, "be": {}, "have": {}, "i": {},
	"in": {}, "of": {}, "that": {}, "the": {}, "to": {},
}

func stopwordFilter(tokens []string) []string {
	r := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if _, ok := stopwords[token]; !ok {
			r = append(r, token)
		}
	}
	return r
}

func analyze(text string) []string {
	tokens := tokenize(text)
	tokens = lowercaseFilter(tokens)
	tokens = stopwordFilter(tokens)
	return tokens
}
