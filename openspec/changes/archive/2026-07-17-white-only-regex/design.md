## Context

The markdown builder (`src/builders/markdown_builder/build_markdown.go`) accepts CLI flags to filter which PGN chapters appear in the output. The current `--white-only` flag does exact string matching (`==`). A regex variant is needed for flexible pattern matching.

The flag parsing lives in `main()`, and the filtering happens in the `chapterLoop` (lines 772–943). Each filter is an independent `if` block with `continue chapterLoop` for exclusion.

## Goals / Non-Goals

**Goals:**
- Add `--white-only-regex` CLI flag that filters chapters by regex match on `White` field
- Support pipe-separated multiple regex patterns (same convention as `--white-only`)
- Exit with clear error on invalid regex
- Pre-compile regexes before the chapter loop for performance
- Place filter between `--white-only` and `--black-only` in the filter chain

**Non-Goals:**
- No `--black-only-regex` (can be added later in the same pattern)
- No regex caching or persistence across invocations
- No regex-based exclude flags (same reason)

## Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| **Matching semantics** | `regexp.MatchString` (Go stdlib) | The user confirmed regex, not substring contains. No new deps needed. |
| **Pre-compilation** | Compile all patterns once before `chapterLoop` | Regex compilation is expensive; doing it per-chapter would be wasteful. |
| **Invalid regex handling** | `os.Exit(1)` with error message | Consistent with existing error handling in the codebase (e.g., missing src/out flags). |
| **Filter position** | After `--white-only`, before `--black-only` | Sibling filter on the `White` field. Both are inclusion filters — if both set, a chapter must pass both (AND). |
| **Pipe separation** | Reuse `splitNames()` | Zero new code for parsing; user already knows this convention from `--white-only`. |

## Risks / Trade-offs

- **[AND semantics with --white-only]** If both `--white-only` and `--white-only-regex` are set, a chapter must pass both. This could surprise users who expect regex to override exact. Mitigation: document clearly in help text, and consider that in practice users will use one or the other.
- **[Case sensitivity]** Go regex is case-sensitive by default. Users who need case-insensitive matching must write `(?i)pattern`. This is standard regex behavior and matches existing case-sensitive `--white-only`.
- **[Edge case: empty string]** `regexp.Compile("")` matches everything. We skip empty strings from the pipe-split list to avoid this surprise.
