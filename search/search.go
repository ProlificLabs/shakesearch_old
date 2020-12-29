package search

import (
	"fmt"
	"github.com/sajari/fuzzy"
	"index/suffixarray"
	"io/ioutil"
	"strings"
)

// Scanner struct ...
type Scanner struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
	dictionary    *fuzzy.Model
}

// Scanner constructor ...
func NewScanner() *Scanner {
	return &Scanner{}
}

// training model to track words in a dictionary for misspell correction ...
func (s *Scanner) trainDictionary(words []string) {
	fmt.Println("Training Dictionary for spell correction ...")
	s.dictionary = fuzzy.NewModel()
	s.dictionary.SetThreshold(1)
	s.dictionary.SetDepth(2)
	s.dictionary.Train(words)
}

// Load the text file & save indexes ...
func (s *Scanner) Load(filename string) error {
	fmt.Println("Loading CompleteWorks ...")
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("load err: %w", err)
	}
	str := string(dat)
	s.CompleteWorks = str
	fmt.Println("Creating Indexes ...")
	s.SuffixArray = suffixarray.New([]byte(strings.ToLower(str)))
	s.trainDictionary(strings.Split(str, " "))
	return nil
}

// Search for a given query ...
func (s *Scanner) Search(query string) ([]string, string) {
	var correction string
	idxes := s.SuffixArray.Lookup([]byte(strings.ToLower(query)), -1)
	var result []string
	for _, idx := range idxes {
		raw := s.CompleteWorks[idx-250:idx+250]
		initialStr := trimSentence("start", raw)
		resultStr := trimSentence("end", initialStr)
		result = append(result, resultStr)
	}
	if len(result) == 0 {
		correction = s.dictionary.SpellCheck(query)
	}
	return result, correction
}

// trim the given content based on full stops ...
func trimSentence(from, str string) string {
	var (
		resultStr = ""
		splitAt = -1
	)
	if str == "" {
		return str
	}
	if from == "start" {
		splitAt = strings.Index(str, ".")
		if splitAt == -1 {
			splitAt = 0
		}
		splitAt = splitAt + 1
		resultStr = str[splitAt:]
	} else {
		splitAt = strings.LastIndex(str, ".")
		if splitAt == -1 {
			splitAt = 0
		}
		splitAt = splitAt + 1
		resultStr = str[:splitAt]
	}
	return resultStr
}