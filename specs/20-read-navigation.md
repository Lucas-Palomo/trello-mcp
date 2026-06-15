# Spec: Read Navigation Tools

## Goal

Expose the core read path for Trello navigation in the MCP server:

1. list boards
2. list lists in a board
3. list cards in a list

## Problem

The current board-scoped `trello_list_cards(board_id)` contract is too broad and does not match Trello's structure. MCP clients need narrower, stepwise navigation.

## Scope

- Keep `trello_list_boards`
- Add `trello_list_lists`
- Refactor `trello_list_cards` to accept `list_id`
- Update test client and docs to match the new flow

## Tool Contracts

### `trello_list_boards`

List boards visible to the configured Trello credentials.

#### Input

- none

#### Output

```json
{
  "boards": [
    {
      "id": "b1",
      "name": "Product Roadmap",
      "url": "https://trello.com/b/b1/product-roadmap"
    }
  ]
}
```

### `trello_list_lists`

List lists in a board.

#### Input

- `board_id` (required, string)

#### Output

```json
{
  "lists": [
    {
      "id": "l1",
      "name": "Todo"
    }
  ]
}
```

### `trello_list_cards`

List cards in a list.

#### Input

- `list_id` (required, string)

#### Output

```json
{
  "cards": [
    {
      "id": "c1",
      "name": "Setup CI/CD pipeline",
      "url": "https://trello.com/c/c1",
      "desc": "Configure GitHub Actions for automated builds"
    }
  ]
}
```

## Contract Change

This phase intentionally replaces legacy behavior:

- old: `trello_list_cards(board_id)`
- new: `trello_list_cards(list_id)`

No compatibility layer is required unless a real client explicitly depends on the old contract.

## Error Model

- missing required argument -> MCP tool error
- board not found -> clear tool error for `trello_list_lists`
- list not found -> clear tool error for `trello_list_cards`
- authentication failure -> clear tool error
- timeout/network failure -> distinguishable tool error
- other Trello remote error -> include status code
- invalid upstream response -> parse error

## Acceptance Criteria

- The MCP server exposes the three navigation tools.
- `trello_list_cards` is list-scoped, not board-scoped.
- Tests cover success, missing input, auth failure, not found, and invalid upstream response.
- Operational docs describe the `boards -> lists -> cards` flow.

## Dependencies

- `10-trello-client-foundation.md`
