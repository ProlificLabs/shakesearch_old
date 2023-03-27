package main

import (
	"testing"
)

// Check the matches for each token in the snippet
func TestGetMatches(t *testing.T) {
	tokens := []string{"queen", "king", "eggs"}
	snippet := "The queen is making ham and eggs to all the kings."

	// Exact match
	tokenType := "exact"
	expectedMatches := []string{"queen", "eggs"}
	// Get matches in snippet
	matches := getMatches(tokens, tokenType, snippet)
	checkExpected(matches, expectedMatches, t)

	// regular match
	tokenType = "simple"
	expectedMatches = []string{"queen", "king", "eggs", "king"}
	// Get matches in snippet
	matches = getMatches(tokens, tokenType, snippet)
	checkExpected(matches, expectedMatches, t)
}

// Check the regex from a list of tokens
func TestGenerateRegex(t *testing.T) {
	exactTokens := []string{"king.", "alive$", "hamlet", "of denamark"}
	exactRegex := generateRegex(exactTokens, "exact")
	exactMatch := (exactRegex.String() == ("\\b(king\\.|alive\\$|hamlet|of denamark)\\b"))
	if !exactMatch {
		t.Errorf("Expected exact match, but did not find one")
	}

	simpleTokens := []string{"king.", "alive$", "hamlet", "of denamark"}
	simpleRegex := generateRegex(simpleTokens, "simple")
	simpleMatch := simpleRegex.String() == ("(king\\.|alive\\$|hamlet|of denamark)")
	if !simpleMatch {
		t.Errorf("Expected simple match, but did not find one")
	}
}

func checkExpected(matches []string, expectedMatches []string, t *testing.T) {
	if len(matches) != len(expectedMatches) {
		t.Errorf("Expected %d matches, but got %d", len(expectedMatches), len(matches))
	}
	for i, match := range matches {
		if match != expectedMatches[i] {
			t.Errorf("Expected match '%s', but got '%s'", expectedMatches[i], match)
		}
	}
}
