---
name: cwctl-cli
description: Operate Chatwoot from the terminal with the `cwctl` CLI — list/read/reply to conversations (attachments, private notes), assign and resolve, manage contacts (search/filter/merge/labels), agents, teams, inboxes and members, labels, canned responses, automation rules, webhooks, help-center portals, pull analytics/reports, provision accounts/users via the platform API, and drive the public client API. Use whenever the user wants to read or answer customer conversations, triage or resolve tickets, look up or edit contacts, pull support metrics, or automate any Chatwoot workflow. Prefer it over raw curl to the Chatwoot API.
version: 0.1.0
homepage: https://github.com/jjuanrivvera/cwctl
license: MIT
allowed-tools: Bash(cwctl:*)
metadata: {"openclaw":{"category":"customer-support","emoji":"💬","requires":{"bins":["cwctl"],"env":["CWCTL_API_KEY"]},"install":[{"kind":"brew","formula":"jjuanrivvera/cwctl/cwctl-cli","bins":["cwctl"]},{"kind":"go","package":"github.com/jjuanrivvera/cwctl/cmd/cwctl@latest","bins":["cwctl"]}]}}
---

# cwctl — Chatwoot CLI

## Prerequisites

- `cwctl` on PATH (`brew install jjuanrivvera/cwctl/cwctl-cli` or
  `go install github.com/jjuanrivvera/cwctl/cmd/cwctl@latest`).
- A configured profile: `cwctl auth status` must succeed. If it doesn't, the human runs
  `cwctl auth login` (token goes to the OS keyring; `CWCTL_API_KEY` also works for CI).
- `cwctl platform …` needs a platform app token (`auth login --platform-token`);
  `cwctl client …` needs no token at all (public API).

## Golden rules

1. **Prefer cwctl over raw curl** — it handles auth, account scoping, pagination,
   retries, and rate limits; `cwctl api METHOD PATH` is the escape hatch if an
   endpoint is somehow missing.
2. **Customer-visible writes are real.** `messages create` posts to a real person
   unless `--private` (internal note). When drafting, show the human the text first
   or use `--dry-run`.
3. **Never guess ids** — resolve them: `conversations list`, `contacts search --q`,
   `agents list -o json`.
4. **Destructive verbs are blocked for you** if the operator installed
   `cwctl agent guard` (deletes, `contacts merge`). Don't try to work around it;
   ask the human.
5. **JSON for parsing, table for humans.** `-o json` never drops fields; `--jq`
   extracts inline.

## Workflow: auth → discover → act → verify

```bash
cwctl auth status                                  # who am I / which account
cwctl conversations list --status open --assignee-type me
cwctl messages list 42                             # read before replying
cwctl messages create 42 --content "Hola Ana, ya lo reviso."
cwctl conversations toggle-status 42 --status resolved
cwctl conversations get 42 -o json | jq .status    # verify
```

## Cheatsheet

```bash
# conversations
cwctl conversations list --status open --labels vip
cwctl conversations meta                            # open/unassigned counts
cwctl conversations assign 42 --assignee-id 7       # or --team-id 2
cwctl conversations toggle-priority 42 --priority urgent
cwctl conversations add-labels 42 --labels billing  # REPLACES the label set

# messages
cwctl messages create 42 --content "..." [--private] [--attachment ./f.pdf]

# contacts
cwctl contacts search --q "+57300"
cwctl contacts filter --payload '[{"attribute_key":"country_code","filter_operator":"equal_to","values":["CO"]}]'
cwctl contacts update 12 --email ana@example.com

# team / setup
cwctl agents list · cwctl teams members 3 · cwctl inboxes list
cwctl canned-responses list · cwctl labels list · cwctl webhooks list

# analytics
cwctl reports summary --since 2026-06-01 --until 2026-07-01
cwctl reports agent-summary --since 2026-06-01 --business-hours

# several instances
cwctl --profile staging conversations list
```

## Troubleshooting

- `401 … run cwctl auth login` → token missing/expired; the human re-runs login.
- `403` → the token's agent role lacks access (admin-only endpoint, or Enterprise
  feature like audit logs / SLA).
- `404 … verify the id` → wrong id OR wrong account: check `--account-id` /
  `cwctl auth status`.
- `platform API needs a platform app token` → `auth login --platform-token <t>`.
- Rate limited (429) → cwctl backs off automatically; lower `--rps` for bulk loops.
- Anything unclear: `cwctl doctor` first, `--dry-run` to see the exact request.
