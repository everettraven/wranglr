package jira

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	gojira "github.com/andygrunwald/go-jira"
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

func (c *Client) Issues(queries ...string) ([]gojira.Issue, error) {
	jiraClient, err := gojira.NewClient(nil, c.host)
	if err != nil {
		return nil, err
	}

	results := []gojira.Issue{}

	for _, query := range queries {
		u := url.URL{
			Path: "rest/api/latest/search",
		}
		uv := url.Values{}
		uv.Add("jql", query)
		u.RawQuery = uv.Encode()

		req, _ := jiraClient.NewRequest("GET", u.String(), nil)

		if token := getAuth(); token != "" {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		}

		searchRes := &searchResult{}
		resp, err := jiraClient.Do(req, searchRes)
		if err != nil {
			return nil, fmt.Errorf("fetching jira issues using query %q: %w", query, err)
		}
		defer func() { _ = resp.Body.Close() }()

		results = append(results, searchRes.Issues...)

	}

	return results, nil
}

// TODO: implement some kind of multi-host authentication
// reading strategy.
// For now, assume someone will only be authenticating
// against one instance at a time.
func getAuth() string {
	auth := os.Getenv("WRANGLR_JIRA_TOKEN")
	return auth
}
