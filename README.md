<p align="center">
  <img src="banner.png" alt="Miku Researcher Banner" width="100%">
</p>

<h1 align="center">Miku Researcher (miku-researcher)</h1>

<p align="center">
  <strong>Ultra-fast, token-efficient, autonomous web research MCP server written in pure Go.</strong>
</p>

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-teal.svg" alt="License: MIT"></a>
  <a href="https://golang.org"><img src="https://img.shields.io/badge/Go-1.22+-00ADD8.svg" alt="Go Version"></a>
  <a href="https://nixos.org"><img src="https://img.shields.io/badge/Nix-Flake_Ready-5277C3.svg" alt="Nix Flake"></a>
  <img src="https://img.shields.io/badge/Dependencies-Zero_External-brightgreen.svg" alt="Zero Dependencies">
</p>

---

## 🌟 Overview

**Miku Researcher** is a high-performance Model Context Protocol (MCP) server designed specifically for autonomous AI agents requiring deep technical research, documentation ingestion, and web synthesis without bloating context windows or wasting tokens.

Built with **Pure Go (Zero-CGO, Zero external dependencies)**, it integrates seamlessly with local **SearXNG** instances, **Tavily Search**, and leverages **Defuddle** for extracting pristine, distraction-free markdown articles.

---

## ⚡ Key Capabilities

- **`deep_research` (Autonomous Intelligence Dossier):** Takes a topic, searches premier sources, concurrently reads and cleans the destination pages with Defuddle, and synthesizes a high-density intelligence brief in a single agent call.
- **`read_page` (Boilerplate-free Reader):** Strips navigation bars, ads, cookie banners, and CSS clutter from any webpage, outputting clean, token-efficient markdown.
- **`searxng_search` (Local Self-Hosted Search):** Unlimited, quota-free queries against your private SearXNG gateway. Formatted as concise markdown, never raw noisy JSON.
- **`tavily_search` (Real-Time AI Web Search):** Quick instant answers and ranked sources when external depth is needed.
- **Zero-CGO & Minimalist:** Compiles to a lightning-fast static single binary with negligible RAM footprint.

---

## 🛠️ MCP Tools Reference

| Tool | Parameters | Description |
| :--- | :--- | :--- |
| `deep_research` | `topic` *(req)*, `max_sources`, `max_chars_per_source` | Autonomous multi-source deep research, concurrent reading, and dossier assembly. |
| `read_page` | `url` *(req)*, `max_chars` | Extracts distilled markdown from a webpage via Defuddle. |
| `searxng_search` | `query` *(req)*, `categories`, `time_range`, `language`, `limit` | High-density search via SearXNG without consuming external API tokens. |
| `tavily_search` | `query` *(req)*, `search_depth`, `max_results` | Real-time web search and quick answers with Tavily. |

---

## 🚀 Installation & Build

### Using Nix Flakes (Recommended)

```bash
# Run directly
nix run github:surtr85/miku-researcher

# Or build the binary locally
nix build
./result/bin/miku-researcher
```

### From Source (Go 1.22+)

```bash
git clone https://github.com/surtr85/miku-researcher.git
cd miku-researcher
CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/miku-researcher ./cmd/miku-researcher
```

---

## ⚙️ Configuration

Set environment variables in your MCP host configuration or shell environment:

```env
SEARXNG_URL=https://searxng.surtr.ir   # Optional (Defaults to https://searxng.surtr.ir)
TAVILY_API_KEY=tvly-xxxxxxxxxxxx       # Optional (For Tavily integration)
```

Ensure `defuddle` is available in your `$PATH` for webpage cleaning capabilities.

---

## 🤝 Contributing & License

Contributions, issues, and feature suggestions are welcome! Distributed under the **MIT License**.
