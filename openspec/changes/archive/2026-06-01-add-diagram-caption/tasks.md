## 1. Go Mono font integration

- [x] 1.1 Add `golang.org/x/image/font/gofont/gomono` and `font/opentype` imports
- [x] 1.2 Implement lazy font face init with `sync.Once` via `opentype.Parse(gomono.TTF)`
- [x] 1.3 Implement `getCaptionFace() (font.Face, error)` — parse + create face at 10pt/72DPI
- [x] 1.4 Create `captionRenderer` struct wrapping face, exposing `lineHeight` and `charWidth`

## 2. Word-wrap and centering with MeasureString

- [x] 2.1 Implement `wrapText(face font.Face, caption string, maxPixels fixed.Int26_6) []string` — use `font.Drawer.MeasureString`, greedy line building
- [x] 2.2 Implement `drawCaption(img *image.RGBA, caption string, imgWidth int)` — compute centered baseline X per line, use `font.Drawer.DrawString`

## 3. Modify GenerateBoard

- [x] 3.1 Remove bitmap font constants from fenlib.go (fontWidth, fontHeight, charSpacing etc.)
- [x] 3.2 Update GenerateBoard to use face-based wrapText + drawCaption
- [x] 3.3 Line height from face.Metrics().Height, padding 6px
- [x] 3.4 Enlarge canvas and render caption below board
- [x] 3.5 Keep variadic signature, backward compat

## 4. Update forwarding functions

- [x] 4.1 Update `GenerateDiagram` with optional caption
- [x] 4.2 Update `GenerateDiagramBase64` with optional caption
- [x] 4.3 Update `GenerateDiagramToWriter` with optional caption

## 5. Update fen_to_diagram CLI

- [x] 5.1 Add `--caption` string flag parsing in fen_to_diagram.go
- [x] 5.2 Forward caption value to GenerateDiagram/GenerateDiagramBase64 calls
- [x] 5.3 Update usage/help text to document --caption flag

## 6. Update markdown_builder callers

- [x] 6.1 Find call sites in build_markdown.go that call GenerateDiagramBase64
- [x] 6.2 Verify backward compatibility — no changes needed, variadic API works as-is

## 7. Tests

- [x] 7.1 Test single-line caption renders and image height increases
- [x] 7.2 Test multi-line word-wrapped caption
- [x] 7.3 Test empty/missing caption produces identical output
- [x] 7.4 Test Russian Cyrillic caption text renders without errors
- [x] 7.5 Test single-word-longer-than-width edge case
- [x] 7.6 Test unsupported glyph replacement with `?`

## 8. Documentation

- [x] 8.1 Update README.md with `--caption` flag examples (if needed)
- [x] 8.2 Update AGENTS.md fen_builder section (if needed)
- [x] 8.3 Run `go mod tidy` and final validation
