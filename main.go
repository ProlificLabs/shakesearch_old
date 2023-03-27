package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
)

func main() {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./client/build"))
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

type Result struct {
	Snippet string
	Matches int
}
type Searcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		encodedTokens, ok := r.URL.Query()["q"]
		if !ok || len(encodedTokens[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}
		var tokens []string
		for _, encodedToken := range encodedTokens {
			decodedToken, ok := url.QueryUnescape(encodedToken)
			if ok != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("error decoding search query"))
			}
			tokens = append(tokens, decodedToken)
		}
		results := searcher.Search(tokens)
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(buf.Bytes())
	}
}

func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.CompleteWorks = string(dat)
	s.SuffixArray = suffixarray.New(bytes.ToLower(dat))
	return nil
}

func (s *Searcher) Search(tokens []string) []Result {
	mergedIdxs := []int{}
	simpleTokens := []string{}
	exactTokens := []string{}
	// group by token type
	for _, token := range tokens {
		if strings.HasPrefix(token, "\"") && strings.HasSuffix(token, "\"") {
			exactTokens = append(exactTokens, token)
		} else {
			simpleTokens = append(simpleTokens, token)
		}
	}
	// Find all indexes matches with "exact" tokens
	indexesOfExactTokens, indexesOfSimpleTokens := [][]int{}, [][]int{}
	if len(exactTokens) > 0 {
		regexExactTokens := generateRegex(exactTokens, "exact")
		indexesOfExactTokens = s.SuffixArray.FindAllIndex(regexExactTokens, -1)
	}
	// Find all indexes matches with "simple" tokens
	if len(simpleTokens) > 0 {
		regexSimpleTokens := generateRegex(simpleTokens, "simple")
		indexesOfSimpleTokens = s.SuffixArray.FindAllIndex(regexSimpleTokens, -1)
	}

	tokenIdxs := []int{}
	for _, match := range indexesOfExactTokens {
		tokenIdxs = append(tokenIdxs, match[0])
	}
	for _, match := range indexesOfSimpleTokens {
		tokenIdxs = append(tokenIdxs, match[0])
	}
	mergedIdxs = append(mergedIdxs, tokenIdxs...)
	// sort merged indexes
	sort.Slice(mergedIdxs, func(i, j int) bool {
		return mergedIdxs[i] < mergedIdxs[j]
	})
	// consolidate results
	results := []Result{}
	for i := 0; i < len(mergedIdxs); i++ {
		lenOfMatches := 1
		// seek for matches in the following 250 characters in order to mitigate overlapping snippets
		if (i + 1) < len(mergedIdxs) {
			var snippet string
			if len(s.CompleteWorks) < mergedIdxs[i]+250 {
				snippet = s.CompleteWorks[mergedIdxs[i]:]
			} else {
				snippet = s.CompleteWorks[mergedIdxs[i] : mergedIdxs[i]+250]
			}
			//snippet := s.CompleteWorks[mergedIdxs[i] : mergedIdxs[i]+250]
			matchesOfExactTokens, matchesOfSimpleTokens := []string{}, []string{}
			if len(exactTokens) > 0 {
				matchesOfExactTokens = getMatches(exactTokens, "exact", snippet)
			}
			if len(simpleTokens) > 0 {
				matchesOfSimpleTokens = getMatches(simpleTokens, "simple", snippet)
			}
			lenOfMatches = len(matchesOfExactTokens) + len(matchesOfSimpleTokens)

			// create the result with 260 characters range in order to trim at the nearest space
			result := Result{}
			var lowerOffset, upperOffset int
			if mergedIdxs[i] < 260 {
				lowerOffset = mergedIdxs[i]
			} else {
				lowerOffset = 260
			}
			if mergedIdxs[i] > len(s.CompleteWorks)-260 {
				upperOffset = len(s.CompleteWorks) - mergedIdxs[i]
			} else {
				upperOffset = 260
			}
			result.Snippet = s.CompleteWorks[mergedIdxs[i]-lowerOffset : mergedIdxs[i]+upperOffset]
			result.Matches = lenOfMatches
			results = append(results, result)

			// skip the next N matches
			if lenOfMatches > 1 {
				i = i + (lenOfMatches - 1)
			}
		}
	}
	return results
}

func getMatches(tokens []string, tokenType string, snippet string) []string {
	regex := generateRegex(tokens, tokenType)
	matches := []string{}
	if len(tokens) > 0 {
		matches = regex.FindAllString(strings.ToLower(snippet), -1)
	}
	return matches
}

func generateRegex(tokens []string, tokenType string) *regexp.Regexp {
	escapedTokens := []string{}
	for _, token := range tokens {
		if tokenType == "exact" {
			escapedTokens = append(escapedTokens, strings.ReplaceAll(regexp.QuoteMeta(token), "\"", ""))
		} else if tokenType == "simple" {
			escapedTokens = append(escapedTokens, regexp.QuoteMeta(token))
		}
	}
	var regexStr string
	if tokenType == "exact" {
		regexStr = "\\b(" + strings.Join(escapedTokens, "|") + ")\\b"
	} else if tokenType == "simple" {
		regexStr = "(" + strings.Join(escapedTokens, "|") + ")"
	}
	regex := regexp.MustCompile(regexStr)
	return regex
}
