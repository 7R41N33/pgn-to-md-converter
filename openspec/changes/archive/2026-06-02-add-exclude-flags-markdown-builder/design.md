## Context

The markdown builder's chapter filtering currently supports `--skip`, `--chapters`, `--chapter-numbers`, and `--exam` (include-only by player name). There is no exclude-by-player mechanism. The chapter loop processes `games` from `splitGames()`, parses tags via `parseGameTags()`, and applies filters sequentially.

## Goals / Non-Goals

**Goals:**
- Add `--white-except` flag to exclude chapters by exact White field match
- Add `--black-except` flag to exclude chapters by exact Black field match
- Same exact-match semantics as `--exam` (case-sensitive, trimmed)
- Compose cleanly with all existing flags (`--exam`, `--skip`, `--chapters`, etc.)

**Non-Goals:**
- Regex/glob/wildcard matching — exact match only, same as `--exam`
- Case-insensitive matching
- Exclusion by other fields (Event, Site, Date, etc.)
- Partial or substring matching

## Decisions

### 1. Pipe `|` as multi-value delimiter
**Choice:** Split the flag value on `|` to support multiple names: `--white-except "Carlsen|Nakamura"`.
**Rationale:** Comma `,` can appear in PGN names (`"Carlsen, Magnus"`). Pipe never appears in player names and is a standard alternative delimiter in CLI tools.
**Alternatives considered:** comma (breaks on `LastName, FirstName`), repeated flags (`flag.Var` requires custom type — more code for same UX), semicolon (less conventional).

### 2. Exact string match (same as --exam)
**Choice:** `strings.TrimSpace(value) == name` for each value in the pipe-split list.
**Rationale:** Consistency with `--exam`. Each individual name is compared exactly after trimming.
**Alternatives considered:** regex, prefix, substring — would break least-surprise principle.

### 3. Apply exclusion before --exam
**Choice:** Insert the except check immediately after tag parsing, before the `--exam` include check.
**Rationale:** Exclusions should act as a pre-filter. The logical order is: skip excluded chapters → apply include filter. This also means `--exam` only sees chapters that already passed exclusion.
**Pipeline order in chapter loop:**
1. `--skip` / `--chapters` / `--chapter-numbers` (position-based)
2. `--white-except` / `--black-except` (exclusion by player)
3. `--exam` (inclusion by player)

### 3. `rawWhite` / `rawBlack` before translation
**Choice:** Use `rawWhite` and `rawBlack` from `tags` before `translateTags()` is called, matching how `--exam` accesses them.
**Rationale:** `translateTags()` is currently a no-op and may eventually translate names. Excluding by translated names would be confusing. Using raw PGN fields keeps behavior predictable.

## Risks / Trade-offs

- [Risk] User misspells player name → exclusion silently fails (no match). Mitigation: same as `--exam` — exact match is predictable; case sensitivity is documented.
- [Risk] Combined with `--exam` in confusing ways (e.g., `--white-except "Carlsen" --exam "Carlsen"` → zero chapters). Mitigation: user responsibility; no guard added to avoid over-engineering.
