package semanticscholar

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchAuthors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}

			if r.URL.Path != "/author/search" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}

			if got := r.URL.Query().Get("query"); got != "Geoffrey Hinton" {
				t.Errorf("unexpected query: %s", got)
			}

			if got := r.URL.Query().Get("limit"); got != "10" {
				t.Errorf("unexpected limit: %s", got)
			}

			if got := r.URL.Query().Get("offset"); got != "0" {
				t.Errorf("unexpected offset: %s", got)
			}

			if got := r.URL.Query().Get("fields"); got != "authorId,name" {
				t.Errorf("unexpected fields: %s", got)
			}

			w.Header().Set("Content-Type", "application/json")

			json.NewEncoder(w).Encode(AuthorSearchResponse{
				Total:  1,
				Offset: 0,
				Data: []Author{
					{
						AuthorID: "12345",
						Name:     "Geoffrey Hinton",
					},
				},
			})
		},
	))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))

	result, err := client.SearchAuthors(
		context.Background(),
		"Geoffrey Hinton",
		10,
		0,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Total != 1 {
		t.Errorf("expected total 1, got %d", result.Total)
	}

	if len(result.Data) != 1 {
		t.Fatalf("expected 1 author, got %d", len(result.Data))
	}

	if result.Data[0].Name != "Geoffrey Hinton" {
		t.Errorf("unexpected author name: %s", result.Data[0].Name)
	}
}

func TestSearchAuthorsInvalidInput(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		limit  int
		offset int
	}{
		{
			name:   "empty query",
			query:  "",
			limit:  10,
			offset: 0,
		},
		{
			name:   "whitespace query",
			query:  "   ",
			limit:  10,
			offset: 0,
		},
		{
			name:   "zero limit",
			query:  "Geoffrey Hinton",
			limit:  0,
			offset: 0,
		},
		{
			name:   "negative limit",
			query:  "Geoffrey Hinton",
			limit:  -1,
			offset: 0,
		},
		{
			name:   "limit above maximum",
			query:  "Geoffrey Hinton",
			limit:  101,
			offset: 0,
		},
		{
			name:   "negative offset",
			query:  "Geoffrey Hinton",
			limit:  10,
			offset: -1,
		},
	}

	client := NewClient()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.SearchAuthors(
				context.Background(),
				tt.query,
				tt.limit,
				tt.offset,
			)

			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
		})
	}
}

func TestGetAuthor(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}

			if r.URL.Path != "/author/123" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}

			expectedFields := "authorId,name,url,paperCount,citationCount,hIndex"
			if got := r.URL.Query().Get("fields"); got != expectedFields {
				t.Errorf("expected fields %q, got %q", expectedFields, got)
			}

			response := Author{
				AuthorID:      "123",
				Name:          "John Doe",
				URL:           "https://example.com/author/123",
				PaperCount:    50,
				CitationCount: 1200,
				HIndex:        15,
			}

			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
		},
	))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))

	author, err := client.GetAuthor(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if author.AuthorID != "123" {
		t.Errorf("unexpected author ID: %s", author.AuthorID)
	}

	if author.Name != "John Doe" {
		t.Errorf("unexpected author name: %s", author.Name)
	}

	if author.PaperCount != 50 {
		t.Errorf("expected 50 papers, got %d", author.PaperCount)
	}

	if author.CitationCount != 1200 {
		t.Errorf("expected 1200 citations, got %d", author.CitationCount)
	}

	if author.HIndex != 15 {
		t.Errorf("expected h-index 15, got %d", author.HIndex)
	}
}

func TestGetAuthorEmptyID(t *testing.T) {
	client := NewClient()

	_, err := client.GetAuthor(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty author ID")
	}
}

func TestGetAuthorNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "author not found", http.StatusNotFound)
		},
	))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))

	_, err := client.GetAuthor(context.Background(), "unknown")
	if err == nil {
		t.Fatal("expected 404 error")
	}
}

func TestGetAuthorInvalidFields(t *testing.T) {
	client := NewClient()

	_, err := client.GetAuthor(
		context.Background(),
		"123",
		"invalidField",
	)

	if err == nil {
		t.Fatal("expected error for invalid field")
	}
}
