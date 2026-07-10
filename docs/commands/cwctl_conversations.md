## cwctl conversations

Manage conversations

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

* [cwctl](cwctl.md)	 - A fast, scriptable CLI for the full Chatwoot API
* [cwctl conversations add-labels](cwctl_conversations_add-labels.md)	 - Add labels to a conversation (replaces the label set)
* [cwctl conversations assign](cwctl_conversations_assign.md)	 - Assign a conversation to an agent or a team
* [cwctl conversations create](cwctl_conversations_create.md)	 - Create a conversation
* [cwctl conversations filter](cwctl_conversations_filter.md)	 - Filter conversations with the query DSL
* [cwctl conversations get](cwctl_conversations_get.md)	 - Get a single conversation
* [cwctl conversations labels](cwctl_conversations_labels.md)	 - List a conversation's labels
* [cwctl conversations list](cwctl_conversations_list.md)	 - List conversations
* [cwctl conversations meta](cwctl_conversations_meta.md)	 - Conversation counts (mine, unassigned, assigned, all)
* [cwctl conversations reporting-events](cwctl_conversations_reporting-events.md)	 - List a conversation's reporting events (first response, resolved, …)
* [cwctl conversations set-custom-attributes](cwctl_conversations_set-custom-attributes.md)	 - Set custom attributes on a conversation
* [cwctl conversations toggle-priority](cwctl_conversations_toggle-priority.md)	 - Change a conversation's priority
* [cwctl conversations toggle-status](cwctl_conversations_toggle-status.md)	 - Change a conversation's status (open/resolved/pending/snoozed)
* [cwctl conversations toggle-typing](cwctl_conversations_toggle-typing.md)	 - Flip the typing indicator on or off
* [cwctl conversations update](cwctl_conversations_update.md)	 - Update a conversation

