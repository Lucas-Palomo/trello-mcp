# Developer Integration Guide

This document is for developers extending `trello-mcp` or integrating it into another MCP environment.

## Architecture

The project is organized around a narrow contract:

- `main.go`: production MCP server entrypoint
- `cmd/mcp-client/main.go`: local CLI entrypoint
- `internal/config/`: environment-based configuration
- `internal/trello/`: Trello HTTP client and models
- `internal/server/`: MCP tool registration and handlers
- `internal/cli/`: Cobra command tree, MCP client bootstrap, and output rendering
- `specs/`: behavior contracts and surface evolution
- `docs/`: operational and contributor-facing documentation

See [filemap.md](/home/palomo/Projects/Personal/mcps/trello-mcp/docs/filemap.md:1) for a quick repo map.

## How The Request Flow Works

1. credentials and timeout are loaded from environment in `internal/config`
2. `main.go` builds a Trello client and MCP server
3. `internal/server` registers MCP tools and validates tool inputs
4. server handlers call `internal/trello`
5. `internal/trello` is the only layer that knows Trello URLs, auth query params, and raw HTTP behavior

The CLI follows the same contract, but currently uses an in-process MCP client for local execution.

## Extending The Server

When adding a new capability:

1. define or update the relevant spec in `specs/`
2. add or update models and HTTP methods in `internal/trello`
3. register the MCP tool and handler in `internal/server`
4. add or extend CLI commands in `internal/cli` if the capability should be available locally
5. add tests for the Trello client, server handler, and CLI path where relevant
6. update docs in `docs/` and top-level `README.md`

## Testing

The project already contains focused tests for:

- config loading
- Trello client methods
- server handlers
- CLI command dispatch

Typical command:

```bash
go test ./...
```

Depending on the execution environment, some tests that use `httptest.NewServer` may require relaxed sandbox/network restrictions.

## Integration Notes For MCP Hosts

- transport is stdio
- the host should launch the server binary as a subprocess
- required env vars must be injected by the host
- the server contract is documented in [entrypoints.md](/home/palomo/Projects/Personal/mcps/trello-mcp/docs/entrypoints.md:1)

## Design Constraints

- do not broaden the API without a clear MCP use case
- keep models intentionally small and useful
- prefer explicit validation at the MCP boundary
- keep Trello-specific HTTP details inside `internal/trello`
- keep CLI commands operating through MCP tools, not direct HTTP calls

## Source Of Truth

- `specs/` defines the intended behavior contract
- `docs/` explains operation and project structure
- `README.md` is the top-level onboarding document
