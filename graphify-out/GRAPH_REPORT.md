# Graph Report - miku-researcher  (2026-10-06)

## Corpus Check
- cluster-only mode — file stats not available

## Summary
- 56 nodes · 97 edges · 8 communities (7 shown, 1 thin omitted)
- Extraction: 98% EXTRACTED · 2% INFERRED · 0% AMBIGUOUS · INFERRED: 2 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `2dccf3d4`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Community 0
- Community 1
- Community 2
- Community 3
- Community 4
- Community 5
- Community 6
- Community 7

## God Nodes (most connected - your core abstractions)
1. `Server` - 13 edges
2. `NewServer()` - 9 edges
3. `Client` - 7 edges
4. `ResearchEngine` - 6 edges
5. `Reader` - 6 edges
6. `NewResearchEngine()` - 5 edges
7. `NewReader()` - 5 edges
8. `ResultItem` - 4 edges
9. `NewClient()` - 4 edges
10. `FormatSearchResults()` - 4 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `NewServer()`  [EXTRACTED]
  cmd/miku-researcher/main.go → internal/mcp/server.go
- `TestMCPServerLifecycle()` --calls--> `NewServer()`  [INFERRED]
  internal/mcp/server_test.go → internal/mcp/server.go
- `main()` --calls--> `AutoLoadEnv()`  [EXTRACTED]
  cmd/miku-researcher/main.go → internal/envutil/env.go
- `ResearchEngine` --references--> `Reader`  [EXTRACTED]
  internal/orchestrator/engine.go → internal/reader/defuddle.go
- `Server` --references--> `ResearchEngine`  [EXTRACTED]
  internal/mcp/server.go → internal/orchestrator/engine.go

## Import Cycles
- None detected.

## Communities (8 total, 1 thin omitted)

### Community 0 - "Community 0"
Cohesion: 0.24
Nodes (10): encoding/json.RawMessage, ContentItem, InputSchema, PropertyDef, Request, Response, RPCError, Tool (+2 more)

### Community 1 - "Community 1"
Cohesion: 0.27
Nodes (7): context.Context, net/http.Client, ResearchEngine, NewResearchEngine(), Client, ResultItem, NewClient()

### Community 2 - "Community 2"
Cohesion: 0.42
Nodes (4): io.Reader, io.Writer, NewServer(), Server

### Community 3 - "Community 3"
Cohesion: 0.33
Nodes (4): testing.T, TestMCPServerLifecycle(), FormatSearchResults(), TestFormatSearchResults()

### Community 4 - "Community 4"
Cohesion: 0.53
Nodes (4): time.Duration, Reader, NewReader(), PageResult

### Community 5 - "Community 5"
Cohesion: 0.50
Nodes (3): main(), AutoLoadEnv(), loadEnvFile()

### Community 6 - "Community 6"
Cohesion: 0.50
Nodes (3): miku-researcher.nix, pkgs.buildGoModule, pkgs.mkShell

## Knowledge Gaps
- **4 isolated node(s):** `miku-researcher.nix`, `pkgs.buildGoModule`, `pkgs.mkShell`, `github.com/surtr85/miku-researcher`
  These have ≤1 connection - possible missing edges or undocumented components.
- **1 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Server` connect `Community 2` to `Community 0`, `Community 1`, `Community 4`?**
  _High betweenness centrality (0.335) - this node is a cross-community bridge._
- **Why does `NewServer()` connect `Community 2` to `Community 1`, `Community 3`, `Community 4`, `Community 5`?**
  _High betweenness centrality (0.255) - this node is a cross-community bridge._
- **Why does `main()` connect `Community 5` to `Community 2`?**
  _High betweenness centrality (0.126) - this node is a cross-community bridge._
- **What connects `miku-researcher.nix`, `pkgs.buildGoModule`, `pkgs.mkShell` to the rest of the system?**
  _4 weakly-connected nodes found - possible documentation gaps or missing edges._