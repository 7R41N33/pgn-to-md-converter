## Context

`build_markdown.go` currently accumulates all processed chapters into a single `output` string and writes it to one file via `--out`. For large PGN files (100+ chapters), a single output file is unwieldy. The flag `--split-by-chapters <N>` splits the result into multiple files, each containing at most N chapters.

## Goals / Non-Goals

**Goals:**
- Add `--split-by-chapters <int>` flag to markdown_builder
- Split output into `base_N.ext` files (e.g., `theory_1.md`, `theory_2.md`)
- Apply split after all filters (skip, chapter-numbers, exam, white/black-except)
- Each file forms a valid contiguous segment of the chapter sequence
- If total chapters ≤ split value, emit single file without suffix (no-op)

**Non-Goals:**
- Parallel/concurrent processing of chapters
- Splitting within a single chapter
- Splitting based on byte size or page count
- Re-splitting existing files (only during generation)

## Decisions

1. **Filename scheme:** Strip extension from `--out`, append `_N` before extension.  
   `theory.md` → `theory_1.md`, `theory_2.md`, ...  
   `.md` is the only expected extension; if absent, append `_1` etc. directly.

2. **Streaming approach:** Write files incrementally during the chapter loop rather than building a giant string and splitting post-hoc. A `strings.Builder` accumulates chapters until it hits the split threshold, then flushes to disk and resets.

3. **Single-part no-suffix rule:** When the total number of chapters (after filters) ≤ `--split-by-chapters`, no suffix is added — output is identical to running without the flag.

4. **Validation:** If `--split-by-chapters` ≤ 0, ignore the flag (or treat as 1). No error.

## Risks / Trade-offs

- [Risk] File overwrite if `base_1.md` already exists → use `os.OpenFile` with O_CREATE|O_EXCL or simply overwrite (consistent with single-file behavior). Decision: overwrite (same as current single-file behavior).
- [Risk] Very large individual files still possible with high `--split-by-chapters` value → acceptable, user controls the threshold.
