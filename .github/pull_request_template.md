## What & why

<!-- What does this change and why? Link any issue. -->

## Checklist

- [ ] `make verify` passes (check + spec-check + coverage + DoD + judge)
- [ ] New code ships with tests in the same commit (coverage ≥ 80%)
- [ ] Adding a Bot API method? Updated `api-manifest.json` and the group file only —
      no edits to the shared generic builder
- [ ] Docs regenerated if the command surface changed (`make docs-gen`)
- [ ] No secret/token in code, tests, fixtures, or commit messages
