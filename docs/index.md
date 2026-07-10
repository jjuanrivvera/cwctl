# cwctl

A fast, scriptable CLI for the full [Chatwoot](https://www.chatwoot.com) API — every
documented operation (144/144) across the application, platform, and public client
APIs, with named profiles, keyring token storage, and an MCP server for AI agents.

```console
$ cwctl conversations list --status open
ID  INBOX_ID  STATUS  PRIORITY
42  3         open    high

$ cwctl messages create 42 --content "On it — checking now."
$ cwctl conversations toggle-status 42 --status resolved
```

## Why cwctl

- **The whole API.** The command surface is pinned to a manifest derived from
  Chatwoot's own OpenAPI spec, and CI fails if anything drifts or goes missing.
- **Several instances, one tool.** Profiles bundle base URL + account + tokens;
  `--profile staging` switches everything at once.
- **Works headless.** OS keyring when there is one, an AES-256-GCM encrypted file
  when there isn't (VPS, containers, CI).
- **Made for pipes.** json/yaml/csv/id output, `--jq`, `--columns`, `--all`
  pagination, `--dry-run` curls.
- **Agent-ready.** `cwctl mcp` exposes annotated tools; `cwctl agent guard`
  generates guardrails that hard-block irreversible operations.

## Install

```bash
brew install jjuanrivvera/cwctl/cwctl-cli          # macOS / Linux
scoop install cwctl                                 # Windows (bucket: scoop-cwctl)
go install github.com/jjuanrivvera/cwctl/cmd/cwctl@latest
```

Then: [Getting started](getting-started.md).
