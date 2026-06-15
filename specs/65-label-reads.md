# Spec: Label Reads

## Goal

Add read-only tools for Trello labels on boards and cards, enabling MCP clients to inspect available labels and per-card label assignments.

## Scope

- Add `trello_list_board_labels`
- Add `trello_list_card_labels`

## Tool Contracts

### `trello_list_board_labels`

List all labels defined on a board.

#### Input

- `board_id` (required, string)

#### Output

```json
{
  "labels": [
    {"id": "lb1", "idBoard": "b1", "name": "urgent", "color": "red"},
    {"id": "lb2", "idBoard": "b1", "name": "bug", "color": "orange"}
  ]
}
```

### `trello_list_card_labels`

List labels applied to a specific card.

#### Input

- `card_id` (required, string)

#### Output

```json
{
  "labels": [
    {"id": "lb1", "idBoard": "b1", "name": "urgent", "color": "red"}
  ]
}
```

## Design Requirements

- Reuse the existing `Label` model (already defined in models.go)
- Keep responses consistent with the existing label model

## Dependencies

- `60-extended-card-model.md`
