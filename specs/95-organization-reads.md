# Spec: Organization Reads

## Goal

Add a read-only tool for listing Trello organizations (workspaces) accessible to the configured account.

## Scope

- Add `trello_list_organizations`

## Tool Contract

### `trello_list_organizations`

List all organizations (workspaces) accessible to the configured credentials.

#### Input

None.

#### Output

```json
{
  "organizations": [
    {"id": "o1", "name": "acmecorp", "displayName": "Acme Corp", "url": "https://trello.com/acmecorp"}
  ]
}
```

## Dependencies

- `20-read-navigation.md`
