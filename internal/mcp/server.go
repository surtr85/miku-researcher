package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/surtr85/miku-researcher/internal/orchestrator"
	"github.com/surtr85/miku-researcher/internal/reader"
	"github.com/surtr85/miku-researcher/internal/search"
)

type Server struct {
	searchClient *search.Client
	pageReader   *reader.Reader
	engine       *orchestrator.ResearchEngine
	in           io.Reader
	out          io.Writer
}

func NewServer(in io.Reader, out io.Writer) *Server {
	return &Server{
		searchClient: search.NewClient(),
		pageReader:   reader.NewReader(0),
		engine:       orchestrator.NewResearchEngine(),
		in:           in,
		out:          out,
	}
}

func (s *Server) GetTools() []Tool {
	return []Tool{
		{
			Name:        "deep_research",
			Description: "Perform autonomous multi-source deep research on any topic or query. Searches top sources, concurrently reads & cleans pages with Defuddle, and synthesizes a high-density intelligence dossier.",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]PropertyDef{
					"topic": {
						Type:        "string",
						Description: "The topic, technical question, or subject to research deeply.",
					},
					"max_sources": {
						Type:        "integer",
						Description: "Number of primary sources to read and analyze (default: 4, max: 7).",
						Default:     4,
					},
					"max_chars_per_source": {
						Type:        "integer",
						Description: "Max character budget per source for token efficiency (default: 3500).",
						Default:     3500,
					},
				},
				Required: []string{"topic"},
			},
		},
		{
			Name:        "read_page",
			Description: "Extract and convert clean markdown from any URL using Defuddle reader engine (stripping ads, navbars, and boilerplate).",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]PropertyDef{
					"url": {
						Type:        "string",
						Description: "Target webpage or article URL to fetch and clean.",
					},
					"max_chars": {
						Type:        "integer",
						Description: "Maximum characters to return to conserve tokens (default: 6000, 0 = unlimited).",
						Default:     6000,
					},
				},
				Required: []string{"url"},
			},
		},
		{
			Name:        "searxng_search",
			Description: "Search the web via self-hosted SearXNG instance without third-party rate limits. Returns high-density markdown.",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]PropertyDef{
					"query": {
						Type:        "string",
						Description: "The search query string.",
					},
					"categories": {
						Type:        "string",
						Description: "Comma-separated categories (e.g. general, news, science, it).",
					},
					"time_range": {
						Type:        "string",
						Description: "Time range filter ('day', 'week', 'month', 'year').",
						Enum:        []string{"day", "week", "month", "year"},
					},
					"language": {
						Type:        "string",
						Description: "Language code (e.g. en, fa, all).",
					},
					"limit": {
						Type:        "integer",
						Description: "Number of search results to return (default: 8).",
						Default:     8,
					},
				},
				Required: []string{"query"},
			},
		},
		{
			Name:        "tavily_search",
			Description: "Search the web for real-time information, summaries, and relevant links using Tavily API.",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]PropertyDef{
					"query": {
						Type:        "string",
						Description: "The search query.",
					},
					"search_depth": {
						Type:        "string",
						Description: "Depth of search: 'basic' (fast) or 'advanced' (thorough).",
						Enum:        []string{"basic", "advanced"},
						Default:     "basic",
					},
					"max_results": {
						Type:        "integer",
						Description: "Number of results to return (default: 5).",
						Default:     5,
					},
				},
				Required: []string{"query"},
			},
		},
	}
}

func (s *Server) Serve(ctx context.Context) error {
	scanner := bufio.NewScanner(s.in)
	// Buffer up to 10MB lines for large inputs/outputs
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var req Request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			s.sendError(nil, -32700, "Parse error: "+err.Error())
			continue
		}

		s.handleRequest(ctx, &req)
	}

	return scanner.Err()
}

