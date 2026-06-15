# Spec: MCP CLI With Cobra

## Goal

Refactor the MCP CLI to use `github.com/spf13/cobra` as the command framework.

## Problem

The current CLI already behaves like a real command-line tool, but its command tree is implemented manually with `flag` parsing and custom dispatch. That creates friction for:

- nested command organization
- help and usage output
- persistent flags such as `--json`
- future command expansion as the MCP tool surface grows

The project now has enough commands that the CLI should adopt a standard Go CLI framework.

## Required Library

The implementation must use:

- `github.com/spf13/cobra`

Do not preserve the current ad hoc command parser if it conflicts with a clean Cobra command tree.

## Scope

- Replace manual command dispatch in the CLI with Cobra commands.
- Preserve the current command surface and semantics unless this spec states otherwise.
- Keep the CLI operating through MCP tools, not Trello HTTP calls.
- Preserve support for JSON output and human-readable output.
- Preserve the current in-process MCP connectivity model unless another spec explicitly changes transport.

## Command Tree

The CLI should expose a command tree equivalent to the current behavior:

- `mcp-client tools list`
- `mcp-client boards list`
- `mcp-client boards get --board-id <id>`
- `mcp-client lists list --board-id <id>`
- `mcp-client cards list --list-id <id>`
- `mcp-client cards get --card-id <id>`
- `mcp-client cards create --list-id <id> --name <name> [--desc <text>]`
- `mcp-client cards update --card-id <id> [--name <name>] [--desc <text>] [--due <value>]`
- `mcp-client cards move --card-id <id> --list-id <id>`
- `mcp-client cards comment --card-id <id> --text <text>`
- `mcp-client checklists list --card-id <id>`
- `mcp-client actions list --card-id <id>`
- `mcp-client members list --card-id <id>`

## Functional Requirements

### Help And Usage

- The root command must provide useful help output.
- Each command must provide command-specific usage.
- Required flags must be visible in help output.

### Output Modes

- Support `--json` as a persistent or shared flag where that results in simpler UX.
- Human-readable output must remain the default.
- JSON output must remain machine-readable and based on MCP tool results.

### Error Handling

- Usage errors must return a non-zero exit code.
- Tool errors must remain readable and concise.
- Command handlers should return errors up to a top-level executor rather than terminating deep in the stack.

### Extensibility

- The command tree must make it straightforward to add future MCP tools without rewriting central dispatch logic.
- Shared behavior such as client creation, output selection, and tool execution should remain reusable.

## Design Requirements

- Keep `cmd/mcp-client/main.go` thin.
- Prefer a root Cobra command constructed in an internal CLI package.
- Separate:
  - Cobra command construction
  - MCP client creation
  - tool execution helpers
  - output rendering
- Avoid coupling Cobra command definitions to Trello-specific HTTP details.

## Suggested File Direction

One acceptable direction is:

- `cmd/mcp-client/main.go` -> top-level execution only
- `internal/cli/root.go` -> root Cobra command
- `internal/cli/<area>.go` -> command groups (`boards`, `cards`, etc.)
- `internal/cli/client.go` -> MCP client bootstrap
- `internal/cli/output.go` -> human/JSON rendering

The exact filenames may vary, but the structure should clearly reflect a Cobra-based command tree.

## Testing Requirements

- Update CLI tests to validate Cobra command execution.
- Cover:
  - successful command execution
  - missing required flags
  - `--json` output mode
  - tool error rendering
  - help/usage behavior for at least one command path

## Documentation Requirements

- Update `docs/entrypoints.md` so the CLI section reflects Cobra-based usage and examples.
- Update any other docs that describe the CLI architecture or command parsing approach.

## Acceptance Criteria

- The CLI uses `github.com/spf13/cobra`.
- The manual `flag`-based top-level dispatch is removed.
- The current command surface still works through the new Cobra command tree.
- `--json` still works consistently.
- Tests cover the required CLI behaviors.

## Dependencies

- `25-mcp-cli.md`
