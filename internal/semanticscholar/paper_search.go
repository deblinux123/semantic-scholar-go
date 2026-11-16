package semanticscholar

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

var allowedPaperFields = map[string]struct{}{
	"paperId":       {},
	"title":         {},
	"abstract":      {},
	"year":          {},
	"url":           {},
	"venue":         {},
	"citationCount": {},
	"authors":       {},
}

func (c *Client) SearchPapers(
	ctx context.Context,
	query string,
	limit int,
	offset int,
	fields ...string,
) (*PaperSearcgResponse, error) {
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty.")
	}

	if limit <= 0 || limit >= 100 {
		return nil, fmt.Errorf("limit must be between 1 and 100")
	}

	if offset < 0 {
		return nil, fmt.Errorf("offset cannot be negative")
	}

	defaultFields := []string{
		"title",
		"abstract",
		"year",
		"url",
		"venue",
		"citationCount",
		"authors",
	}

	if len(fields) == 0 {
		fields = defaultFields
	}

	if err := validatePaperFields(fields); err != nil {
		return nil, err
	}

	params := url.Values{}

	params.Set("query", query)
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(offset))
	params.Set("fields", strings.Join(fields, ","))

	var result PaperSearcgResponse

	err := c.get(
		ctx,
		"/paper/search",
		params,
		&result,
	)

	if err != nil {
		return nil, fmt.Errorf("search papers: %w", err)
	}

	return &result, nil
}

func validatePaperFields(fileds []string) error {
	for _, filed := range fileds {
		if _, ok := allowedPaperFields[filed]; !ok {
			return fmt.Errorf(
				"invalid paper field: %s",
				filed,
			)
		}
	}

	return nil
}
