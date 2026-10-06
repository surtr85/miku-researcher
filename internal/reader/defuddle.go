package reader

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type PageResult struct {
	URL     string `json:"url"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content"`
	Engine  string `json:"engine,omitempty"`
	Error   string `json:"error,omitempty"`
}

type Reader struct {
	defuddlePath string
	firecrawlKey string
	timeout      time.Duration
	httpClient   *http.Client
}

func NewReader(timeout time.Duration) *Reader {
	p, err := exec.LookPath("defuddle")
	if err != nil {
		p = "defuddle"
	}
	if timeout <= 0 {
		timeout = 25 * time.Second
	}
	firecrawlKey := os.Getenv("FIRECRAWL_API_KEY")
	if firecrawlKey == "" {
		firecrawlKey = "fc-deb0899c43da4021b891848cd6fdd372"
	}

	return &Reader{
		defuddlePath: p,
		firecrawlKey: firecrawlKey,
		timeout:      timeout,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// ReadPage fetches and extracts clean markdown using defuddle, falling back to Firecrawl if defuddle fails
func (r *Reader) ReadPage(ctx context.Context, targetURL string, maxChars int) (*PageResult, error) {
	ctxDefuddle, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	// 1. Try defuddle parse --markdown <url>
	cmd := exec.CommandContext(ctxDefuddle, r.defuddlePath, "parse", "--markdown", targetURL)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	content := strings.TrimSpace(stdout.String())

	// Check if defuddle succeeded and extracted meaningful content (> 150 chars)
	if err == nil && len(content) >= 150 {
		if maxChars > 0 && len(content) > maxChars {
			content = content[:maxChars] + "\n\n... [Content truncated for token efficiency] ..."
		}
		return &PageResult{
			URL:     targetURL,
			Content: content,
			Engine:  "defuddle",
		}, nil
	}

	// 2. If defuddle failed or returned too little content, try Firecrawl fallback
	if r.firecrawlKey != "" {
		res, errFc := r.ScrapeWithFirecrawl(ctx, targetURL, maxChars, true)
		if errFc == nil && res != nil && res.Error == "" && len(res.Content) > 0 {
			return res, nil
		}
	}

	// If defuddle had some content even if short, return it
	if err == nil && len(content) > 0 {
		return &PageResult{
			URL:     targetURL,
			Content: content,
			Engine:  "defuddle",
		}, nil
	}

	errMsg := strings.TrimSpace(stderr.String())
	if errMsg == "" && err != nil {
		errMsg = err.Error()
	}
	return &PageResult{
		URL:   targetURL,
		Error: fmt.Sprintf("extraction error: %s", errMsg),
	}, nil
}

// ScrapeWithFirecrawl scrapes via Firecrawl v1 API
func (r *Reader) ScrapeWithFirecrawl(ctx context.Context, targetURL string, maxChars int, onlyMainContent bool) (*PageResult, error) {
	if r.firecrawlKey == "" {
		return &PageResult{
			URL:   targetURL,
			Error: "FIRECRAWL_API_KEY is not configured",
		}, nil
	}

	reqBody, err := json.Marshal(map[string]any{
		"url":             targetURL,
		"formats":         []string{"markdown"},
		"onlyMainContent": onlyMainContent,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.firecrawl.dev/v1/scrape", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+r.firecrawlKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return &PageResult{
			URL:   targetURL,
			Error: fmt.Sprintf("firecrawl request failed: %v", err),
		}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return &PageResult{
			URL:   targetURL,
			Error: fmt.Sprintf("firecrawl returned HTTP %d: %s", resp.StatusCode, string(bodyBytes)),
		}, nil
	}

	var fcResp struct {
		Success bool `json:"success"`
		Data    struct {
			Markdown string `json:"markdown"`
			Metadata struct {
				Title string `json:"title"`
			} `json:"metadata"`
		} `json:"data"`
		Error string `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&fcResp); err != nil {
		return nil, err
	}

	if !fcResp.Success {
		return &PageResult{
			URL:   targetURL,
			Error: fmt.Sprintf("firecrawl error: %s", fcResp.Error),
		}, nil
	}

	content := strings.TrimSpace(fcResp.Data.Markdown)
	if maxChars > 0 && len(content) > maxChars {
		content = content[:maxChars] + "\n\n... [Content truncated for token efficiency] ..."
	}

	return &PageResult{
		URL:     targetURL,
		Title:   fcResp.Data.Metadata.Title,
		Content: content,
		Engine:  "firecrawl",
	}, nil
}
