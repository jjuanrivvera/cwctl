## cwctl client

Public client API (inbox/contact/conversation flows) — no token required

### Options

```
  -h, --help   help for client
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

* [cwctl](cwctl.md)	 - A fast, scriptable CLI for the full Chatwoot API
* [cwctl client contacts](cwctl_client_contacts.md)	 - Public (contact-facing) contact endpoints
* [cwctl client conversations](cwctl_client_conversations.md)	 - Public (contact-facing) conversation endpoints
* [cwctl client inbox](cwctl_client_inbox.md)	 - Public inbox details
* [cwctl client messages](cwctl_client_messages.md)	 - Public (contact-facing) message endpoints

