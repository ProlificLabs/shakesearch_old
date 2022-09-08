package main

import (
	"bufio"
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
	"strings"
)

const (
	AMOUNT_OF_LINES int = 169432
	START_OF_WORKS  int = 133
)

var END_OF_EACH_WORK []int = []int{2909, 7873, 14514, 17306, 20507, 24624, 30499, 37187,
	41902, 45311, 50247, 53518, 57006, 60352, 64001, 66918, 71562, 77598, 80513,
	84662, 87663, 91831, 94789, 98249, 103961, 110231, 114369, 117467, 121875,
	127131, 132005, 135828, 138544, 141420, 147583, 152066, 154478, 159666,
	164681, 165064, 165304, 165398, 167583, 169017}

func main() {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handleSearch(searcher))
	http.HandleFunc("/add-lines/up", handleRequestAddLines(true, searcher))
	http.HandleFunc("/add-lines/down", handleRequestAddLines(false, searcher))
	http.HandleFunc("/read-work", handleReadWork(searcher))

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
	CompleteWorks    string
	SuffixArray      *suffixarray.Index
	SearchLines      []SearchLine
	EndOfLineIndexes []int
}

type SearchLine struct {
	TextResult string
	LineIndex  uint
}

func (s *SearchLine) setTextResult(newTextResult string) {
	s.TextResult = newTextResult
}

type ResultParagraph struct {
	Paragraph []SearchLine
}

func handleReadWork(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		lineNumberString, ok := r.URL.Query()["line"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Wrong request, make sure to add line-number in URL."))
			return
		}
		queryArray, ok2 := r.URL.Query()["q"]
		if !ok2 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Wrong request, make sure to add query in URL."))
			return
		}
		query := queryArray[0]
		lineNumber, error := strconv.Atoi(lineNumberString[0])
		if error != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("The lineNumber in the URL is not a number..."))
			return
		}
		if lineNumber < START_OF_WORKS {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("That paragraph is not part of a book. It belongs to the introduction of the complete works."))
			return
		}

		workId := getWorkIdOfLine(lineNumber)
		if workId == -1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("That paragraph is not part of a book. It belongs to the ending of the complete works. "))
			return
		}
		allLines := getLinesByWorkId(workId, searcher, query, lineNumber)

		results := []ResultParagraph{ResultParagraph{Paragraph: allLines}} // this is a slice because frontend expects an iterable
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

func getWorkIdOfLine(lineNumber int) int {
	var workId int = -1
	for i := 0; i < len(END_OF_EACH_WORK); i++ {
		if lineNumber < END_OF_EACH_WORK[i] {
			workId = i
			break
		}
	}
	return workId
}

func getLinesByWorkId(workId int, s Searcher, query string, lineNumber int) []SearchLine {
	allLines := []SearchLine{}
	beginning := START_OF_WORKS
	if workId != 0 {
		beginning = END_OF_EACH_WORK[workId-1] + 1
	}
	end := END_OF_EACH_WORK[workId]

	for i := beginning; i < end; i++ {
		currentLine := s.SearchLines[i].highlightRegexQuery(query)
		if i == lineNumber {
			currentLine = currentLine.addScrollId()
		}
		allLines = append(allLines, currentLine)
	}
	return allLines
}

func handleRequestAddLines(addLinesUp bool, searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r)
		lineNumberList, ok := r.URL.Query()["q"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Wrong request, make sure to add query in URL."))
			return
		}
		lineNumber, error := (strconv.Atoi(lineNumberList[0]))
		if error != nil {
			fmt.Println(error)
		}

		if lineNumber >= AMOUNT_OF_LINES || lineNumber <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Line number is higher than amount of lines in document..."))
			return
		}

		extraLines := getExtraLines(addLinesUp, lineNumber, searcher)

		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(extraLines)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
}

func getExtraLines(addLinesUp bool, lineNumber int, searcher Searcher) []SearchLine {
	extraLines := []SearchLine{}
	for i := 1; i <= 3; i++ {
		if addLinesUp && lineNumber-i >= 0 {
			extraLines = append(extraLines, searcher.SearchLines[lineNumber-i])
		} else if !addLinesUp && lineNumber+i <= AMOUNT_OF_LINES {
			extraLines = append(extraLines, searcher.SearchLines[lineNumber+i])
		}
	}
	return extraLines
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r)
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}
		results := searcher.SearchWithRegex(query[0])
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

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}

	lines, endOfLineIndexes := loadLinesInfo(filename)

	s.EndOfLineIndexes = endOfLineIndexes
	s.SearchLines = lines
	s.CompleteWorks = string(data)
	s.SuffixArray = suffixarray.New(data)
	return nil
}

