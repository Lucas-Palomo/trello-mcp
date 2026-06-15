# trello-mcp

`trello-mcp` is a Go MCP server for Trello. It exposes a focused tool surface for reading and mutating boards, lists, cards, labels, members, and organizations through the Model Context Protocol, plus a local CLI for manual operation and verification.

## What It Is

This project provides two executable surfaces:

1. an MCP server intended to be launched by an MCP host such as Claude Desktop or Cursor
2. a local CLI for exercising the same MCP contract during development and manual testing

The implementation is intentionally narrow. It does not try to mirror the full Trello API. It focuses on the operations that are useful in MCP workflows and keeps request/response contracts explicit.

## Pre-requisites

- Go `1.26` or newer
- A Trello API key
- A Trello token
- An MCP host if you want to use the server from another MCP client

Required environment variables:

- `TRELLO_API_KEY`
- `TRELLO_TOKEN`

Optional environment variables:

- `HTTP_TIMEOUT` in seconds, default `30`

## Installation With MCP

The recommended packaging path is to distribute the server binary through GitHub releases or another artifact pipeline.

Tagged releases publish binary artifacts automatically through GitHub Actions.

This MCP server uses `stdio`. It does not expose an HTTP port. The MCP host must start the binary as a subprocess and communicate through `stdin`/`stdout`.

Example connection shape:

```json
{
  "name": "trello-mcp",
  "command": "/absolute/path/to/trello-mcp",
  "env": {
    "TRELLO_API_KEY": "your-key",
    "TRELLO_TOKEN": "your-token",
    "HTTP_TIMEOUT": "30"
  }
}
```

For local build, copy, and MCP host registration, see:

- [docs/local-install.md](/home/palomo/Projects/Personal/mcps/trello-mcp/docs/local-install.md:1)

### Local CLI usage

The CLI runs against the same MCP contract and is useful for verification:

```bash
go run ./cmd/mcp-client/ boards list
go run ./cmd/mcp-client/ cards search --query "bug" --json
```

More runtime details are in [docs/entrypoints.md](/home/palomo/Projects/Personal/mcps/trello-mcp/docs/entrypoints.md:1).

## Features

At a high level, the project currently supports:

- navigation across boards, lists, and cards
- detail reads for boards and cards
- checklist, action, member, label, and organization reads
- card creation, updates, moves, comments, archive/unarchive, and delete
- label assignment and removal
- member assignment and removal
- board and list creation
- card search
- a Cobra-based CLI for local operation and manual testing

For the complete tool surface and command overview, see [docs/features.md](/home/palomo/Projects/Personal/mcps/trello-mcp/docs/features.md:1).

## Developer Integration Guide

If you want to extend the server, add new tools, or integrate it into another MCP environment, start here:

- [docs/developer-integration.md](/home/palomo/Projects/Personal/mcps/trello-mcp/docs/developer-integration.md:1)
- [docs/filemap.md](/home/palomo/Projects/Personal/mcps/trello-mcp/docs/filemap.md:1)
- [specs/README.md](/home/palomo/Projects/Personal/mcps/trello-mcp/specs/README.md:1)

## Additional Documentation

- [docs/entrypoints.md](/home/palomo/Projects/Personal/mcps/trello-mcp/docs/entrypoints.md:1): server and CLI usage
- [docs/local-install.md](/home/palomo/Projects/Personal/mcps/trello-mcp/docs/local-install.md:1): local build, copy, and MCP host registration
- [docs/features.md](/home/palomo/Projects/Personal/mcps/trello-mcp/docs/features.md:1): tool and command overview
- [docs/developer-integration.md](/home/palomo/Projects/Personal/mcps/trello-mcp/docs/developer-integration.md:1): implementation and extension guide
- [docs/roadmap.md](/home/palomo/Projects/Personal/mcps/trello-mcp/docs/roadmap.md:1): roadmap and future evolution notes
