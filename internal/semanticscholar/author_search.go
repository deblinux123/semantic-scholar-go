package semanticscholar

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

var allowedAuthorfields = map[string]struct{}{
	"authorId":      {},
	"name":          {},
	"url":           {},
	"paperCount":    {},
	"citationCount": {},
	"hIndex":        {},
}

func (c *Client) SearchAuthors(
	ctx context.Context,
	query string,
	limit int,
	offset int,
) (*AuthorSearchResponse, error) {
	if strings.TrimSpace(query) == "" {
		return nil, fmt.Errorf("author query cannot be empty")
	}

	if limit < 1 || limit > 100 {
		return nil, fmt.Errorf("limit must be between 1 and 100")
	}

	if offset < 0 {
		return nil, fmt.Errorf("offset cannot be negative")
	}

	params := url.Values{}
	params.Set("query", query)
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(offset))
	params.Set("fields", "authorId,name")

	var result AuthorSearchResponse

	if err := c.get(
		ctx,
		"/author/search",
		params,
		&result,
	); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetAuthor(
	ctx context.Context,
	authorId string,
	fields ...string,
) (*Author, error) {
	if strings.TrimSpace(authorId) == "" {
		return nil, fmt.Errorf("author ID cannot be empty")
	}

	if len(fields) == 0 {
		fields = []string{
			"authorId",
			"name",
			"url",
			"paperCount",
			"citationCount",
			"hIndex",
		}
	}

	if err := validateAuthorFields(fields); err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("fields", strings.Join(fields, ","))

	path := "/author/" + url.PathEscape(authorId)

	var author Author

	if err := c.get(ctx, path, params, &author); err != nil {
		return nil, err
	}

	return &author, nil
}

func validateAuthorFields(fields []string) error {
	for _, field := range fields {
		if _, ok := allowedAuthorfields[field]; !ok {
			return fmt.Errorf("invalid author field: %s", field)
		}
	}

	return nil
}
