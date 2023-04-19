package main

type SearchResult struct {
	Work      string `json:"work"`
	Text      string `json:"text"`
	Character string `json:"char"`
}

type SearchQuery struct {
	QueryText string   `json:"query"`
	WorkIds   []string `json:"workIds"`
	CharIds   []string `json:"charIds"`
}

type Work struct {
	WorkID string `json:"workId"`
	Title  string `json:"title"`
}

type Character struct {
	CharID string `json:"charId"`
	Name   string `json:"name"`
}

type HomePageData struct {
	Works      []Work
	Characters []Character
}
