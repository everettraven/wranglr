package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cli/cli/v2/pkg/search"
	"github.com/cli/go-gh/v2/pkg/auth"
)

type Client struct {
	host       string
	httpClient *http.Client
}

func NewClient(host string) *Client {
	return &Client{
		host:       host,
		httpClient: http.DefaultClient,
	}
}

func (c *Client) Issues(ctx context.Context, queries ...string) ([]search.Issue, error) {
	path := fmt.Sprintf("https://api.%s/search/issues", c.host)

	out := []search.Issue{}
	for _, query := range queries {
		qs := url.Values{}
		qs.Set("q", query)

		uri := fmt.Sprintf("%s?%s", path, qs.Encode())
		req, err := http.NewRequest(http.MethodGet, uri, nil)
		if err != nil {
			return nil, fmt.Errorf("building request: %w", err)
		}

		req.Header.Add("Content-Type", "application/json; charset=utf-8")
		req.Header.Add("Accept", "application/vnd.github.v3+json")

		if authToken := getAuthToken(c.host); authToken != "" {
			req.Header.Add("Authorization", fmt.Sprintf("token %s", authToken))
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("doing http request: %w", err)
		}

		// TODO: stream this?
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("reading response body: %w", err)
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, fmt.Errorf("request failed with status %q: %s", resp.Status, bodyBytes)
		}

		results := &search.IssuesResult{}

		err = json.Unmarshal(bodyBytes, results)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling results: %w", err)
		}

		out = append(out, results.Items...)

	}

	return out, nil
}

func getAuthToken(host string) string {
	// TODO: in the future, could continue falling
	// back to other methods in case
	// users do not use the GitHub CLI
	// for interacting with GitHub.
	// For now, this should be sufficiently robust.
	token, _ := auth.TokenForHost(host)
	return token
}
