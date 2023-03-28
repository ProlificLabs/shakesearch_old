package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"path/filepath"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"strconv"
	// "io/fs"
	// "sync"
)

func main() {
	//loop through theworks directory - the directory containing all of the works
	dir := "./theworks/"

    files, err := os.ReadDir(dir)
    if err != nil {
        log.Fatal(err)
        return
    }

	// load files in to a slice of Searchers(struct type) using goroutines
	// update: goroutines are not necessarily better for this task. decided not to use them,
	// but left their code in the comments just in case
	worksSlice := [] Searcher{}
	// var load_wg sync.WaitGroup


	for _, file := range files {
		// load_wg.Add(1)

		// go func(file fs.DirEntry){
		// 	defer load_wg.Done()
			if file.IsDir() {
				log.Fatal("directory found within the works files")
			} else {
				fmt.Printf("Loading File: %s\n", file.Name())
				searcher := Searcher{}
				searcher.title = strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
				err := searcher.Load(dir + file.Name())
				if err != nil {
					log.Fatal(err)
				}
				worksSlice = append(worksSlice, searcher)
			}
		//}(file)
    }
	// load_wg.Wait()
	fmt.Println("File loading complete.")

	//serve static files and run server
	fs := http.FileServer(http.Dir("./build"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handleSearch(worksSlice))

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

// struct to hold content of each work
type Searcher struct {
	title		  string
	works 		  string
	SuffixArray   *suffixarray.Index
}

//handle incoming search queries
func handleSearch(worksSlice []Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// enable CORS
		// enableCors(&w)
		// get url parameters
		query, query_ok := r.URL.Query()["q"]
		if !query_ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}
		caseSensitive_param := r.URL.Query().Get("caseSensitive")
		if caseSensitive_param == "" {
			http.Error(w, "missing caseSensitive parameter", http.StatusBadRequest)
			return
		}
		// Parse the boolean parameter
		caseSensitive, caseSensitive_ok := strconv.ParseBool(caseSensitive_param)
		if caseSensitive_ok != nil {
			http.Error(w, "invalid caseSensitive parameter", http.StatusBadRequest)
			return
		}
		wholeWord_param := r.URL.Query().Get("wholeWord")
		if wholeWord_param == "" {
			http.Error(w, "missing wholeWord parameter", http.StatusBadRequest)
			return
		}
		// Parse the boolean parameter
		wholeWord, wholeWord_ok := strconv.ParseBool(wholeWord_param)
		if wholeWord_ok != nil {
			http.Error(w, "invalid wholeWord parameter", http.StatusBadRequest)
			return
		}
		//search for matches in the works
		results := []map[string]string{}
		for _, work := range worksSlice {
			results = append(results, work.Search(query[0], caseSensitive, wholeWord))
		}
		//write answer to client
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

// load files to searcher struct
func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.works = string(dat)
	s.SuffixArray = suffixarray.New(dat)
	return nil
}

//search for matches within the works using regex patterns
func (s *Searcher) Search(query string, caseSensitive bool, wholeWord bool) map[string]string {
	var pattern *regexp.Regexp
	// Compile a regular expression pattern for the query
	if caseSensitive {
		if wholeWord {
			pattern = regexp.MustCompile("\\b" + regexp.QuoteMeta(query) + "\\b")
		} else{
			pattern = regexp.MustCompile(regexp.QuoteMeta(query))
		}
	} else {
		if wholeWord {
			pattern = regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(query) + `\b`)
		} else{
			query = strings.ToLower(query)
			pattern = regexp.MustCompile("(?i)" + regexp.QuoteMeta(query))
		}
	}

	// Use the FindAllIndex function to find all matches of the query in the longString
	matches := s.SuffixArray.FindAllIndex(pattern, -1)
	// Return the number of matches found
	var builder strings.Builder

	// Loop through the longString and insert HTML tags around the matched words
	lastMatchEnd := 0
	for _, match := range matches {

		// Append the part of the longString before the match
		builder.WriteString(s.works[lastMatchEnd:match[0]])

		// Append the matched word with HTML span tags
		builder.WriteString(fmt.Sprintf("<span>%s</span>", s.works[match[0]:match[1]]))
		lastMatchEnd = match[1]
	}

	// Append the part of the longString after the last match
	builder.WriteString(s.works[lastMatchEnd:])

	// Return the new string with highlighted words
	artPiece := builder.String()
	wordCount := len(matches)
	return map[string]string{"title" : s.title, "artPiece" : artPiece, "wordCount" : strconv.Itoa(wordCount)}
}

// // enables CORS
// func enableCors(w *http.ResponseWriter) {
// (*w).Header().Set("Access-Control-Allow-Origin", "*")
// }