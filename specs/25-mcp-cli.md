# Spec: MCP CLI

## Goal

Turn `cmd/mcp-client` from an in-process test program into a real CLI for operating and inspecting the MCP server.

## Problem

The current `cmd/mcp-client/main.go` behaves like a smoke test:

- it starts the server in-process
- it lists tools automatically
- it walks boards, lists, and cards without explicit user intent
- it mixes bootstrap, demo flow, JSON formatting, and error handling in one file

That is acceptable for early manual testing, but it is not a usable CLI contract.

## Scope

- Redefine `cmd/mcp-client` as a CLI entrypoint.
- Replace the automatic demo flow with explicit commands and flags.
- Keep the CLI focused on calling the MCP server tools that exist in the roadmap.
- Support both human-readable output and raw JSON output.
- Standardize exit codes and error presentation.

## Non-Goals

- Building a rich interactive TUI
- Supporting every Trello capability from day one
- Replacing the stdio MCP server entrypoint in `main.go`

## Functional Requirements

### Command Model

The CLI must support explicit subcommands instead of hardcoded traversal.

Minimum command surface for the first phase:

- `mcp-client tools list`
- `mcp-client boards list`
- `mcp-client lists list --board-id <id>`
- `mcp-client cards list --list-id <id>`

The CLI should be designed so future commands can be added without reshaping the whole program:

- `boards get --board-id <id>`
- `cards get --card-id <id>`
- `cards create --list-id <id> --name <name> [--desc <text>]`
- `cards update --card-id <id> [...]`

### Output Modes

The CLI must support:

- default human-readable output
- `--json` for machine-readable output

Human-readable output should stay compact and task-oriented. JSON output should preserve the MCP tool result payload cleanly.

### Error Handling

- Usage errors must print a concise message and non-zero exit code.
- MCP tool errors must be rendered clearly without stack traces.
- Unexpected runtime errors must remain distinguishable from user input errors.
- The CLI should avoid `log.Fatalf` in deep call paths. Errors should bubble to a top-level runner.

### Server Connectivity Model

For the initial implementation, the CLI may continue using an in-process MCP client if that is the simplest way to exercise the server contract.

However, the program structure must not assume a demo-only lifecycle. The design should make it straightforward to later support:

- spawning the stdio server as a subprocess, or
- connecting to another MCP transport

## Design Requirements

- Keep `main.go` for the CLI thin.
- Prefer a `run(args []string) error` style entrypoint.
- Separate these responsibilities:
  - argument parsing
  - MCP client creation
  - command execution
  - output rendering
- Keep formatting helpers reusable across commands.
- Avoid coupling command handlers to Trello HTTP details; they should operate through MCP tools.

## Suggested Structure

One reasonable structure is:

- `cmd/mcp-client/main.go` -> process entrypoint only
- `cmd/mcp-client/cli.go` -> argument parsing and dispatch
- `cmd/mcp-client/client.go` -> MCP client setup
- `cmd/mcp-client/output.go` -> human/JSON rendering

The exact file split may vary, but the separation of concerns should remain.

## Documentation Requirements

- `docs/entrypoints.md` must describe `cmd/mcp-client` as a CLI, not a test client, once the refactor is complete.
- Command examples must match the actual implemented flags and subcommands.
- Documentation must not describe the legacy `trello_list_cards(board_id)` contract after the CLI is updated.

## Acceptance Criteria

- The CLI no longer auto-traverses boards, lists, and cards on startup.
- The CLI exposes explicit subcommands for tools, boards, lists, and cards.
- `cards list` uses `--list-id`, not `--board-id`.
- The CLI supports `--json`.
- Tests cover at least:
  - successful command dispatch
  - missing required flags
  - tool error rendering
  - JSON output mode

## Dependencies

- `20-read-navigation.md`
