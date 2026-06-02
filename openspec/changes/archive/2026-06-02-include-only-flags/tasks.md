## 1. Flag definitions

- [x] 1.1 Add `--white-only` string flag in the flag definition block
- [x] 1.2 Add `--black-only` string flag in the flag definition block
- [x] 1.3 Add both flags to the usage/help string

## 2. Core implementation

- [x] 2.1 Add inclusion logic after exclusion (exclusion → inclusion → exam): skip chapter if `--white-only` is set and `rawWhite` doesn't match any pipe-split value
- [x] 2.2 Add inclusion logic for `--black-only`: skip chapter if set and `rawBlack` doesn't match
- [x] 2.3 Both flags together = logical AND: each check independently skips non-matching chapters

## 3. Documentation

- [x] 3.1 Add `--white-only`/`--black-only` to AGENTS.md markdown_builder flags section
- [x] 3.2 Add `--white-only`/`--black-only` to README.md flags list and examples

## 4. Tests

- [x] 4.1 Test `--white-only` single value includes matching chapter
- [x] 4.2 Test `--black-only` single value includes matching chapter
- [x] 4.3 Test multiple pipe-separated values with `--white-only`
- [x] 4.4 Test both flags together (logical AND)
- [x] 4.5 Test pipeline order: exclusion before inclusion before exam
- [x] 4.6 Test non-matching chapters are excluded

## 5. Build and verify

- [x] 5.1 Run `go test ./src/builders/markdown_builder/...` — all tests pass
- [x] 5.2 Run `go build -o bin/build_markdown src/builders/markdown_builder/build_markdown.go` — builds cleanly
