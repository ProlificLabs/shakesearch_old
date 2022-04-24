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
	"strconv"
)

var playTitles = []string{
	"ALL’S WELL THAT ENDS WELL",
	"ANTONY AND CLEOPATRA",
	"AS YOU LIKE IT",
	"THE COMEDY OF ERRORS",
	"THE TRAGEDY OF CORIOLANUS",
	"CYMBELINE",
	"THE TRAGEDY OF HAMLET, PRINCE OF DENMARK",
	"THE FIRST PART OF KING HENRY THE FOURTH",
	"THE SECOND PART OF KING HENRY THE FOURTH",
	"THE LIFE OF KING HENRY V",
	"THE FIRST PART OF HENRY THE SIXTH",
	"THE SECOND PART OF KING HENRY THE SIXTH",
	"THE THIRD PART OF KING HENRY THE SIXTH",
	"KING HENRY THE EIGHTH",
	"KING JOHN",
	"THE TRAGEDY OF JULIUS CAESAR",
	"THE TRAGEDY OF KING LEAR",
	"LOVE’S LABOUR’S LOST",
	"MACBETH",
	"MEASURE FOR MEASURE",
	"THE MERCHANT OF VENICE",
	"THE MERRY WIVES OF WINDSOR",
	"A MIDSUMMER NIGHT’S DREAM",
	"MUCH ADO ABOUT NOTHING",
	"OTHELLO, THE MOOR OF VENICE",
	"PERICLES, PRINCE OF TYRE",
	"KING RICHARD THE SECOND",
	"KING RICHARD THE THIRD",
	"THE TRAGEDY OF ROMEO AND JULIET",
	"THE TAMING OF THE SHREW",
	"THE TEMPEST",
	"THE LIFE OF TIMON OF ATHENS",
	"THE TRAGEDY OF TITUS ANDRONICUS",
	"THE HISTORY OF TROILUS AND CRESSIDA",
	"TWELFTH NIGHT: OR, WHAT YOU WILL",
	"THE TWO GENTLEMEN OF VERONA",
	"THE TWO NOBLE KINSMEN",
	"THE WINTER’S TALE",
	"A LOVER’S COMPLAINT",
	"THE PASSIONATE PILGRIM",
	"THE PHOENIX AND THE TURTLE",
	"THE RAPE OF LUCRECE",
	"VENUS AND ADONIS"}

type Searcher struct {
	CompleteWorks string
	Works         []Work
}

type Work struct {
	Title       string
	Text        string
	SuffixArray *suffixarray.Index
}

type WorkResult struct {
	Index   int
	Title   string
	Matches []WorkMatch
}

type WorkMatch struct {
	Text  string
	Index int
}
	
func main() {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
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

func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.CompleteWorks = string(dat)

	s.loadSonnets()
	s.loadPlays()

	return nil
}

func (s *Searcher) Search(query string) []WorkResult {
	results := []WorkResult{}

	for workIndex, work := range s.Works {
		result := WorkResult{}
		result.Index = workIndex
		result.Title = work.Title

		for _, matchIndex := range work.SuffixArray.Lookup([]byte(query), -1) {
			match := WorkMatch{}
			match.Index = matchIndex;
			match.Text = s.getLine(matchIndex, work.Text)

			result.Matches = append(result.Matches, match)
		}

		if (len(result.Matches) > 0) {
			results = append(results, result)
		}
	}
	return results
}

func (s *Searcher) loadSonnets() {
	sonnetsStart := regexp.MustCompile(`\nTHE SONNETS`).FindStringIndex(s.CompleteWorks)[1]
	sonnetsEnd := regexp.MustCompile(`\nTHE END`).FindStringIndex(s.CompleteWorks)[0]
	sonnets := regexp.MustCompile(`(?:[^0-9\r\n]+\r\n)+`).FindAllString(s.CompleteWorks[sonnetsStart:sonnetsEnd], -1)
	for index, sonnet := range sonnets {
		work := Work{}
		work.Title = "Sonnet " + strconv.Itoa(index + 1)
		work.Text = sonnet
		work.SuffixArray = suffixarray.New([]byte(sonnet))
		s.Works = append(s.Works, work);
	}
}

func (s *Searcher) loadPlays() {
	playsStart := regexp.MustCompile(`\nTHE END`).FindStringIndex(s.CompleteWorks)[1]
	playsEnd := regexp.MustCompile(`CONTENT NOTE`).FindStringIndex(s.CompleteWorks)[0]

	for titleIndex, title := range playTitles {
		start := regexp.MustCompile(title).FindStringIndex(s.CompleteWorks[playsStart:playsEnd])[1]
		var end int
		if (titleIndex == len(playTitles) - 1) {
			end = playsEnd
		} else {
			end = regexp.MustCompile(playTitles[titleIndex + 1]).FindStringIndex(s.CompleteWorks[playsStart:playsEnd])[0]
		}
		text := s.CompleteWorks[start:end]

		work := Work{}
		work.Title = title
		work.Text = text
		work.SuffixArray = suffixarray.New([]byte(text))
		s.Works = append(s.Works, work);
	}
}

func (s *Searcher) getLine(index int, source string) string {
	var start, end = index, index

	for source[start] != '\n' && source[start] != '\r' {
		start--
	}
	for source[end] != '\n' && source[end] != '\r' {
		end++
	}

	return source[start:end]
}