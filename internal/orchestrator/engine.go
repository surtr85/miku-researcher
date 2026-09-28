package orchestrator

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/surtr85/miku-researcher/internal/reader"
	"github.com/surtr85/miku-researcher/internal/search"
)

type ResearchEngine struct {
	searchClient *search.Client
	pageReader   *reader.Reader
}

func NewResearchEngine() *ResearchEngine {
	return &ResearchEngine{
		searchClient: search.NewClient(),
		pageReader:   reader.NewReader(25 * time.Second),
	}
}

// DeepResearch performs autonomous multi-stage search, concurrent reading, and distillation
func (e *ResearchEngine) DeepResearch(ctx context.Context, topic string, maxSources int, maxCharsPerSource int) (string, error) {
	if maxSources <= 0 {
		maxSources = 3
	}
	if maxSources > 7 {
		maxSources = 7
	}
	if maxCharsPerSource <= 0 {
		maxCharsPerSource = 3000
	}

	// 1. Search for best sources using SearXNG or Tavily
	var sources []search.ResultItem
	var quickAnswer string
	var err error

	// Try tavily first if available for high relevance, otherwise searxng
	sources, quickAnswer, err = e.searchClient.TavilySearch(ctx, topic, "advanced", maxSources)
	if err != nil || len(sources) == 0 {
		sources, err = e.searchClient.SearxngSearch(ctx, topic, "general,it,science", "", "", maxSources)
		if err != nil {
			return "", fmt.Errorf("search failed: %w", err)
		}
	}

	if len(sources) == 0 {
		return "No relevant sources found for this topic.", nil
	}

	// 2. Concurrently read and extract markdown using Defuddle
	type fetchedDoc struct {
		title   string
		url     string
		content string
		err     string
	}

	results := make([]fetchedDoc, len(sources))
	var wg sync.WaitGroup

	for i, s := range sources {
		wg.Add(1)
		go func(idx int, item search.ResultItem) {
			defer wg.Done()
			page, err := e.pageReader.ReadPage(ctx, item.URL, maxCharsPerSource)
			if err != nil || (page != nil && page.Error != "") {
				errMsg := ""
				if err != nil {
					errMsg = err.Error()
				} else {
					errMsg = page.Error
				}
				results[idx] = fetchedDoc{
					title:   item.Title,
					url:     item.URL,
					content: item.Snippet, // fallback to search snippet
					err:     errMsg,
				}
			} else {
				results[idx] = fetchedDoc{
					title:   item.Title,
					url:     item.URL,
					content: page.Content,
				}
			}
		}(i, s)
	}

	wg.Wait()

	// 3. Assemble high-density Markdown synthesis
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Research Dossier: %s\n\n", topic))

	if quickAnswer != "" {
		sb.WriteString("## Executive Summary\n")
		sb.WriteString(quickAnswer)
		sb.WriteString("\n\n")
	}

	sb.WriteString(fmt.Sprintf("## Extracted Intelligence (%d Sources Analyzed)\n\n", len(results)))

	for i, doc := range results {
		sb.WriteString(fmt.Sprintf("### Source %d: [%s](%s)\n", i+1, doc.title, doc.url))
		if doc.err != "" {
			sb.WriteString(fmt.Sprintf("> *Note: Full extraction degraded (%s). Utilizing indexed summary.*\n\n", doc.err))
		}
		sb.WriteString(doc.content)
		sb.WriteString("\n\n---\n\n")
	}

	return strings.TrimSpace(sb.String()), nil
}

// FormatSearchResults returns high-density Markdown representation instead of wasteful JSON
func FormatSearchResults(items []search.ResultItem, answer string) string {
	var sb strings.Builder
	if answer != "" {
		sb.WriteString("### Quick Answer\n")
		sb.WriteString(answer)
		sb.WriteString("\n\n")
	}
	sb.WriteString("### Search Results\n")
	for i, item := range items {
		sb.WriteString(fmt.Sprintf("%d. **[%s](%s)**\n", i+1, item.Title, item.URL))
		if item.Snippet != "" {
			sb.WriteString(fmt.Sprintf("   > %s\n", strings.ReplaceAll(item.Snippet, "\n", " ")))
		}
	}
	return strings.TrimSpace(sb.String())
}
