# Spec: Read Detail Tools

## Goal

Add focused read tools for a single board and a single card so MCP clients can inspect one resource without fetching larger collections.

## Scope

- Add `trello_get_board`
- Add `trello_get_card`

## Tool Contracts

### `trello_get_board`

Get details for one board.

#### Input

- `board_id` (required, string)

#### Output

```json
{
  "board": {
    "id": "b1",
    "name": "Product Roadmap",
    "url": "https://trello.com/b/b1/product-roadmap"
  }
}
```

Initial board detail scope should stay small. Only add fields when MCP consumers need them.

### `trello_get_card`

Get details for one card.

#### Input

- `card_id` (required, string)

#### Output

```json
{
  "card": {
    "id": "c1",
    "name": "Setup CI/CD pipeline",
    "url": "https://trello.com/c/c1",
    "desc": "Configure GitHub Actions for automated builds"
  }
}
```

Initial card detail scope may later grow to include fields like due date or list ID, but this phase should keep the contract intentionally small.

## Design Requirements

- Reuse the Trello client foundation rather than bypassing it with direct server-layer HTTP.
- Keep result envelopes explicit (`board`, `card`) instead of returning raw objects.

## Acceptance Criteria

- Both tools are available in the MCP server.
- Both tools validate required input at the MCP boundary.
- Tests cover success, missing input, auth failure, 404, and invalid upstream response.

## Dependencies

- `20-read-navigation.md`
