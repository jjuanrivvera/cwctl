## cwctl portals

Manage help-center portals, articles, and categories

### Options

```
  -h, --help   help for portals
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
* [cwctl portals create](cwctl_portals_create.md)	 - Create a portal
* [cwctl portals create-article](cwctl_portals_create-article.md)	 - Add an article to a portal
* [cwctl portals create-category](cwctl_portals_create-category.md)	 - Add a category to a portal
* [cwctl portals list](cwctl_portals_list.md)	 - List portals
* [cwctl portals update](cwctl_portals_update.md)	 - Update a portal

