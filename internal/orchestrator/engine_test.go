package orchestrator

import (
	"strings"
	"testing"

	"github.com/surtr85/miku-researcher/internal/search"
)

func TestFormatSearchResults(t *testing.T) {
	items := []search.ResultItem{
		{
			Title:   "Go Language",
			URL:     "https://golang.org",
			Snippet: "Go is an open source programming language.",
		},
		{
			Title:   "NixOS",
			URL:     "https://nixos.org",
			Snippet: "Declarative builds and deployments.",
		},
	}

	res := FormatSearchResults(items, "Quick summary of topics")
	if !strings.Contains(res, "Quick summary of topics") {
		t.Errorf("expected summary to be included")
	}
	if !strings.Contains(res, "https://golang.org") {
		t.Errorf("expected golang url")
	}
	if !strings.Contains(res, "NixOS") {
		t.Errorf("expected nixos title")
	}
}
