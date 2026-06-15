# Agents Guide

This document is for coding agents working in `trello-mcp`.

## Project Context

`trello-mcp` is a Go MCP server for Trello with two executable surfaces:

1. the production MCP server over `stdio`
2. a local Cobra-based CLI that exercises the same MCP contract

The repository is not a bootstrap skeleton anymore. It already contains:

- a Trello HTTP client in `internal/trello`
- MCP tool registration and handlers in `internal/server`
- CLI command handling in `internal/cli`
- focused tests for config, client, server, and CLI behavior
- specs that define the intended contract in `specs/`
- operational and contributor documentation in `docs/`

## What To Read First

Before changing behavior:

1. read [docs/filemap.md](/home/palomo/Projects/Personal/mcps/trello-mcp/docs/filemap.md:1)
2. read [docs/project.md](/home/palomo/Projects/Personal/mcps/trello-mcp/docs/project.md:1)
3. read the relevant files in `specs/`

If your change affects runtime usage or developer workflow, also read:

- [docs/entrypoints.md](/home/palomo/Projects/Personal/mcps/trello-mcp/docs/entrypoints.md:1)
- [docs/developer-integration.md](/home/palomo/Projects/Personal/mcps/trello-mcp/docs/developer-integration.md:1)

## Core Rules

- Treat `specs/` as the source of truth for server, CLI, and contract behavior.
- If implementation diverges from a spec, either fix the code or update the spec intentionally.
- When a spec changes a user-facing contract, update the relevant docs in the same change.
- Keep the MCP surface task-oriented and intentionally small.
- Do not bypass `internal/trello` with direct HTTP from other layers.
- Do not bypass MCP tools from the CLI with direct Trello calls.

## Expected Working Style

- Prefer small, verifiable changes.
- Preserve existing repository patterns unless there is a clear technical reason not to.
- Validate inputs at the boundary:
  - environment/config
  - MCP tool arguments
  - Trello response parsing
- Keep result envelopes explicit and stable.
- Add tests that scale with the change:
  - Trello client tests for HTTP behavior
  - server tests for MCP handler behavior
  - CLI tests for command dispatch and output mode where relevant

## Repository Structure

- `main.go`: production MCP server entrypoint
- `cmd/mcp-client/main.go`: CLI entrypoint
- `internal/config/`: environment-based config loading
- `internal/trello/`: Trello client, request helpers, and models
- `internal/server/`: MCP tool registration and handlers
- `internal/cli/`: Cobra commands, local MCP client bootstrap, and renderers
- `docs/`: user and contributor documentation
- `specs/`: contract definitions and surface evolution

## Documentation Rules

When changing behavior, update the directly affected documentation:

- `README.md` for top-level onboarding
- `docs/entrypoints.md` for runtime usage and current command/tool surface
- `docs/features.md` for capability overview
- `docs/developer-integration.md` for developer-facing architecture/integration guidance
- `docs/filemap.md` if the repo structure changed materially

## Testing Notes

Typical command:

```bash
go test ./...
```

Some environments may restrict local networking or writable build cache paths. If a test failure looks environment-specific, separate that from actual code regressions before drawing conclusions.

## What To Avoid

- Treating the repository as an unimplemented scaffold
- Expanding the surface without updating specs and docs
- Introducing broad abstractions without real complexity pressure
- Coupling transport concerns, Trello HTTP details, and MCP contract logic into the same layer
