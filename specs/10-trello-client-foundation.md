# Spec: Trello Client Foundation

## Goal

Define the Trello HTTP client contract that all MCP tools will use.

## Problem

The project needs a single place to handle Trello authentication, timeouts, request building, response parsing, and upstream error mapping. Without that, each MCP tool will duplicate fragile HTTP behavior.

## Scope

- Create or preserve a dedicated `trello` package for API access.
- Support authenticated GET, POST, and PUT requests needed by the roadmap.
- Return typed Go models for MCP-facing data.
- Translate Trello and network failures into actionable Go errors.

## Functional Requirements

- The client must authenticate every request using API key and token from configuration.
- The client must enforce request timeouts via the configured HTTP client.
- The client must distinguish:
  - authentication failures
  - not found errors
  - network/timeout failures
  - unexpected upstream status codes
  - invalid JSON responses

## Initial Client Surface

The client should ultimately support these operations:

- `ListBoards()`
- `ListLists(boardID string)`
- `ListCards(listID string)`
- `GetBoard(boardID string)`
- `GetCard(cardID string)`
- `ListCardChecklists(cardID string)`
- `ListCardActions(cardID string)`
- `ListCardMembers(cardID string)`
- `CreateCard(input CreateCardInput)`
- `UpdateCard(cardID string, input UpdateCardInput)`
- `MoveCard(cardID string, listID string)`
- `AddComment(cardID string, text string)`

## Design Requirements

- Keep request construction inside the client package.
- Keep models aligned with MCP needs, not the full Trello API.
- Add helpers only where they reduce repeated request/response code.
- Do not introduce broad mocking infrastructure unless it is needed for tests.

## Acceptance Criteria

- The client can be instantiated from config and used by the server layer.
- Client tests cover at least one success path and key error paths for each implemented method.
- The client contract is the only place where Trello URLs and auth query params are assembled.

## Dependencies

- `00-bootstrap-and-config.md`
