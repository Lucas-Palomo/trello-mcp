# Project Summary

## Name

`trello-mcp`

## Summary

`trello-mcp` is a Go MCP server that exposes a focused Trello tool surface for MCP clients, plus a companion CLI for manual operation and verification. The project is contract-driven: Trello operations are modeled as MCP tools with explicit inputs, typed outputs, validation at the boundary, and predictable operational behavior.

## Current State

The repository currently includes:

- a stdio MCP server entrypoint in `main.go`
- environment-based configuration loading
- a dedicated Trello client in `internal/trello`
- an MCP server layer in `internal/server`
- a CLI in `cmd/mcp-client` / `internal/cli`
- tests for config, CLI behavior, Trello client behavior, and server handlers

The current tool surface includes:

- navigation reads:
  - `trello_list_boards`
  - `trello_list_lists`
  - `trello_list_cards`
- detail reads:
  - `trello_get_board`
  - `trello_get_card`
- expansion reads:
  - `trello_list_checklists`
  - `trello_list_card_actions`
  - `trello_list_card_members`
- card mutations:
  - `trello_create_card`
  - `trello_update_card`
  - `trello_move_card`
  - `trello_add_comment`

See `docs/entrypoints.md` for runtime usage and examples.

## Product Shape

The project has two entry surfaces:

1. MCP server
   - production entrypoint
   - stdio transport
   - intended to be launched by an MCP host

2. MCP CLI
   - local operational/manual-testing surface
   - executes commands against the MCP server contract
   - currently uses an in-process MCP client

## Technical Direction

- Language: Go
- Style: small, explicit, contract-driven implementation
- Server transport: stdio MCP server
- Trello integration: isolated in `internal/trello`
- Functional source of truth: `specs/`
- Operational documentation: `docs/`

## Current Contract Direction

- The canonical read flow is `boards -> lists -> cards`.
- `trello_list_cards(list_id)` is the supported card-listing contract.
- Detailed reads and narrow mutations are preferred over broad API mirroring.
- Models should remain intentionally small and MCP-oriented.

## Constraints

- The project does not aim to cover the full Trello API.
- Complex local persistence is out of scope.
- Multi-tenant credential management is out of scope.
- A rich UI is out of scope; the CLI remains operational rather than interactive.

## Near-Term Evolution

- keep expanding the MCP tool surface conservatively, following the contracts in `specs/`
- preserve the Cobra-based CLI structure as new commands are added
- keep the current in-process CLI connectivity model easy to replace with another MCP transport later
