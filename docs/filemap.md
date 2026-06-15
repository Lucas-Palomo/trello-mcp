# File Map

Quick map to avoid unnecessary scouting.

## Source Files

- `main.go`: MCP server entrypoint (stdio)
- `cmd/mcp-client/main.go`: CLI entrypoint, delegates to `internal/cli`
- `internal/cli/`: CLI (Cobra-based), organized by area (boards.go, cards.go, etc.), output rendering, MCP client factory
- `internal/config/`: Environment-based configuration loader
- `internal/server/`: MCP tool registration and request handlers
- `internal/trello/`: Trello HTTP client (`doGet`/`doPost`/`doPut`), domain models

## Documentation

- `README.md`: top-level project overview and onboarding.
- `docs/agents.md`: instructions for coding agents in this repository.
- `docs/developer-integration.md`: guide for developers extending or integrating the project.
- `docs/project.md`: project summary, goals, and current state.
- `docs/entrypoints.md`: server and CLI usage reference.
- `docs/features.md`: current MCP tool and CLI capability overview.
- `docs/filemap.md`: this map.
- `docs/local-install.md`: local build, copy, and MCP host registration guide.
- `docs/roadmap.md`: roadmap and future evolution notes.

## Specifications

- `specs/README.md`: how to use the specifications folder.
- `specs/00-bootstrap-and-config.md`: bootstrap and configuration contract.
- `specs/10-trello-client-foundation.md`: Trello client responsibilities and error model.
- `specs/20-read-navigation.md`: board -> list -> card navigation tools.
- `specs/25-mcp-cli.md`: CLI contract for `cmd/mcp-client`.
- `specs/30-read-details.md`: focused board and card detail reads.
- `specs/40-read-expansion.md`: checklists, actions, and members reads.
- `specs/50-card-mutations.md`: minimal write surface for cards.
- `specs/60-extended-card-model.md`: extended card model with due, labels, members.
- `specs/65-label-reads.md`: label read tools (board labels, card labels).
- `specs/70-label-mutations.md`: add/remove labels on cards.
- `specs/75-card-member-mutations.md`: add/remove members on cards.
- `specs/80-board-list-mutations.md`: create board and list.
- `specs/85-card-archiving.md`: archive, unarchive, delete card.
- `specs/90-card-search.md`: search cards by query.
- `specs/95-organization-reads.md`: list organizations (workspaces).
