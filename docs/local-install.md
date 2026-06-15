# Local Installation

This document covers local build and installation of the `trello-mcp` binary.

## Build Locally

Clone the repository and build the server binary:

```bash
git clone <your-fork-or-repo-url>
cd trello-mcp
make build
```

This generates the binary at `dist/trello-mcp`.

## Copy The Binary To Another Directory

To build and copy the binary to another location:

```bash
make install DEST=/absolute/path/to/target
```

This copies the server binary to:

```text
/absolute/path/to/target/trello-mcp
```

## Configure Credentials

Required environment variables:

- `TRELLO_API_KEY`
- `TRELLO_TOKEN`

Optional environment variables:

- `HTTP_TIMEOUT` in seconds, default `30`

Example:

```bash
export TRELLO_API_KEY="your-key"
export TRELLO_TOKEN="your-token"
export HTTP_TIMEOUT="30"
```

## Register In An MCP Host

This server uses stdio transport. Your MCP host should launch the binary as a subprocess.

Example shape:

```json
{
  "command": "/absolute/path/to/trello-mcp",
  "env": {
    "TRELLO_API_KEY": "your-key",
    "TRELLO_TOKEN": "your-token",
    "HTTP_TIMEOUT": "30"
  }
}
```

Exact configuration depends on the MCP host.
