# Spec: Card Member Mutations

## Goal

Add write operations for managing member assignments on cards.

## Scope

- Add `trello_add_card_member`
- Add `trello_remove_card_member`

## Tool Contracts

### `trello_add_card_member`

Assign a member to a card.

#### Input

- `card_id` (required, string)
- `member_id` (required, string)

#### Output

```json
{"success": true}
```

### `trello_remove_card_member`

Remove a member from a card.

#### Input

- `card_id` (required, string)
- `member_id` (required, string)

#### Output

```json
{"success": true}
```

## Dependencies

- `40-read-expansion.md`
