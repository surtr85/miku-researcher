---
name: miku-researcher
description: Ultra-fast, token-efficient, autonomous web research MCP server written in pure Go (Zero-CGO). Provides deep technical research, documentation ingestion, and web synthesis via SearXNG, Tavily, and Defuddle markdown extraction.
---

# Miku Researcher (`miku-researcher`)

**Miku Researcher** is a high-performance Model Context Protocol (MCP) server for deep technical research, documentation ingestion, and web synthesis without bloating context windows or wasting tokens.

Built with **Pure Go (Zero-CGO)**, it integrates with local **SearXNG** instances, **Tavily Search**, and leverages **Defuddle** for extracting pristine, distraction-free markdown articles.

---

## 🛠️ MCP Tools Reference (5 Tools)

### 1. `deep_research`

Autonomous multi-source deep research, concurrent reading, and comprehensive dossier assembly in a single agent step.

- **Parameters**:
  - `topic` _(string, required)_: Research topic or question.
  - `max_sources` _(integer, optional)_: Maximum number of sources to explore (default: 3).
  - `max_chars_per_source` _(integer, optional)_: Cap character count per extracted source to prevent token bloat (default: 8000).
- **Codemode Invocation**:
  ```javascript
  await tools.mcp__miku_researcher__deep_research({
    topic: "Vulkan dynamic rendering architecture in Rust",
    max_sources: 3,
  });
  ```

### 2. `read_page`

Extracts distilled, boilerplate-free markdown from any webpage via Defuddle. Strips navigation bars, cookie banners, tracking scripts, and ads.

- **Parameters**:
  - `url` _(string, required)_: Webpage URL to fetch and clean.
  - `max_chars` _(integer, optional)_: Maximum characters to return.
- **Codemode Invocation**:
  ```javascript
  await tools.mcp__miku_researcher__read_page({
    url: "https://docs.kernel.org/gpu/amdgpu/index.html",
  });
  ```

### 3. `searxng_search`

High-density search via self-hosted SearXNG without consuming external API quota or tokens. Output is formatted as clean markdown links and summaries.

- **Parameters**:
  - `query` _(string, required)_: Search keywords.
  - `categories` _(string, optional)_: Category filter (e.g., `general`, `it`, `science`).
  - `time_range` _(string, optional)_: Filter by time (`day`, `week`, `month`, `year`).
  - `language` _(string, optional)_: Preferred language code.
  - `limit` _(integer, optional)_: Maximum result count.
- **Codemode Invocation**:
  ```javascript
  await tools.mcp__miku_researcher__searxng_search({
    query: "NixOS flakes Lanzaboote secure boot setup",
    limit: 5,
  });
  ```

### 4. `tavily_search`

Real-time web search and quick AI direct answers powered by Tavily API.

- **Parameters**:
  - `query` _(string, required)_: Search query string.
  - `search_depth` _(string, optional)_: Search depth (`basic` or `advanced`).
  - `max_results` _(integer, optional)_: Result count.
- **Codemode Invocation**:
  ```javascript
  await tools.mcp__miku_researcher__tavily_search({
    query: "Stockfish 18 release notes and neural net benchmarks",
  });
  ```

### 5. `firecrawl_scrape`

Direct webpage scraping and extraction of clean markdown via Firecrawl API for complex or dynamic JavaScript sites.

- **Parameters**:
  - `url` _(string, required)_: Target webpage URL.
  - `max_chars` _(integer, optional)_: Maximum characters to return (default: 6000).
  - `only_main_content` _(boolean, optional)_: Only extract main content (default: true).
- **Codemode Invocation**:
  ```javascript
  await tools.mcp__miku_researcher__firecrawl_scrape({
    url: "https://example.com",
  });
  ```

---

## ⚡ Agent Usage Protocol

- For fast verification and technical queries, favor `searxng_search` first.
- For deep technical synthesis across multiple articles, use `deep_research`.
- When reading external documentation or blog posts from search results, use `read_page`.
