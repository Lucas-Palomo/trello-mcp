# Roadmap

## Goal

Track follow-up evolution after the current Trello MCP surface became viable for real use.

## Current State

The repository already includes:

- navigation across boards, lists, and cards
- board and card detail reads
- checklist, action, member, label, and organization reads
- card create/update/move/comment/archive/unarchive/delete
- label and member assignment mutations
- board and list creation
- card search
- a Cobra-based CLI that exercises the same MCP contract locally

## Near-Term Priorities

### 1. Harden Test Coverage

- close any gaps between the newest tools and their focused tests
- keep client, server, and CLI coverage aligned when new tools are added
- make sandbox-sensitive tests easier to run in restricted environments when possible

### 2. Tighten Contract Consistency

- keep `specs/`, `docs/`, and exposed MCP tools in sync
- preserve small, task-oriented result envelopes
- standardize error behavior where newer tools expanded the surface

### 3. Improve Runtime Integration

- keep the current in-process CLI connectivity model easy to replace
- evaluate a future CLI mode that spawns the stdio server or supports another MCP transport

### 4. Expand Carefully

Potential future capabilities should stay consistent with the project style:

- more card lifecycle operations if they clearly improve MCP workflows
- additional board context reads if they reduce multi-step navigation cost
- narrow write operations only when they can be validated and documented cleanly

## Design Notes

- keep the MCP surface task-oriented and small
- prefer explicit contracts over broad Trello mirroring
- evolve the current architecture instead of bypassing it with ad hoc integrations
