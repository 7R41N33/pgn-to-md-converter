## 1. Add flag definitions

- [x] 1.1 Add `--white-except` string flag in the flag definition block
- [x] 1.2 Add `--black-except` string flag in the flag definition block
- [x] 1.3 Add both flags to the usage/help string

## 2. Helper function

- [x] 2.1 Add `splitNames(s string) []string` — split on `|`, trim spaces, filter empties

## 3. Implement exclusion logic

- [x] 3.1 After `tags := parseGameTags(g)`, extract `rawWhite` and `rawBlack` (trimmed)
- [x] 3.2 Skip chapter if `whiteExcept` is set and `rawWhite` matches any of the pipe-split values
- [x] 3.3 Skip chapter if `blackExcept` is set and `rawBlack` matches any of the pipe-split values
- [x] 3.4 Ensure exclusion runs before `--exam` filter

## 4. Tests

- [x] 4.1 Test single value excludes matching chapters
- [x] 4.2 Test multiple pipe-separated values exclude all matching chapters
- [x] 4.3 Test non-matching chapters are kept
- [x] 4.4 Test both `--white-except` and `--black-except` together
- [x] 4.5 Test combination with `--exam` (exclusion before inclusion)
- [x] 4.6 Test whitespace trimming in flag values and PGN fields
- [x] 4.7 Test `splitNames` helper directly (empty string, single value, multiple values, whitespace)

## 5. Build and verify

- [x] 5.1 Run `go test ./src/builders/markdown_builder/...` — all tests pass
- [x] 5.2 Run `go build -o bin/build_markdown src/builders/markdown_builder/build_markdown.go` — builds cleanly
