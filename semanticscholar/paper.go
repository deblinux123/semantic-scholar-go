package semanticscholar

type Paper struct {
	PaperID       string   `json:"paperId"`
	Title         string   `json:"title"`
	Abstract      string   `json:"abstract,omitempty"`
	Year          int      `json:"year"`
	URL           string   `json:"url"`
	Venue         string   `json:"venue,omitempty"`
	CitationCount int      `json:"citationCount,omitempty"`
	Authors       []Author `json:"author"`
}
type PaperSearchResponse struct {
	Total  int     `json:"total"`
	Offset int     `json:"offset"`
	Next   int     `json:"nxt,omitempty"`
	Data   []Paper `json:"data"`
}

type Author struct {
	AuthorID      string `json:"authorId"`
	Name          string `json:"name"`
	URL           string `json:"url"`
	PaperCount    int    `json:"paperCount"`
	CitationCount int    `json:"citationCount"`
	HIndex        int    `json:"hIndex"`
}

type AuthorSearchResponse struct {
	Total  int      `json:"total"`
	Offset int      `json:"offset"`
	Next   *int     `json:"next"`
	Data   []Author `json:"data"`
}
