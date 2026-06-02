## 1. Flag definition

- [x] 1.1 Add `--split-by-chapters int` flag in the flag definition block
- [x] 1.2 Validate the value: ignore (treat as no split) if ≤ 0
- [x] 1.3 Add the flag to the usage/help string

## 2. Helper function

- [x] 2.1 Add `splitOutputPath(basePath string, part int) string` — inserts `_N` before extension (e.g., `theory.md`, `1` → `theory_1.md`)

## 3. Core implementation

- [x] 3.1 Replace `output := ""` with `[]string` per-chapter output
- [x] 3.2 Use per-chapter `strings.Builder` to accumulate each chapter's content
- [x] 3.3 Append each chapter's builder content to `chaptersOutput` slice
- [x] 3.4 After loop, write single file if no split needed
- [x] 3.5 If splitting, iterate through `chaptersOutput` in threshold-sized groups
- [x] 3.6 Replace final `os.WriteFile` with split-aware write logic

## 4. Documentation

- [x] 4.1 Add `--split-by-chapters` flag to the markdown_builder section in AGENTS.md
- [x] 4.2 Add `--split-by-chapters` flag to README.md flags list and examples

## 5. Tests

- [x] 4.1 Test `splitOutputPath` helper (no extension, .md extension, path with directories)
- [x] 4.2 Test `--split-by-chapters` splits into correct number of files
- [x] 4.3 Test each split file contains the expected chapters
- [x] 4.4 Test zero/negative split value produces single unsplit file
- [x] 4.5 Test split composes with `--skip`
- [x] 4.6 Test split composes with `--exam` + `--white-except`
- [x] 4.7 Test split composes with `--chapter-numbers`

## 6. Build and verify

- [x] 6.1 Run `go test ./src/builders/markdown_builder/...` — all tests pass
- [x] 6.2 Run `go build -o bin/build_markdown src/builders/markdown_builder/build_markdown.go` — builds cleanly
