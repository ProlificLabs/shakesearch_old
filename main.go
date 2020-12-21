package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"pulley.com/shakesearch/pkg/searcher"
)

func main() {
	dat, err := ioutil.ReadFile("./completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	suffixarraySearcher := searcher.NewSuffixArraySearcher(dat)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handleSearch(suffixarraySearcher))

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

func handleSearch(s searcher.Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}

		searchRequest := SearchRequest{}
		if err := searchRequest.Bind(r); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		req := searcher.Request{
			Query:            searchRequest.Q,
			CaseSensitive:    searchRequest.Sensitive,
			ExactMatch:       searchRequest.ExactMatch,
			CharBeforeQuery:  searchRequest.Before,
			CharAfterQuery:   searchRequest.After,
			HighlightPreTag:  "<em>",
			HighlightPostTag: "</em>",
		}

		res, err := s.Search(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
			return
		}

		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		if err := enc.Encode(res); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
}

type SearchRequest struct {
	Q          string
	Sensitive  bool
	ExactMatch bool
	Before     int
	After      int
}

func (sr *SearchRequest) Bind(r *http.Request) error {
	const (
		defaultSensitive = false
		defaultBefore    = 215
		defaultAfter     = 215
	)

	query := r.URL.Query()

	q := query.Get("q")
	if q == "" {
		return errors.New("missing search query in URL params")
	}

	if strings.HasPrefix(q, `"`) && strings.HasSuffix(q, `"`) {
		sr.ExactMatch = true
		q = strings.TrimPrefix(q, `"`)
		q = strings.TrimSuffix(q, `"`)
	}

	sr.Q = q

	sensitive, err := strconv.ParseBool(query.Get("sensitive"))
	if err != nil {
		sr.Sensitive = defaultSensitive
	} else {
		sr.Sensitive = sensitive
	}

	before, err := strconv.Atoi(query.Get("before"))
	if err != nil {
		sr.Before = defaultBefore
	} else {
		sr.Before = before
	}

	after, err := strconv.Atoi(query.Get("after"))
	if err != nil {
		sr.After = defaultAfter
	} else {
		sr.After = after
	}

	return nil

}
