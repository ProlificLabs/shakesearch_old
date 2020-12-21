package main

import (
	"net/http"
	"strconv"
	"testing"
)

func Test_SearchRequest_Bind(t *testing.T) {
	req, err := http.NewRequest("GET", "/search", nil)
	if err != nil {
		t.Fatal(err)
	}

	query := `"Pulley"`
	sensitive := true
	before := 10
	after := 20

	q := req.URL.Query()
	q.Add("q", query)
	q.Add("sensitive", strconv.FormatBool(sensitive))
	q.Add("before", strconv.Itoa(before))
	q.Add("after", strconv.Itoa(after))
	req.URL.RawQuery = q.Encode()

	searchRequest := &SearchRequest{}
	err = searchRequest.Bind(req)
	if err != nil {
		t.Fatal(err)
	}

	if searchRequest.Q != "Pulley" {
		t.Errorf("got %s, want %s", searchRequest.Q, query)
	}

	if !searchRequest.ExactMatch {
		t.Errorf("want ExactMatch but is false")
	}

	if !searchRequest.Sensitive {
		t.Errorf("want Sensitive but is false")
	}

	if searchRequest.After != after {
		t.Errorf("after: got %d, want %d", searchRequest.After, after)
	}

	if searchRequest.Before != before {
		t.Errorf("before: got %d, want %d", searchRequest.Before, before)
	}
}
