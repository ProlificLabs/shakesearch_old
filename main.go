package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"pulley.com/shakesearch/pkg/render"
	"pulley.com/shakesearch/pkg/searcher"
)

func main() {
	dat, err := ioutil.ReadFile("./completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	suffixarraySearcher := searcher.NewSuffixArraySearcher(dat)
	handler := NewHandler(suffixarraySearcher)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handler.Search)

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

type Handler struct {
	Searcher searcher.Searcher
	Render   render.Render
}

func NewHandler(s searcher.Searcher) *Handler {
	return &Handler{
		Searcher: s,
		Render:   render.NewJsonRender(),
	}
}

func (h Handler) ChangerRender(r render.Render) {
	h.Render = r
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Render.Error(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	searchRequest := SearchRequest{}
	if err := searchRequest.Bind(r); err != nil {
		h.Render.Error(w, http.StatusBadRequest, err)
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

	res, err := h.Searcher.Search(req)
	if err != nil {
		h.Render.Error(w, http.StatusBadRequest, err)
		return
	}

	h.Render.Success(w, res)
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
