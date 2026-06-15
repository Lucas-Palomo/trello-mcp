# Spec: Card Search

## Goal

Add a search tool that queries Trello cards across the workspace.

## Scope

- Add `trello_search_cards`

## Tool Contract

### `trello_search_cards`

Search for cards by text query, optionally scoped to a board.

#### Input

- `query` (required, string) — text to search for
- `board_id` (optional, string) — scope search to a specific board

#### Output

```json
{
  "cards": [
    {
      "id": "c1",
      "name": "Fix login bug",
      "desc": "Users cannot log in",
      "idList": "l1",
      "idBoard": "b1"
    }
  ]
}
```

## Dependencies

- `20-read-navigation.md`
- `60-extended-card-model.md`
