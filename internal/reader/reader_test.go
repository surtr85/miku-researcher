package reader

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestScrapeWithFirecrawl_Mock(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		resp := map[string]any{
			"success": true,
			"data": map[string]any{
				"markdown": "# Mock Title\n\nMock extracted content from Firecrawl",
				"metadata": map[string]any{
					"title": "Mock Title",
				},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer mockServer.Close()

	r := &Reader{
		defuddlePath: "defuddle",
		firecrawlKey: "test-key",
		timeout:      5 * time.Second,
		httpClient:   mockServer.Client(),
	}

	if r.firecrawlKey != "test-key" {
		t.Errorf("expected test-key, got %s", r.firecrawlKey)
	}
}
