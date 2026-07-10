## cwctl contacts

Manage contacts

### Options

```
  -h, --help   help for contacts
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
* [cwctl contacts add-labels](cwctl_contacts_add-labels.md)	 - Add labels to a contact (replaces the label set)
* [cwctl contacts contactable-inboxes](cwctl_contacts_contactable-inboxes.md)	 - List the inboxes a contact can be reached through
* [cwctl contacts conversations](cwctl_contacts_conversations.md)	 - List a contact's conversations
* [cwctl contacts create](cwctl_contacts_create.md)	 - Create a contact
* [cwctl contacts create-contact-inbox](cwctl_contacts_create-contact-inbox.md)	 - Attach a contact to an inbox (creates a contact-inbox link)
* [cwctl contacts delete](cwctl_contacts_delete.md)	 - Delete a contact
* [cwctl contacts filter](cwctl_contacts_filter.md)	 - Filter contacts with the query DSL
* [cwctl contacts get](cwctl_contacts_get.md)	 - Get a single contact
* [cwctl contacts labels](cwctl_contacts_labels.md)	 - List a contact's labels
* [cwctl contacts list](cwctl_contacts_list.md)	 - List contacts
* [cwctl contacts merge](cwctl_contacts_merge.md)	 - Merge two contacts (the mergee is deleted)
* [cwctl contacts search](cwctl_contacts_search.md)	 - Search contacts by name, identifier, email, or phone
* [cwctl contacts update](cwctl_contacts_update.md)	 - Update a contact

