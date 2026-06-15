# Spec: Read Expansion Tools

## Goal

Add high-value supporting reads around cards after the core navigation and detail flow is stable.

## Scope

- Add `trello_list_checklists`
- Add `trello_list_card_actions`
- Add `trello_list_card_members`

## Tool Contracts

### `trello_list_checklists`

List checklists attached to a card.

#### Input

- `card_id` (required, string)

#### Output

```json
{
  "checklists": [
    {
      "id": "chk1",
      "name": "Release checklist"
    }
  ]
}
```

Checklist item modeling can be added in this phase only if it is straightforward and clearly useful.

### `trello_list_card_actions`

List relevant actions for a card.

#### Input

- `card_id` (required, string)

#### Output

```json
{
  "actions": [
    {
      "id": "a1",
      "type": "commentCard",
      "text": "Ready for review"
    }
  ]
}
```

The initial scope should prefer comments and other obviously useful user-facing actions over the full Trello action payload.

### `trello_list_card_members`

List members assigned to a card.

#### Input

- `card_id` (required, string)

#### Output

```json
{
  "members": [
    {
      "id": "m1",
      "full_name": "Jane Doe",
      "username": "jane"
    }
  ]
}
```

## Design Requirements

- Keep each response narrowly modeled.
- Avoid passing through the full Trello action or member schema.
- Add only the fields needed for useful MCP interaction.

## Acceptance Criteria

- All three tools are available and documented.
- Tests cover success and core error cases for each tool.
- The card action contract stays filtered and understandable.

## Dependencies

- `30-read-details.md`
