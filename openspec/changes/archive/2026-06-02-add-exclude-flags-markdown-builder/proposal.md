## Why

The markdown builder's `--exam` filter only includes chapters matching a specific player name, but there is no way to exclude chapters by player name. Users who want to skip all games by a specific player (e.g., in a large collection) have no option but to manually edit the PGN or use external tools.

## What Changes

- Add `--white-except` flag to exclude chapters where the `White` field matches any of the given values
- Add `--black-except` flag to exclude chapters where the `Black` field matches any of the given values
- Multiple values are separated by pipe `|` (e.g., `--white-except "Carlsen|Nakamura"`)
- Comparison is exact-match after trimming (same semantics as `--exam`)
- Can be combined with `--exam`, `--skip`, `--chapters`, `--chapter-numbers`

## Capabilities

### New Capabilities
- `chapter-exclusion`: Exclude chapters by exact match on White or Black player name

### Modified Capabilities
*(none)*

## Impact

- **`src/builders/markdown_builder/build_markdown.go`**: Add two `flag.String` definitions, add `splitNames()` helper, insert skip logic after tag parsing in the chapter loop
- **`src/builders/markdown_builder/build_markdown_test.go`**: Add tests for `--white-except` and `--black-except` independently and combined with `--exam`
