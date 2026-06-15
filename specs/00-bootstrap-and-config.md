# Spec: Bootstrap And Config

## Goal

Provide a stable executable entry point for the MCP server with explicit environment-based configuration and clear startup failures.

## Problem

The coding agent will be asked to rebuild parts of the project. That work needs a small but durable bootstrap contract so MCP wiring and Trello integration do not leak into `main.go`.

## Scope

- Load configuration from environment variables.
- Initialize the Trello client.
- Initialize the MCP server.
- Start stdio serving.
- Fail fast on invalid configuration.

## Functional Requirements

- Required environment variables:
  - `TRELLO_API_KEY`
  - `TRELLO_TOKEN`
- Optional environment variables:
  - `HTTP_TIMEOUT` with default `30s`
- The process must exit non-zero on configuration failure.
- Error messages must clearly distinguish configuration problems from runtime API errors.

## Design Requirements

- Keep `main.go` thin.
- Keep config loading in its own package.
- Keep MCP server construction outside `main.go`.
- Avoid introducing heavy frameworks.

## Acceptance Criteria

- Running the binary with valid environment variables starts the MCP server.
- Missing required credentials produce a readable startup error.
- Configuration loading has focused tests.

## Dependencies

- none
