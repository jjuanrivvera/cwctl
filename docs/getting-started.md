# Getting started

## 1. Connect

```bash
cwctl auth login
```

You are asked for:

- **Base URL** — your instance root (`https://app.chatwoot.com`, or your self-hosted
  domain).
- **api_access_token** — from Chatwoot: Profile Settings → Access Token. Input is
  hidden and the token goes to your OS keyring, never to a file in plaintext.

The token is verified against `GET /api/v1/profile` and the account is captured (a
multi-account token prompts you to pick one). `cwctl init` runs the same flow plus a
smoke check.

## 2. Look around

```bash
cwctl auth status            # who am I, which account, is the token valid
cwctl inboxes list
cwctl conversations meta     # open/unassigned counts
cwctl conversations list --status open --assignee-type me
```

## 3. Work a conversation

```bash
cwctl messages list 42
cwctl messages create 42 --content "On it."
cwctl messages create 42 --content "internal note" --private
cwctl messages create 42 --content "see attached" --attachment ./invoice.pdf
cwctl conversations assign 42 --assignee-id 7
cwctl conversations toggle-status 42 --status resolved
```

## 4. Script it

```bash
cwctl contacts search --q ana -o json | jq '.[0].id'
cwctl conversations list --all -o csv > backlog.csv
cwctl labels list -o id | xargs -n1 -I{} cwctl labels get {}
```

Every command honors `--dry-run` (prints the exact curl, token redacted), so you can
always see what would happen before it does.

## Environment variables

| Variable | Effect |
|---|---|
| `CWCTL_API_KEY` | token override (CI) — beats the keyring |
| `CWCTL_PLATFORM_TOKEN` | platform app token override |
| `CWCTL_BASE_URL` / `CWCTL_ACCOUNT_ID` | instance/account override |
| `CWCTL_PROFILE` | profile selection per shell |
| `CWCTL_KEYRING_PASSWORD` | key for the encrypted-file fallback on headless hosts |
| `NO_COLOR` | disable table colors |
