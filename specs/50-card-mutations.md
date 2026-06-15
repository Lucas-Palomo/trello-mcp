# Spec: Minimal Card Mutations

## Goal

Add a small, controlled set of write operations after the read path is stable.

## Scope

- Add `trello_create_card`
- Add `trello_update_card`
- Add `trello_move_card`
- Add `trello_add_comment`

## Tool Contracts

### `trello_create_card`

Create a new card in a list.

#### Input

- `list_id` (required, string)
- `name` (required, string)
- `desc` (optional, string)

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

### `trello_update_card`

Update editable fields on a card.

#### Input

- `card_id` (required, string)
- `name` (optional, string)
- `desc` (optional, string)
- `due` (optional, string in RFC3339 or a format the implementation standardizes explicitly)

At least one mutable field must be provided.

### `trello_move_card`

Move a card to another list.

#### Input

- `card_id` (required, string)
- `list_id` (required, string)

### `trello_add_comment`

Add a comment to a card.

#### Input

- `card_id` (required, string)
- `text` (required, string)

## Mutation Rules

- Validate required inputs before calling Trello.
- Return the updated or created card when Trello provides enough data to do so cleanly.
- Keep mutation scope intentionally small. Do not broaden into labels, attachments, or checklist editing in this phase.

## Error Model

- missing required argument -> MCP tool error
- invalid mutation payload -> MCP tool error
- auth failure -> clear tool error
- not found -> clear tool error
- Trello validation failure -> include useful status/message
- timeout/network failure -> distinguishable error

## Acceptance Criteria

- All four mutation tools are implemented with focused tests.
- Input validation prevents obviously invalid writes from reaching Trello.
- Docs clearly distinguish read tools from write tools.

## Dependencies

- `30-read-details.md`
