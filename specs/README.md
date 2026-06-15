# Specs

This folder defines the target contracts for the current `trello-mcp` project. These specs are not bootstrap notes anymore; they describe the implementation that now exists or the next constrained evolutions that should be applied on top of it.

## How To Use These Specs

- Read the relevant spec before changing behavior.
- Treat specs as the contract for code, tests, CLI behavior, and operational documentation.
- When implementation diverges from a spec, either update the spec intentionally or bring the code back into alignment.
- When a spec changes a user-facing contract, update the affected files in `docs/` during the same change.

## Current Architecture Coverage

The specs are organized by product area and implementation phase:

1. `00-bootstrap-and-config.md`
   - server startup
   - environment-based configuration
2. `10-trello-client-foundation.md`
   - Trello HTTP client contract
   - request/error handling model
3. `20-read-navigation.md`
   - board -> list -> card navigation tools
4. `25-mcp-cli.md`
   - MCP CLI behavior and command surface
5. `26-mcp-cli-cobra.md`
   - migration of the CLI to `github.com/spf13/cobra`
6. `30-read-details.md`
   - board and card detail reads
7. `40-read-expansion.md`
   - checklists, actions, and members
8. `50-card-mutations.md`
   - minimal write support for cards
9. `60-extended-card-model.md`
   - extended card fields such as due dates, labels, members, and archived state
10. `65-label-reads.md`
   - board label and card label read tools
11. `70-label-mutations.md`
   - add and remove labels on cards
12. `75-card-member-mutations.md`
   - add and remove member assignments on cards
13. `80-board-list-mutations.md`
   - create boards and lists
14. `85-card-archiving.md`
   - archive, unarchive, and delete cards
15. `90-card-search.md`
   - search cards by text, optionally scoped to a board
16. `95-organization-reads.md`
   - list organizations (workspaces)

## Current Product Direction

- The MCP server is the production entrypoint.
- The CLI is a real operational/manual-testing surface, not a one-off smoke test.
- The canonical read flow is:
  1. `trello_list_boards`
  2. `trello_list_lists(board_id)`
  3. `trello_list_cards(list_id)`
- Detailed reads, label operations, board/list creation, card search, organization reads, and card mutations are part of the supported surface.

## Notes For The Coding Agent

- Do not preserve the old `trello_list_cards(board_id)` contract.
- Keep MCP contracts intentionally small and typed.
- Prefer evolving the current architecture over reintroducing bootstrap-era assumptions.
