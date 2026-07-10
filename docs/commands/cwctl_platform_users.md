## cwctl platform users

Provision users (platform app token)

### Options

```
  -h, --help   help for users
```

### Options inherited from parent commands

```
      --account-id string   override the profile's account id for this invocation
      --all                 fetch all pages (list commands)
      --base-url string     override the instance base URL
      --columns strings     comma-separated columns to show
      --dry-run             print the equivalent curl and make no request
      --filter strings      client-side field=value filters (list commands)
      --jq string           gojq expression applied to the response before rendering
      --limit int           max items to output, applied client-side (list commands)
      --no-color            disable colored output
  -o, --output string       output format: table|json|yaml|csv|id
      --page int            page number to fetch (list commands; Chatwoot pages are server-sized)
      --profile string      named profile to use (instance + account + token)
      --quiet               suppress non-essential chatter
      --rps rps             max requests per second (default 5; also per-profile rps in config)
      --show-token          reveal the API token in dry-run output
      --sort string         sort field, prefix with - for descending (where the API supports it)
  -v, --verbose             verbose request logging (stderr)
```

### SEE ALSO

* [cwctl platform](cwctl_platform.md)	 - Platform API (accounts, users, agent bots) — needs a platform app token
* [cwctl platform users create](cwctl_platform_users_create.md)	 - Create a user
* [cwctl platform users delete](cwctl_platform_users_delete.md)	 - Delete a user
* [cwctl platform users get](cwctl_platform_users_get.md)	 - Get a single user
* [cwctl platform users sso-link](cwctl_platform_users_sso-link.md)	 - Get a user's one-time SSO login URL
* [cwctl platform users update](cwctl_platform_users_update.md)	 - Update a user

