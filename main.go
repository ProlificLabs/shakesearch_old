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
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var (
	indexMap    sync.Map //map[string]*Searcher // title: *Searcher
	titlesList  []string
	contentDir  string = filepath.Join("static", "completeWorks")
	tocFilename string = "0_toc"
)

// Searcher is the struct holding the title-wise content and indexes
type Searcher struct {
	id              int
	title           string
	work            string
	toc             string
	characters      string
	workIndex       *suffixarray.Index // assert non-nil before every access
	tocIndex        *suffixarray.Index // assert non-nil before every access
	charactersIndex *suffixarray.Index // assert non-nil before every access
}

// SearcherResult is a struct to encapsulate the sub-category results from each dataset
type SearcherResult struct {
	TOC        []string `json:"toc"`
	Characters []string `json:"characters"`
	Work       []string `json:"work"`
}

func searchWithFilter(s Searcher, filters []string, queries []string, limit int) SearcherResult {
	sresult := SearcherResult{}
	for _, filter := range filters {
		switch filter {
		case "toc":
			if s.tocIndex == nil {
				continue
			}
			for _, query := range queries {
				idxs := s.tocIndex.Lookup([]byte(strings.ToLower(query)), limit)
				for _, idx := range idxs {
					l := len(s.toc)
					var lb, ub int
					if idx-40 < 0 {
						lb = 0
					} else {
						lb = idx - 40
					}
					if idx+40 > l {
						ub = l
					} else {
						ub = idx + 40
					}
					sresult.TOC = append(sresult.TOC, s.toc[lb:ub])
				}
			}
		case "work":
			if s.workIndex == nil {
				continue
			}
			for _, query := range queries {
				idxs := s.workIndex.Lookup([]byte(strings.ToLower(query)), limit)
				for _, idx := range idxs {
					l := len(s.work)
					var lb, ub int
					if idx-250 < 0 {
						lb = 0
					} else {
						lb = idx - 250
					}
					if idx+250 > l {
						ub = l
					} else {
						ub = idx + 250
					}
					sresult.Work = append(sresult.Work, s.work[lb:ub])
				}
			}
		case "characters":
			if s.charactersIndex == nil {
				continue
			}
			for _, query := range queries {
				idxs := s.charactersIndex.Lookup([]byte(strings.ToLower(query)), limit)
				for _, idx := range idxs {
					l := len(s.characters)
					var lb, ub int
					if idx-20 < 0 {
						lb = 0
					} else {
						lb = idx - 20
					}
					if idx+20 > l {
						ub = l
					} else {
						ub = idx + 20
					}
					sresult.Characters = append(sresult.Characters, s.characters[lb:ub])
				}
			}
		}
	}
	return sresult
}

func handleSearch() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		q, ok := r.URL.Query()["q"]
		if !ok || len(q) < 1 || q[0] == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query"))
			return
		}
		queries := strings.Split(q[0], " ")
		filters, ok := r.URL.Query()["f"]
		if !ok || filters[0] == "" || filters[0] == "none" || filters[0] == "undefined" {
			filters = append(filters, "toc") // If no filter is set, explicitly setting search scope to all
			filters = append(filters, "work")
			filters = append(filters, "characters")
		} else {
			filters = strings.Split(filters[0], " ")
		}

		dataset, ok := r.URL.Query()["d"]
		if !ok || len(dataset) < 1 {
			dataset = append(dataset, "all")
		} else {
			dataset = strings.Split(dataset[0], ",")
		}
		limit, err := strconv.Atoi(r.URL.Query().Get("l"))
		if err != nil {
			log.Printf("Setting result limit to -1. Error while fetching the result limit from query: %s", err.Error())
			limit = -1
		}
		if limit == 0 {
			limit = -1
		}

		wg := &sync.WaitGroup{}

		if dataset[0] == "all" || dataset[0] == "" {
			dataset = titlesList
			dataset = append(dataset, "mindex")
		}

		var resultsMap sync.Map // map[string][]string // dataset: SearcherResult

		for _, d := range dataset {
			wg.Add(1)
			go func(wg *sync.WaitGroup, d string) {
				defer wg.Done()
				v, ok := indexMap.Load(d)
				if !ok {
					log.Printf("Unknown dataset: %s", d)
					return
				}
				searcher, ok := v.(Searcher)
				if !ok {
					log.Printf("Error during type assertion of dataset %s", d)
					return
				}
				result := searchWithFilter(searcher, filters, queries, limit)
				resultsMap.Store(d, result)
			}(wg, d)
		}

		wg.Wait()

		m := make(map[string]SearcherResult)
		fn := func(key, value interface{}) bool {
			d := key.(string)
			res := value.(SearcherResult)
			m[d] = res
			return true
		}
		resultsMap.Range(fn)

		resultsJson, err := json.Marshal(finalizeResponse(m))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("json encoding failure"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(resultsJson)
	}
}

