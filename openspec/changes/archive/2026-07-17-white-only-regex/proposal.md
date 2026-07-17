## Why

The `--white-only` flag filters chapters by exact match on the `White` PGN field. This is too rigid for cases where White names contain suffixes, alternate spellings, or tournament metadata (e.g., `"Carlsen, Magnus"` vs `"Carlsen"`). A regex-based filter gives users flexible pattern matching — they can match by surname, prefix, or complex patterns without knowing the exact string.

## What Changes

- Add `--white-only-regex` CLI flag to `build_markdown.go`
- Accept pipe-separated regex patterns (same convention as `--white-only`)
- Include a chapter if its `White` field matches **any** of the given regex patterns
- Exit with a clear error if any regex pattern is invalid
- Pre-compile all regex patterns once before the chapter loop (performance)
- Update usage/help string

No breaking changes. Existing `--white-only` behavior is unchanged.

## Capabilities

### New Capabilities

_(none — this extends an existing capability)_

### Modified Capabilities

- `chapter-inclusion`: Add regex-based white-only filtering as an additional include mode alongside exact-match `--white-only`. The new flag uses `regexp.MatchString` instead of `==` comparison.

## Impact

- **One file changed**: `src/builders/markdown_builder/build_markdown.go`
- No new dependencies (Go standard library `regexp` already imported)
- No API surface changes (CLI flag only)
- No schema or storage changes
