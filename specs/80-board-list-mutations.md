# Spec: Board and List Mutations

## Goal

Add write operations for creating boards and lists.

## Scope

- Add `trello_create_board`
- Add `trello_create_list`

## Tool Contracts

### `trello_create_board`

Create a new Trello board.

#### Input

- `name` (required, string)
- `default_lists` (optional, bool, default `true`) — whether to include default lists

#### Output

```json
{
  "board": {
    "id": "b1",
    "name": "My Board",
    "url": "https://trello.com/b/b1"
  }
}
```

### `trello_create_list`

Create a new list inside a board.

#### Input

- `name` (required, string)
- `board_id` (required, string)

#### Output

```json
{
  "list": {
    "id": "l1",
    "name": "To Do"
  }
}
```

## Dependencies

- `20-read-navigation.md`
