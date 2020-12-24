package main

import (
	"bytes"
	"encoding/json"
	"fmt"
    "html/template"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"net/http"
	"os"
    "regexp"
    "strings"
    "unicode"
)


type Pair struct {
    Name string
    Value string
}

type Person struct {
    Value string
}

func main() {
	searcher := Searcher{}
	err := searcher.Load("workslist.txt", "worksbody.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

    http.HandleFunc("/view", handleDropdown(searcher))
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
    WorksTitles []string
	CompleteWorks string
	SuffixArray   *suffixarray.Index
}

func handleDropdown(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
        var fruits = map[string]interface{}{
				 "Apple":  "apple",
				 "Orange": "orange",
				 "Pear":   "pear",
				 "Grape":  "grape",
			}
        t, err := template.ParseFiles("static/index.html")
        if err != nil {
            panic(err)
        }
        t.Execute(w, fruits)
    }
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

func (s *Searcher) Load(titlefile string, bodyfile string) error {
    titles, err := ioutil.ReadFile(titlefile)
    if err!= nil {
        return fmt.Errorf("Load: %w", err)
    }
    s.WorksTitles = strings.Split(string(titles), "\n\n")
	dat, err := ioutil.ReadFile(bodyfile)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.CompleteWorks = string(dat)
	s.SuffixArray = suffixarray.New(dat)
	return nil
}

func (s *Searcher) Search(query string) []string {
    results := []string{}
    if (query[0] == '"' && query[len(query)-1] == '"') {
        idxs := s.SuffixArray.Lookup([]byte(query[1:len(query)-1]), -1)
        for _, idx := range idxs {
            lines := s.GetLines(idx)
            results = append(results, lines)
        }
    } else {
        query = regexp.QuoteMeta(query)
        if !s.ContainsUpper(query) {
            query = "(?i)" + query
        }
        if (query[len(query)-2:] == "ed") {
            query = query[:len(query)-2] + "[e']d"
        }

        query = strings.ReplaceAll(query, " ", "[.!?'\"\\s\\[\\]\\(\\)-]")
        reg := regexp.MustCompile(query)
        idxs := s.SuffixArray.FindAllIndex(reg, -1)
        for _, idx := range idxs {
            lines := s.GetLines(idx[0])
            results = append(results, lines)
        }
	}
	return results
}

func (s *Searcher) GetLines(idx int) string {
    result := s.CompleteWorks[idx-100 : idx+100]
    lines := strings.Split(result, "\n")
    count := 0
    for i,line := range lines {
        count += len(line)
        if count > 100 {
            result = strings.Join(lines[i-1 : i+2],"<br>")
            break
        }
    }
    return result

}

func (s *Searcher) ContainsUpper(str string) bool {
    for _,r := range str {
        if unicode.IsUpper(r) && unicode.IsLetter(r) {
            return true
        }
    }
    return false
}
