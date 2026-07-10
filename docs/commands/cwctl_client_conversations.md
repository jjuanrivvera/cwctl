## cwctl client conversations

Public (contact-facing) conversation endpoints

### Options

```
  -h, --help   help for conversations
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

* [cwctl client](cwctl_client.md)	 - Public client API (inbox/contact/conversation flows) — no token required
* [cwctl client conversations create](cwctl_client_conversations_create.md)	 - Start a conversation as a contact
* [cwctl client conversations get](cwctl_client_conversations_get.md)	 - Get one of the contact's conversations
* [cwctl client conversations list](cwctl_client_conversations_list.md)	 - List a contact's conversations in a public inbox
* [cwctl client conversations resolve](cwctl_client_conversations_resolve.md)	 - Resolve a conversation as the contact (toggle_status)
* [cwctl client conversations toggle-typing](cwctl_client_conversations_toggle-typing.md)	 - Flip the contact-side typing indicator
* [cwctl client conversations update-last-seen](cwctl_client_conversations_update-last-seen.md)	 - Mark the conversation read up to now (contact side)

