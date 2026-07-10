# Authentication & profiles

## Token classes

Chatwoot authenticates everything with one header (`api_access_token`) but three
token classes exist:

| Class | Used by | Stored as |
|---|---|---|
| user token | application API (`cwctl <resource> …`) | keyring key `<profile>` |
| platform app token | `cwctl platform …` (self-hosted provisioning) | keyring key `<profile>/platform` |
| none | `cwctl client …` (public contact-facing API), `csat` | — |

cwctl picks the class from the path automatically — including for the raw
`cwctl api METHOD PATH` escape hatch.

```bash
cwctl auth login                                   # user token (verified live)
cwctl auth login --platform-token <token>          # add a platform token
cwctl auth status                                  # identity + validity + backend
cwctl auth logout                                  # remove both tokens
```

## Where tokens live

1. **OS keyring** — macOS Keychain, Linux Secret Service, Windows Credential Manager.
2. **Encrypted file fallback** — when no keyring is reachable (VPS, container), an
   AES-256-GCM `credentials.enc` under the config dir. Set `CWCTL_KEYRING_PASSWORD`
   to key it (scrypt-derived); without it the key is host-bound, which resists casual
   copying but is not a hard boundary.
3. **Env override** — `CWCTL_API_KEY` beats everything (CI).

Non-secret settings live in `~/.cwctl-cli/config.yaml`
(or `$XDG_CONFIG_HOME/cwctl/config.yaml`), written 0600.

## Profiles

A profile = base URL + account id + tokens. Use them for staging vs production, or
two companies on one Cloud login:

```bash
cwctl --profile staging auth login
cwctl config list-profiles
cwctl config use staging
cwctl --profile prod conversations list     # one-off override
CWCTL_PROFILE=staging cwctl doctor          # per-shell
```

`cwctl config set base_url|account_id|rps <value>` edits the active profile;
`--account-id` overrides the account for a single invocation.

## Diagnosing

```bash
cwctl doctor          # config, profile, base URL, token presence + validity, account
cwctl doctor --json   # for scripts; exits non-zero on any failure
```
