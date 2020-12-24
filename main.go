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
    "regexp"
    "strings"
    "unicode"
)

func main() {
	searcher := Searcher{}
	err := searcher.Load("workslist.txt", "worksbody.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))
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

type Searcher struct {
    WorksTitles []string
	CompleteWorks string
	SuffixArray   *suffixarray.Index
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
    //filter := regexp.MustCompile("[^a-zA-Z0-9\s.!?;:\'\"\\-]+")
    //query = filter.ReplaceAllString(query, "")
    if !s.ContainsUpper(query) {
        query = "(?i)" + query
    }
    reg := regexp.MustCompile(query)
	idxs := s.SuffixArray.FindAllIndex(reg, -1)
	results := []string{}
    curr_idx := 0
	for _, idx := range idxs {
        if idx[0] >= curr_idx {
            lines := s.GetLines(idx[0])
            results = append(results, lines)
        }
	}
	return results
}

func (s *Searcher) GetLines(idx int) string {
    result := s.CompleteWorks[idx-150 : idx+150]
    lines := strings.Split(result, "\n")
    count := 0
    for i,line := range lines {
        count += len(line)
        if count > 150 {
            result = strings.Join(lines[i-1 : i+2],"")
        }
    }
    result = strings.Join(lines[:],"")
    result = strings.Replace(result, "\n", "<br>", -1)
    return result

//    start, end := -1, -1
//    sn, en := 0, 0
//    for i := 0; start!=-1 && end!=-1 && i < 250; i++ {
//        if string(s.CompleteWorks[idx-i]) == string("\n") {
//            sn += 1
//            if sn == 1 {
//                start = idx-i+1
//            }
//        }
//        if string(s.CompleteWorks[idx+i]) == string("\n") {
//            en += 1
//            if en == 1 {
//                end = idx+i
//            }
//        }
//    }
//    if start == -1 || end == -1 {
//        start,end = idx-100, idx+100
//    }
//    result := s.CompleteWorks[start : end]
//    result = strings.Replace(result, "\n", "<br>", -1)
//    return result, end
}

func (s *Searcher) ContainsUpper(str string) bool {
    for _,r := range str {
        if unicode.IsUpper(r) && unicode.IsLetter(r) {
            return true
        }
    }
    return false
}
