# Spec: Card Archiving

## Goal

Add card archive, unarchive, and delete operations.

## Scope

- Add `trello_archive_card`
- Add `trello_unarchive_card`
- Add `trello_delete_card`

## Tool Contracts

### `trello_archive_card`

Archive (close) a card.

#### Input

- `card_id` (required, string)

#### Output

```json
{
  "card": {
    "id": "c1",
    "closed": true
  }
}
```

### `trello_unarchive_card`

Unarchive (re-open) a card.

#### Input

- `card_id` (required, string)

#### Output

```json
{
  "card": {
    "id": "c1",
    "closed": false
  }
}
```

### `trello_delete_card`

Permanently delete a card.

#### Input

- `card_id` (required, string)

#### Output

```json
{"success": true}
```

## Dependencies

- `50-card-mutations.md`
