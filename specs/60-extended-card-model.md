# Spec: Extended Card Model

## Goal

Extend the Card model and its MCP responses to include due dates, start dates, labels, members, list ID, and archived status so MCP clients have richer card data without additional API calls.

## Scope

- Extend the `Card` struct with `Due`, `Start`, `Labels`, `IDMembers`, `IDList`, `Closed`, `DueComplete`
- Add a `Label` model
- Add a `CardLabel` model for inline labels on cards

## Model Changes

### Card (extended)

```go
type Card struct {
    ID          string      `json:"id"`
    Name        string      `json:"name"`
    URL         string      `json:"url"`
    Desc        string      `json:"desc,omitempty"`
    Due         string      `json:"due,omitempty"`
    Start       string      `json:"start,omitempty"`
    DueComplete bool        `json:"dueComplete"`
    Closed      bool        `json:"closed"`
    IDList      string      `json:"idList,omitempty"`
    IDMembers   []string    `json:"idMembers,omitempty"`
    Labels      []CardLabel `json:"labels,omitempty"`
}
```

### CardLabel (new)

```go
type CardLabel struct {
    ID      string `json:"id"`
    IDBoard string `json:"idBoard,omitempty"`
    Name    string `json:"name"`
    Color   string `json:"color"`
}
```

### Label (standalone, for future use)

```go
type Label struct {
    ID      string `json:"id"`
    IDBoard string `json:"idBoard"`
    Name    string `json:"name"`
    Color   string `json:"color"`
}
```

## Impact Analysis

Every tool that returns a Card or list of Cards will now include these additional fields. No new server tools are introduced in this spec.

### Affected Tools

- `trello_list_cards` — Cards in responses now include extended fields
- `trello_get_card` — Card detail now includes extended fields
- `trello_create_card` — Created card response includes extended fields
- `trello_update_card` — Updated card response includes extended fields
- `trello_move_card` — Moved card response includes extended fields

### CLI Output

The CLI printers for card list and card detail should display due date and labels if present.

## Dependencies

- `30-read-details.md`
- `50-card-mutations.md`
