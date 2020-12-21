package searcher

type Searcher interface {
	Search(req Request) (*Response, error)
}

type Request struct {
	Query            string
	CaseSensitive    bool
	ExactMatch       bool
	CharBeforeQuery  int
	CharAfterQuery   int
	HighlightPreTag  string
	HighlightPostTag string
}

type Response struct {
	Query      string       `json:"query"`
	Hits       []*Hit       `json:"hits"`
	Highlights []*Highlight `json:"highlights"`
}

type Highlight struct {
	Text         string   `json:"text"`
	MatchedWords []string `json:"matched_words"`
}

type Hit struct {
	Query      string `json:"query"`
	Order      int    `json:"order"`
	Start      int    `json:"start"`
	End        int    `json:"end"`
	ExactMatch bool   `json:"exact_match"`
}

type Hits []*Hit

func (h Hits) Len() int           { return len(h) }
func (h Hits) Less(i, j int) bool { return h[i].Start < h[j].Start }
func (h Hits) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].Order = i + 1
	h[j].Order = j + 1
}
