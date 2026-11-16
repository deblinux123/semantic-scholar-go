package semanticscholar

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchPapers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}

			if r.URL.Path != "/graph/v1/paper/search" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}

			if r.URL.Query().Get("query") != "golang" {
				t.Errorf("unexpected query: %s", r.URL.Query().Get("query"))
			}

			if r.URL.Query().Get("limit") != "10" {
				t.Errorf("unexpected llimit")
			}

			if got := r.Header.Get("x-api-key"); got != "" {
				t.Errorf("expected empty API key, got %q", got)
			}

			response := PaperSearchResponse{
				Total:  1,
				Offset: 0,
				Data: []Paper{
					{
						PaperID: "paper-123",
						Title:   "Go Programming",
						Year:    2025,
					},
				},
			}

			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Errorf("failed to encode response: %v", err)
			}
		},
	))

	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL + "/graph/v1"),
	)

	result, err := client.SearchPapers(
		context.Background(),
		"golang",
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
		t.Fatalf("expected 1 paper, got %d", len(result.Data))
	}

	paper := result.Data[0]

	if paper.PaperID != "paper-123" {
		t.Errorf("unexpected paper ID: %s", paper.PaperID)
	}

	if paper.Title != "Go Programming" {
		t.Errorf("unexpected title: %s", paper.Title)
	}
}

func TestSearchPapersEmptyQuery(t *testing.T) {
	client := NewClient()

	result, err := client.SearchPapers(
		context.Background(),
		"",
		3,
		0,
	)

	if err == nil {
		t.Fatal("expected error for empty query")
	}

	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestSearchPapersPagenation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query()

			if got := query.Get("limit"); got != "20" {
				t.Errorf("expected limit 20, got %s", got)
			}

			if got := query.Get("offset"); got != "40" {
				t.Errorf("expected offset 40, got %s", got)
			}

			response := PaperSearchResponse{
				Total:  100,
				Offset: 40,
				Data:   []Paper{},
			}

			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			json.NewEncoder(w).Encode(response)
		},
	))

	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL + "/graph/v1"),
	)

	result, err := client.SearchPapers(
		context.Background(),
		"golang",
		20,
		40,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Offset != 40 {
		t.Errorf("expected offset 40, got %d", result.Offset)
	}

	if result.Total != 100 {
		t.Errorf("expected total 100, got %d", result.Total)
	}
}

func TestSearchPapersInvalidPagination(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name   string
		limit  int
		offset int
	}{
		{"zero limit", 0, 0},
		{"negative limit", -1, 0},
		{"limit above maximum", 101, 0},
		{"negative offset", 10, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := client.SearchPapers(
				context.Background(),
				"golang",
				tt.limit,
				tt.offset,
			)

			if err == nil {
				t.Fatal("expected validation error")
			}

			if result != nil {
				t.Fatal("expected nil result")
			}
		})
	}
}

func TestSearchPapersCustomFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fields := r.URL.Query().Get("fields")

			expected := "title,year,authors"

			if fields != expected {
				t.Errorf(
					"expected fields %q, got %q",
					expected,
					fields,
				)
			}

			response := PaperSearchResponse{
				Total: 1,
				Data: []Paper{
					{
						PaperID: "paper-123",
						Title:   "Go Programming",
						Year:    2025,
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		},
	))

	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL + "/graph/v1"),
	)

	result, err := client.SearchPapers(
		context.Background(),
		"golang",
		10,
		0,
		"title",
		"year",
		"authors",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Data) != 1 {
		t.Fatalf("expected 1 paper, got %d", len(result.Data))
	}
}

func TestValidatePaperFields(t *testing.T) {
	tests := []struct {
		name    string
		fields  []string
		wantErr bool
	}{
		{
			name:    "valid fields",
			fields:  []string{"title", "year", "authors"},
			wantErr: false,
		},
		{
			name:    "invalid field",
			fields:  []string{"title", "invalidField"},
			wantErr: true,
		},
		{
			name:    "empty fields",
			fields:  []string{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePaperFields(tt.fields)

			if (err != nil) != tt.wantErr {
				t.Errorf(
					"validatePaperFields() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestGetPaper(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}

			if r.URL.Path != "/paper/paper-123" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}

			if got := r.URL.Query().Get("fields"); got != "title,year,abstract" {
				t.Errorf("unexpected fields: %s", got)
			}

			w.Header().Set("Content-Type", "application/json")

			json.NewEncoder(w).Encode(Paper{
				PaperID: "paper-123",
				Title:   "Learning Go",
				Year:    2025,
			})
		},
	))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))

	paper, err := client.GetPaper(
		context.Background(),
		"paper-123",
		"title",
		"year",
		"abstract",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if paper.PaperID != "paper-123" {
		t.Errorf("unexpected paper ID: %s", paper.PaperID)
	}

	if paper.Title != "Learning Go" {
		t.Errorf("unexpected title: %s", paper.Title)
	}

	if paper.Year != 2025 {
		t.Errorf("unexpected year: %d", paper.Year)
	}
}

func TestGetPaperEmptyID(t *testing.T) {
	client := NewClient()

	_, err := client.GetPaper(
		context.Background(),
		"   ",
	)

	if err == nil {
		t.Fatal("expected error for empty paper ID")
	}
}

func TestGetPaperNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(
				w,
				`{"error":"Paper not found"}`,
				http.StatusNotFound,
			)
		},
	))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))

	_, err := client.GetPaper(
		context.Background(),
		"unknown-paper",
	)

	if err == nil {
		t.Fatal("expected 404 error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}

	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf(
			"expected status 404, got %d",
			apiErr.StatusCode,
		)
	}
}
