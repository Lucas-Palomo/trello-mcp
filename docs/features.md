# Features

This document summarizes the current functional surface of `trello-mcp`.

## MCP Server

The server exposes tools in these groups:

### Navigation

- `trello_list_boards`
- `trello_list_lists`
- `trello_list_cards`

### Details

- `trello_get_board`
- `trello_get_card`

### Expansion Reads

- `trello_list_checklists`
- `trello_list_card_actions`
- `trello_list_card_members`

### Organization Reads

- `trello_list_organizations`

### Label Reads

- `trello_list_board_labels`
- `trello_list_card_labels`

### Card Mutations

- `trello_create_card`
- `trello_update_card`
- `trello_move_card`
- `trello_add_comment`

### Card Archiving

- `trello_archive_card`
- `trello_unarchive_card`
- `trello_delete_card`

### Search

- `trello_search_cards`

### Label Mutations

- `trello_add_card_label`
- `trello_remove_card_label`

### Card Member Mutations

- `trello_add_card_member`
- `trello_remove_card_member`

### Board and List Mutations

- `trello_create_board`
- `trello_create_list`

## CLI

The local CLI mirrors the MCP contract through Cobra commands.

### Read Commands

- `tools list`
- `boards list`
- `boards get --board-id <id>`
- `boards create --name <name> [--no-default-lists]`
- `lists list --board-id <id>`
- `lists create --name <name> --board-id <id>`
- `cards list --list-id <id>`
- `cards get --card-id <id>`
- `checklists list --card-id <id>`
- `actions list --card-id <id>`
- `members list --card-id <id>`
- `members add --card-id <id> --member-id <id>`
- `members remove --card-id <id> --member-id <id>`
- `labels board-list --board-id <id>`
- `labels card-list --card-id <id>`
- `orgs list`

### Write Commands

- `cards create --list-id <id> --name <name> [--desc <text>] [--due <rfc3339>]`
- `cards update --card-id <id> [--name <name>] [--desc <text>] [--due <rfc3339>]`
- `cards move --card-id <id> --list-id <id>`
- `cards comment --card-id <id> --text <text>`
- `cards archive --card-id <id>`
- `cards unarchive --card-id <id>`
- `cards delete --card-id <id>`
- `cards search --query <text> [--board-id <id>]`
- `labels add --card-id <id> --label-id <id>`
- `labels remove --card-id <id> --label-id <id>`

### Output Modes

- human-readable output by default
- `--json` for machine-readable output

## Canonical Flow

The canonical read flow remains:

1. `trello_list_boards`
2. `trello_list_lists(board_id)`
3. `trello_list_cards(list_id)`

For command examples and runtime usage, see [entrypoints.md](entrypoints.md).
