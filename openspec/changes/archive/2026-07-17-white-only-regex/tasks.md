## 1. Flag Definition

- [x] 1.1 Add `whiteOnlyRegex` flag variable with `flag.String` next to `whiteOnly` in `main()`
- [x] 1.2 Update usage string to include `[-white-only-regex <pattern>]`

## 2. Filter Implementation

- [x] 2.1 Before the `chapterLoop`, split pipe-separated patterns via `splitNames()`, compile each as `*regexp.Regexp`, and collect into a slice; exit with error on invalid pattern
- [x] 2.2 Add filter block after the `whiteOnly` block (around line 826) that skips chapters where `rawWhite` doesn't match any compiled regex

## 3. Tests

- [x] 3.1 Add integration test: `TestMain_WhiteOnlyRegex_Match` — regex matches White field
- [x] 3.2 Add integration test: `TestMain_WhiteOnlyRegex_NoMatch` — regex doesn't match, chapter excluded
- [x] 3.3 Add integration test: `TestMain_WhiteOnlyRegex_Pipe` — multiple pipe-separated patterns
- [x] 3.4 Add integration test: `TestMain_WhiteOnlyRegex_InvalidPattern` — invalid regex causes exit error
- [x] 3.5 Add integration test: `TestMain_WhiteOnlyRegex_WithWhiteExcept` — composition with exclusion filter
- [x] 3.6 Add integration test: `TestMain_WhiteOnlyRegex_CombinedWithWhiteOnly` — both `--white-only` and `--white-only-regex` set

## 4. Verification

- [x] 4.1 Run `go test ./src/builders/markdown_builder/...` and confirm all tests pass
- [x] 4.2 Build binary: `go build -o bin/build_markdown src/builders/markdown_builder/build_markdown.go`