func loadLinesInfo(filename string) ([]SearchLine, []int) {
	file, _ := os.Open(filename)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []SearchLine{}
	lineIndex := 0
	currentEndOfLineIndex := 0
	endOfLineIndexes := []int{}
	for scanner.Scan() {
		text := scanner.Text()

		currentEndOfLineIndex += len(text) + 2 // The + 2 is because the delimiter of the scanner (/n) is left out
		endOfLineIndexes = append(endOfLineIndexes, currentEndOfLineIndex)

		searchLine := SearchLine{
			TextResult: text,
			LineIndex:  uint(lineIndex),
		}
		lines = append(lines, searchLine)
		lineIndex++
	}

	return lines, endOfLineIndexes
}

func (s *Searcher) SearchWithRegex(query string) []ResultParagraph {
	re := regexp.MustCompile(query)
	restultIndeces := re.FindAllStringIndex(s.CompleteWorks, -1)
	allSearchLinesIndeces := getSearchLinesByIndex(restultIndeces, s)
	allParagraphs := makeParagraphsOutOfLines(allSearchLinesIndeces, s, query)
	return allParagraphs
}

func makeParagraphsOutOfLines(searchLinesIndeces []int, s *Searcher, query string) []ResultParagraph {
	allParagraphs := []ResultParagraph{}

	for _, lineIndex := range searchLinesIndeces {
		paragraph := []SearchLine{}
		if lineIndex != 0 {
			paragraph = append(paragraph, s.SearchLines[lineIndex-1].highlightRegexQuery(query))
		}

		paragraph = append(paragraph, s.SearchLines[lineIndex].highlightRegexQuery(query))

		if lineIndex != AMOUNT_OF_LINES-2 {
			paragraph = append(paragraph, s.SearchLines[lineIndex+1].highlightRegexQuery(query))
		}

		allParagraphs = append(allParagraphs, ResultParagraph{Paragraph: paragraph})
	}

	return allParagraphs
}

func getSearchLinesByIndex(indecesOfMatches [][]int, s *Searcher) []int {
	currentLine := 0
	results := []int{}

	for _, indexOfMatch := range indecesOfMatches {
		for currentLine < len(s.EndOfLineIndexes) {
			if indexOfMatch[0] < s.EndOfLineIndexes[currentLine] {
				results = append(results, currentLine)
				break
			} else {
				currentLine++
			}
		}
	}
	return results
}

func (s *Searcher) Search(query string) []ResultParagraph {
	results := []ResultParagraph{}

	for i, searchLine := range s.SearchLines {
		if strings.Contains(searchLine.TextResult, query) {
			if i == 0 {
				i = 1 // to prevent indexOutOfBoundsError
			}

			searchLines := []SearchLine{
				s.SearchLines[i-1].highlightQuery(query),
				s.SearchLines[i].highlightQuery(query),
				s.SearchLines[i+1].highlightQuery(query),
			}

			resultParagaph := ResultParagraph{
				Paragraph: searchLines,
			}

			results = append(results, resultParagaph)
		}
	}

	return results
}

func (s *SearchLine) highlightQuery(query string) SearchLine {
	newText := strings.ReplaceAll(s.TextResult, query, `<span style="color: #FD5F00;">`+query+`</span>`)
	s.setTextResult(newText)
	return *s
}

func (s *SearchLine) addScrollId() SearchLine {
	return SearchLine{
		TextResult: `<span id="scroll-here"></span>` + s.TextResult,
		LineIndex:  s.LineIndex,
	}
}

func (s *SearchLine) highlightRegexQuery(query string) SearchLine {
	re := regexp.MustCompile(query)
	newText := re.ReplaceAllStringFunc(s.TextResult,
		func(original string) string {
			return `<span style="color: #FD5F00;">` + original + `</span>`
		},
	)
	return SearchLine{
		TextResult: newText,
		LineIndex:  s.LineIndex,
	}
}
