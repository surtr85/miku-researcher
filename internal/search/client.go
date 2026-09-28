package search

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type ResultItem struct {
	Title         string `json:"title"`
	URL           string `json:"url"`
	Snippet       string `json:"snippet"`
	Engine        string `json:"engine,omitempty"`
	PublishedDate string `json:"published_date,omitempty"`
}

type Client struct {
	searxngURL string
	tavilyKey  string
	httpClient *http.Client
}

func NewClient() *Client {
	searxng := os.Getenv("SEARXNG_URL")
	if searxng == "" {
		searxng = "http://localhost:8080"
	}
	tavily := os.Getenv("TAVILY_API_KEY")

	return &Client{
		searxngURL: strings.TrimRight(searxng, "/"),
		tavilyKey:  tavily,
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

// SearxngSearch queries the SearXNG instance
func (c *Client) SearxngSearch(ctx context.Context, query string, categories, timeRange, language string, limit int) ([]ResultItem, error) {
	u, err := url.Parse(c.searxngURL + "/search")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("q", query)
	q.Set("format", "json")
	if categories != "" {
		q.Set("categories", categories)
	}
	if timeRange != "" {
		q.Set("time_range", timeRange)
	}
	if language != "" {
		q.Set("language", language)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "MikuResearcher/1.0 (Autonomous; Pure-Go)")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("searxng returned HTTP %d", resp.StatusCode)
	}

	var data struct {
		Results []struct {
			Title         string `json:"title"`
			URL           string `json:"url"`
			Content       string `json:"content"`
			Engine        string `json:"engine"`
			PublishedDate string `json:"publishedDate"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 10
	}

	var items []ResultItem
	for i, r := range data.Results {
		if i >= limit {
			break
		}
		items = append(items, ResultItem{
			Title:         r.Title,
			URL:           r.URL,
			Snippet:       r.Content,
			Engine:        r.Engine,
			PublishedDate: r.PublishedDate,
		})
	}

	return items, nil
}

// TavilySearch queries Tavily API if configured
func (c *Client) TavilySearch(ctx context.Context, query string, depth string, maxResults int) ([]ResultItem, string, error) {
	if c.tavilyKey == "" {
		return nil, "", fmt.Errorf("TAVILY_API_KEY is not configured")
	}

	if depth == "" {
		depth = "basic"
	}
	if maxResults <= 0 {
		maxResults = 5
	}

	reqBody, _ := json.Marshal(map[string]any{
		"api_key":        c.tavilyKey,
		"query":          query,
		"search_depth":   depth,
		"max_results":    maxResults,
		"include_answer": true,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.tavily.com/search", strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("tavily returned HTTP %d", resp.StatusCode)
	}

	var data struct {
		Answer  string `json:"answer"`
		Results []struct {
			Title   string `json:"title"`
			URL     string `json:"url"`
			Content string `json:"content"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, "", err
	}

	var items []ResultItem
	for _, r := range data.Results {
		items = append(items, ResultItem{
			Title:   r.Title,
			URL:     r.URL,
			Snippet: r.Content,
			Engine:  "tavily",
		})
	}

	return items, data.Answer, nil
}
