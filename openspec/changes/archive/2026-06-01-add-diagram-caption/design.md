## Context

The `fenlib` package renders chess diagrams as RGBA images using a 60px cell size and a 45px border, producing a 570×570px image. Text is currently limited to coordinate labels (a–h, 1–8) rendered with a 5×7 bitmap font. The `GenerateBoard` function is the single render entry point; all other functions (`GenerateDiagram`, `GenerateDiagramBase64`, `GenerateDiagramToWriter`) call it.

Captions need to handle arbitrary Russian text (as used in training materials) with automatic word-wrapping and centering. The project avoids external dependencies where possible.

## Goals / Non-Goals

**Goals:**
- Add an optional caption parameter to `GenerateBoard` and all forwarding functions
- Render caption text centered below the board diagram
- Word-wrap text that exceeds the board width
- Support Latin letters, digits, punctuation, and Russian Cyrillic (А–Я, а–я, Ё, ё)
- Use Go Mono through existing `golang.org/x/image` dependency
- Backward compatible — empty caption produces identical output

**Non-Goals:**
- Variable font sizes, bold/italic, or rich text formatting
- Caption above the board or on the border
- Unicode beyond Cyrillic (e.g., Chinese, Arabic)
- Animated or interactive diagrams

## Decisions

### 1. Go Mono via golang.org/x/image (TrueType, WGL4)
**Choice:** Use the Go Mono font from `golang.org/x/image/font/gofont/gomono` rendered via `font/opentype`.
**Rationale:** `golang.org/x/image v0.39.0` is already an indirect dependency in `go.mod`. Go Mono is a professional-quality monospaced font by Bigelow & Holmes with WGL4 charset (Latin + Cyrillic + Greek, 650+ glyphs). Rasterization uses `font.Drawer` (anti-aliased). Zero maintenance — no hand-crafted glyph data.
**Alternatives considered:**
- Hand-crafted 7×10 bitmap — zero deps but ugly, no anti-alias, manual maintenance
- `golang.org/x/image/font/basicfont` — no Cyrillic support
- TrueType font loading from disk — fragile in Docker, adds file dependency

### 2. Caption parameter as `string` (not pointer)
**Choice:** `GenerateBoard(fen string, caption ...string)` — variadic for backward compatibility.
**Rationale:** All existing callers pass `fen` only. A variadic string parameter means zero changes at call sites. Internally, default to empty string if no caption given.
**Alternatives considered:**
- `*string` pointer — more explicit but awkward to call with `nil`
- `Config` struct — breaks existing API more aggressively; unnecessary for one optional field

### 3. Word-wrap algorithm
**Choice:** Split caption text on whitespace, measure each word with `font.Drawer.MeasureString`, build lines greedily.
**Rationale:** TrueType glyphs vary slightly even in monospace (kerning, wide chars). `MeasureString` is the correct API; it handles all font metrics implicitly. No hyphenation needed since Russian text doesn't break mid-word.
**Edge case:** Single word longer than board width → still render on one line (overflow right).

### 4. Image sizing
**Choice:** Increase canvas height dynamically based on line count. Line height comes from `face.Metrics().Height` (≈13-14px at 10pt) + 2px spacing. Padding 6px from board bottom. New total height = original 570px + 6 + (lineHeight × lineCount).
**Rationale:** Produces compact images. Fixed padding keeps positioning predictable.

## Risks / Trade-offs

- **Binary size** → `gomono.TTF` is ≈200KB of embedded TTF data. For a CLI tool that's negligible. Mitigation: only one font face loaded.
- **Font init cost** → First call with caption parses TTF and rasterizes a face. A few milliseconds, one-shot via `sync.Once`. Mitigation: lazy init, no cost if caption unused.
- **Image size increase** → A multi-line caption can double image height (e.g., 4 lines = ~60px → 630px total). Acceptable for educational use. Mitigation: caption is optional, so images stay compact when unused.