func finalizeResponse(m map[string]SearcherResult) map[string]map[string][]string { // input title:results struct | returns filter:title:result
	// Create finalResponse map
	finalResponse := make(map[string]map[string][]string)
	// Create maps for each tab/filter
	finalResponse["mindex"] = make(map[string][]string)
	finalResponse["toc"] = make(map[string][]string)
	finalResponse["characters"] = make(map[string][]string)
	finalResponse["work"] = make(map[string][]string)

	for title, sresult := range m {
		if title == "mindex" && len(sresult.TOC) > 0 {
			finalResponse["mindex"][""] = sresult.TOC
			continue
		}
		if len(sresult.TOC) > 0 {
			finalResponse["toc"][title] = sresult.TOC
		}
		if len(sresult.Characters) > 0 {
			finalResponse["characters"][title] = sresult.Characters
		}
		if len(sresult.Work) > 0 {
			finalResponse["work"][title] = sresult.Work
		}
	}
	for k, v := range finalResponse {
		if len(v) < 1 {
			delete(finalResponse, k)
		}
	}
	return finalResponse
}

func load(contentDir string) {
	toc, err := ioutil.ReadFile(path.Join(contentDir, tocFilename))
	if err != nil {
		log.Fatalf("Unable to load toc from directory: %s", err.Error())
	}

	titlesList = strings.Split(string(toc), "\n")
	wg := &sync.WaitGroup{}
	// Load main index to indexMap
	wg.Add(1)
	loadSearcherToMap(0, "mindex", wg)
	// Load all work to indexMap
	for i, title := range titlesList {
		wg.Add(1)
		go loadSearcherToMap(i+1, title, wg)
	}
	wg.Wait()
}

func loadSearcherToMap(id int, title string, wg *sync.WaitGroup) {
	defer wg.Done()
	files, err := filepath.Glob(fmt.Sprintf("%s/%d_*", contentDir, id))
	if err != nil {
		log.Printf("Error while fetching filenames for id: %d", id)
	}
	searchIndex := Searcher{
		id:    id,
		title: title,
	}
	for _, file := range files {
		switch {
		case strings.Contains(file, "toc"):
			dat, err := ioutil.ReadFile(file)
			if err != nil {
				log.Printf("Error while reading %s: %s", file, err.Error())
			}
			searchIndex.toc = string(dat)
			searchIndex.tocIndex = suffixarray.New(bytes.ToLower(dat))
			continue
		case strings.Contains(file, "characters"):
			dat, err := ioutil.ReadFile(file)
			if err != nil {
				log.Printf("Error while reading %s: %s", file, err.Error())
			}
			searchIndex.characters = string(dat)
			searchIndex.charactersIndex = suffixarray.New(bytes.ToLower(dat))
			continue
		default:
			dat, err := ioutil.ReadFile(file)
			if err != nil {
				log.Printf("Error while reading %s: %s", file, err.Error())
			}
			searchIndex.work = string(dat)
			searchIndex.workIndex = suffixarray.New(bytes.ToLower(dat))
		}
	}
	indexMap.Store(title, searchIndex)
}

func main() {
	load(contentDir)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/search", handleSearch())
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	fmt.Printf("Listening on port %s...", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