func (s *Server) handleRequest(ctx context.Context, req *Request) {
	// Notifications (no ID) require no response
	if len(req.ID) == 0 || string(req.ID) == "null" {
		return
	}

	switch req.Method {
	case "initialize":
		s.sendResult(req.ID, map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]any{"tools": map[string]any{}},
			"serverInfo": map[string]any{
				"name":    "miku-researcher",
				"version": "1.0.0",
			},
		})

	case "ping":
		s.sendResult(req.ID, map[string]any{})

	case "tools/list":
		s.sendResult(req.ID, map[string]any{
			"tools": s.GetTools(),
		})

	case "tools/call":
		var params ToolCallParams
		if err := json.Unmarshal(req.Params, &params); err != nil {
			s.sendError(req.ID, -32602, "Invalid params: "+err.Error())
			return
		}
		res, err := s.callTool(ctx, params.Name, params.Arguments)
		if err != nil {
			s.sendResult(req.ID, ToolResult{
				IsError: true,
				Content: []ContentItem{{Type: "text", Text: err.Error()}},
			})
			return
		}
		s.sendResult(req.ID, res)

	default:
		s.sendError(req.ID, -32601, fmt.Sprintf("Method not found: %s", req.Method))
	}
}

func (s *Server) callTool(ctx context.Context, name string, rawArgs json.RawMessage) (*ToolResult, error) {
	switch name {
	case "deep_research":
		var args struct {
			Topic              string `json:"topic"`
			MaxSources         int    `json:"max_sources"`
			MaxCharsPerSource  int    `json:"max_chars_per_source"`
		}
		if err := json.Unmarshal(rawArgs, &args); err != nil {
			return nil, err
		}
		if strings.TrimSpace(args.Topic) == "" {
			return nil, fmt.Errorf("topic parameter cannot be empty")
		}
		dossier, err := s.engine.DeepResearch(ctx, args.Topic, args.MaxSources, args.MaxCharsPerSource)
		if err != nil {
			return nil, err
		}
		return &ToolResult{
			Content: []ContentItem{{Type: "text", Text: dossier}},
		}, nil

	case "read_page":
		var args struct {
			URL      string `json:"url"`
			MaxChars int    `json:"max_chars"`
		}
		if err := json.Unmarshal(rawArgs, &args); err != nil {
			return nil, err
		}
		if strings.TrimSpace(args.URL) == "" {
			return nil, fmt.Errorf("url parameter cannot be empty")
		}
		res, err := s.pageReader.ReadPage(ctx, args.URL, args.MaxChars)
		if err != nil {
			return nil, err
		}
		text := res.Content
		if res.Error != "" {
			text = fmt.Sprintf("Warning: %s\n\n%s", res.Error, text)
		}
		return &ToolResult{
			Content: []ContentItem{{Type: "text", Text: text}},
		}, nil

	case "searxng_search":
		var args struct {
			Query      string `json:"query"`
			Categories string `json:"categories"`
			TimeRange  string `json:"time_range"`
			Language   string `json:"language"`
			Limit      int    `json:"limit"`
		}
		if err := json.Unmarshal(rawArgs, &args); err != nil {
			return nil, err
		}
		if strings.TrimSpace(args.Query) == "" {
			return nil, fmt.Errorf("query parameter cannot be empty")
		}
		items, err := s.searchClient.SearxngSearch(ctx, args.Query, args.Categories, args.TimeRange, args.Language, args.Limit)
		if err != nil {
			return nil, err
		}
		output := orchestrator.FormatSearchResults(items, "")
		return &ToolResult{
			Content: []ContentItem{{Type: "text", Text: output}},
		}, nil

	case "tavily_search":
		var args struct {
			Query       string `json:"query"`
			SearchDepth string `json:"search_depth"`
			MaxResults  int    `json:"max_results"`
		}
		if err := json.Unmarshal(rawArgs, &args); err != nil {
			return nil, err
		}
		if strings.TrimSpace(args.Query) == "" {
			return nil, fmt.Errorf("query parameter cannot be empty")
		}
		items, answer, err := s.searchClient.TavilySearch(ctx, args.Query, args.SearchDepth, args.MaxResults)
		if err != nil {
			return nil, err
		}
		output := orchestrator.FormatSearchResults(items, answer)
		return &ToolResult{
			Content: []ContentItem{{Type: "text", Text: output}},
		}, nil

	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}

func (s *Server) sendResult(id json.RawMessage, result any) {
	resp := Response{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	data, _ := json.Marshal(resp)
	fmt.Fprintf(s.out, "%s\n", data)
}

func (s *Server) sendError(id json.RawMessage, code int, msg string) {
	resp := Response{
		JSONRPC: "2.0",
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: msg,
		},
	}
	data, _ := json.Marshal(resp)
	fmt.Fprintf(s.out, "%s\n", data)
}
