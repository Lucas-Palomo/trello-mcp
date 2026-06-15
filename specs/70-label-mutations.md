# Spec: Label Mutations

## Goal

Add write operations for managing labels on cards.

## Scope

- Add `trello_add_card_label`
- Add `trello_remove_card_label`

## Tool Contracts

### `trello_add_card_label`

Attach a label to a card.

#### Input

- `card_id` (required, string)
- `label_id` (required, string)

#### Output

```json
{
  "success": true
}
```

### `trello_remove_card_label`

Remove a label from a card.

#### Input

- `card_id` (required, string)
- `label_id` (required, string)

#### Output

```json
{
  "success": true
}
```

## Mutation Rules

- Validate required inputs before calling Trello.
- Return success confirmation on 200 OK.
- Missing resource or auth failures produce clear tool errors.

## Dependencies

- `65-label-reads.md`
