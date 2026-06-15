# Entrypoints

This repository exposes two executable surfaces:

1. the production MCP server
2. the local MCP CLI

## MCP Server

```bash
export $(grep -v '^#' .env | xargs)
go run main.go
```

Starts a stdio-based MCP server. An MCP host such as Claude Desktop or Cursor should launch this binary as a subprocess.

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `TRELLO_API_KEY` | yes | — | Trello API key |
| `TRELLO_TOKEN` | yes | — | Trello API token |
| `HTTP_TIMEOUT` | no | `30` | HTTP timeout in seconds |

### Tool Surface

#### Navigation

| Tool | Description | Arguments |
|------|-------------|-----------|
| `trello_list_boards` | List visible Trello boards | none |
| `trello_list_lists` | List lists inside a board | `board_id` |
| `trello_list_cards` | List cards inside a list (extended: due, labels) | `list_id` |

#### Details

| Tool | Description | Arguments |
|------|-------------|-----------|
| `trello_get_board` | Get one board | `board_id` |
| `trello_get_card` | Get one card (extended: due, labels, members) | `card_id` |

#### Expansion Reads

| Tool | Description | Arguments |
|------|-------------|-----------|
| `trello_list_checklists` | List checklists on a card | `card_id` |
| `trello_list_card_actions` | List actions on a card | `card_id` |
| `trello_list_card_members` | List members on a card | `card_id` |

#### Organization Reads

| Tool | Description | Arguments |
|------|-------------|-----------|
| `trello_list_organizations` | List organizations (workspaces) | none |

#### Label Reads

| Tool | Description | Arguments |
|------|-------------|-----------|
| `trello_list_board_labels` | List labels defined on a board | `board_id` |
| `trello_list_card_labels` | List labels applied to a card | `card_id` |

#### Card Mutations

| Tool | Description | Arguments |
|------|-------------|-----------|
| `trello_create_card` | Create a card | `list_id`, `name`, `desc` (optional), `due` (optional, RFC3339) |
| `trello_update_card` | Update a card | `card_id`, plus at least one of `name`, `desc`, `due` |
| `trello_move_card` | Move a card to another list | `card_id`, `list_id` |
| `trello_add_comment` | Add a comment to a card | `card_id`, `text` |

#### Card Archiving

| Tool | Description | Arguments |
|------|-------------|-----------|
| `trello_archive_card` | Archive (close) a card | `card_id` |
| `trello_unarchive_card` | Unarchive (re-open) a card | `card_id` |
| `trello_delete_card` | Permanently delete a card | `card_id` |

#### Search

| Tool | Description | Arguments |
|------|-------------|-----------|
| `trello_search_cards` | Search cards by text | `query`, `board_id` (optional) |

#### Label Mutations

| Tool | Description | Arguments |
|------|-------------|-----------|
| `trello_add_card_label` | Attach a label to a card | `card_id`, `label_id` |
| `trello_remove_card_label` | Remove a label from a card | `card_id`, `label_id` |

#### Card Member Mutations

| Tool | Description | Arguments |
|------|-------------|-----------|
| `trello_add_card_member` | Assign a member to a card | `card_id`, `member_id` |
| `trello_remove_card_member` | Remove a member from a card | `card_id`, `member_id` |

#### Board and List Mutations

| Tool | Description | Arguments |
|------|-------------|-----------|
| `trello_create_board` | Create a new board | `name`, `default_lists` (optional) |
| `trello_create_list` | Create a new list | `name`, `board_id` |

### Canonical Read Flow

1. `trello_list_boards`
2. `trello_list_lists(board_id)`
3. `trello_list_cards(list_id)`

## MCP CLI

```bash
export $(grep -v '^#' .env | xargs)
go run ./cmd/mcp-client/ <command> <action> [flags] [--json]
```

The CLI is the local operational/manual-testing interface for the MCP server contract. It currently connects to the server in-process.

### Current Command Surface

#### Read Commands

| Command | Flags |
|---------|-------|
| `tools list` | — |
| `boards list` | — |
| `boards get` | `--board-id` |
| `boards create` | `--name`, optional `--no-default-lists` |
| `lists list` | `--board-id` |
| `lists create` | `--name`, `--board-id` |
| `cards list` | `--list-id` |
| `cards get` | `--card-id` |
| `checklists list` | `--card-id` |
| `actions list` | `--card-id` |
| `members list` | `--card-id` |
| `members add` | `--card-id`, `--member-id` |
| `members remove` | `--card-id`, `--member-id` |
| `labels board-list` | `--board-id` |
| `labels card-list` | `--card-id` |
| `orgs list` | — |

#### Write Commands

| Command | Flags |
|---------|-------|
| `cards create` | `--list-id`, `--name`, optional `--desc`, optional `--due` |
| `cards update` | `--card-id`, plus at least one of `--name`, `--desc`, `--due` |
| `cards move` | `--card-id`, `--list-id` |
| `cards comment` | `--card-id`, `--text` |
| `cards archive` | `--card-id` |
| `cards unarchive` | `--card-id` |
| `cards delete` | `--card-id` |
| `cards search` | `--query`, optional `--board-id` |
| `labels add` | `--card-id`, `--label-id` |
| `labels remove` | `--card-id`, `--label-id` |

### Output

- default: human-readable output
- `--json`: machine-readable JSON output

Example:

```bash
go run ./cmd/mcp-client/ boards list --json
```

### CLI Implementation Notes

The CLI is implemented with `github.com/spf13/cobra` and keeps `--json` as a persistent root flag.

Current structure:

- `cmd/mcp-client/main.go`: thin process entrypoint
- `internal/cli/root.go`: root Cobra command and shared flags
- `internal/cli/*.go`: command groups, tool execution, and output rendering

The design remains intentionally simple: commands operate through MCP tools and the local client currently connects to the server in-process.

## Build

```bash
make build
./dist/trello-mcp
```
