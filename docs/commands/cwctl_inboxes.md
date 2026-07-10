## cwctl inboxes

Manage inboxes and their members

### Options

```
  -h, --help   help for inboxes
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
* [cwctl inboxes add-members](cwctl_inboxes_add-members.md)	 - Add agents to an inbox
* [cwctl inboxes agent-bot](cwctl_inboxes_agent-bot.md)	 - Show the agent bot attached to an inbox
* [cwctl inboxes create](cwctl_inboxes_create.md)	 - Create a inboxe
* [cwctl inboxes get](cwctl_inboxes_get.md)	 - Get a single inboxe
* [cwctl inboxes list](cwctl_inboxes_list.md)	 - List inboxes
* [cwctl inboxes members](cwctl_inboxes_members.md)	 - List the agents in an inbox
* [cwctl inboxes remove-members](cwctl_inboxes_remove-members.md)	 - Remove agents from an inbox
* [cwctl inboxes set-agent-bot](cwctl_inboxes_set-agent-bot.md)	 - Attach an agent bot to an inbox (0 detaches)
* [cwctl inboxes update](cwctl_inboxes_update.md)	 - Update a inboxe
* [cwctl inboxes update-members](cwctl_inboxes_update-members.md)	 - Replace an inbox's agents

