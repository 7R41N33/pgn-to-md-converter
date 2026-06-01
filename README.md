# Chess PGN to Markdown Converter

A project for converting chess PGN files into formatted Markdown.

## Project Structure

- `bin/` — Compiled binary files (use these for running)
- `src/builders/fenlib/` — Library for generating chess diagrams (used by fen_builder and markdown_builder)
- `src/builders/fen_builder/` — Source code for the standalone diagram generator
- `src/builders/markdown_builder/` — Source code for the PGN to Markdown converter
- `PGNs/` — Input PGN files
- `Markdowns/` — Output Markdown files

## Docker

The project includes a `Dockerfile` and `docker-compose.yml` for containerized development.

### Prerequisites

- [Docker](https://docs.docker.com/engine/install/)
- [Docker Compose](https://docs.docker.com/compose/install/) (included with Docker Desktop)

### Build binaries inside container

```bash
make docker-build
# or directly:
docker compose run build
```

This compiles both binaries (`build_markdown` and `fen_builder`) and places them in `bin/` on the host.

### Interactive development shell

```bash
docker compose run dev bash
```

Opens a shell with the Go toolchain and the entire project mounted. Run any Go commands:

```bash
go build -o bin/build_markdown src/builders/markdown_builder/build_markdown.go
go test ./src/...
```

### Build the image only

```bash
docker compose build
```

### Project structure inside container

| Host path | Container path |
|-----------|---------------|
| `./PGNs/` | `/app/PGNs/` |
| `./Markdowns/` | `/app/Markdowns/` |
| `./tmp/` | `/app/tmp/` |
| `./bin/` | `/app/bin/` |

All directories are bind-mounted — changes on the host are immediately visible in the container and vice versa.

## Using markdown_builder

The `markdown_builder` script converts PGN files into formatted Markdown. Supports embedding diagrams as base64 images.

### Flags

- `-src` — Path to input PGN file (required)
- `-out` — Path to output Markdown file (required)
- `-skip` — Skip the first N chapters (optional)
- `-chapters` — Process only the first N chapters (after skip) (optional)
- `-chapter-numbers` — Process only the specified chapter numbers (1-based, comma-separated) (optional)
- `--inline-images` — Replace FEN strings with data URI format: `data:image/png;base64,...` (optional)
- `-numbered-lists` — Format numbered lists (default: true) (optional)
- `-dash-lists` — Format dash lists (default: false) (optional)
- `-exam` — Exam mode: processes only chapters where White or Black **exactly** match the given string. If the chapter contains a FEN — outputs only the FEN. If FEN is absent — applies normal formatting (optional)

### Examples

```bash
# Basic conversion
./bin/build_markdown -src input.pgn -out output.md

# Conversion with embedded base64 diagrams
./bin/build_markdown -src input.pgn -out output.md --inline-images

# Skip first 2 chapters, process next 5
./bin/build_markdown -src input.pgn -out output.md -skip 2 -chapters 5

# Process only chapters 1, 3 and 5
./bin/build_markdown -src input.pgn -out output.md -chapter-numbers 1,3,5

# Exam mode — only chapters with White="Exam Time!"
./bin/build_markdown -src input.pgn -out output.md -exam "Exam Time!"

# Exam mode with inline images
./bin/build_markdown -src input.pgn -out output.md -exam "Exam Time!" --inline-images

# Disable numbered list formatting
./bin/build_markdown -src input.pgn -out output.md -numbered-lists=false

# Enable dash list formatting
./bin/build_markdown -src input.pgn -out output.md -dash-lists=true
```

### Description

- **Input:** PGN file with chess games
- **Output:** Markdown file with formatted moves
- **FEN strings:** By default preserved as `**FEN:** \`FEN string\``. With `--inline-images` flag replaced with `data:image/png;base64,...` (without the `**FEN:**` marker)
- **Chapter headings:** Automatically formats chapter headings: `White` → `##`, `Black` → `###`. If `White` repeats, skips `##` and uses only `###`
- **Numbered lists:** By default formatted on new lines with indentation. Can be disabled via `-numbered-lists=false`
- **Dash lists:** By default not formatted. Can be enabled via `-dash-lists=true`
- **Exam mode:** Useful for creating study materials — exercise chapters show only the position (FEN), while introductory chapters are fully formatted
- **Dependencies:** Uses the `fenlib` package for diagram generation

## Using fen_builder

The `fen_builder` script generates chess position images from FEN strings. Supports saving to a PNG file or outputting in base64 format.

### Via binary (recommended)

```bash
./bin/fen_builder <fen_string> [output_path] [--base64] [--caption "text"]
```

### Via source code

```bash
cd src/builders/fen_builder/
go run fen_to_diagram.go <fen_string> [output_path] [--base64] [--caption "text"]
```

### Parameters

- `<fen_string>` — FEN string (required)
- `[output_path]` — Path to save the PNG file (default: `tmp/images/diagram.png`)
- `--base64` — Output the image as base64 to stdout instead of saving to file (optional)
- `--caption "text"` — Add centered, word-wrapped caption text below the diagram (optional)

### Examples

```bash
# Generate a diagram and save to file
./bin/fen_builder "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBKQBNR w KQkq - 0 1" output.png

# Generate and output as base64 (for embedding in HTML/Markdown)
./bin/fen_builder "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBKQBNR w KQkq - 0 1" --base64

# Generate a diagram with a Russian caption (word-wraps if too long)
./bin/fen_builder "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1" output.png --caption "Черные сыграли h5. Какова их идея?"
```

### Description

- **Input:** FEN string (6 fields: position, active color, castling, en passant, halfmove, fullmove)
- **Output:** PNG image with enlarged border (8 squares at 60px + 45px borders = ~570x570 px)
- **Move indicator:** A triangle pointing upward in the bottom-right corner indicates whose turn it is. White to move — white triangle with black border, black to move — black triangle.
- **FEN format:** Starts with rank 8 (closest to black), ends with rank 1 (closest to white)
- **Square colors:** Light (#F0D9B5) and dark (#B58863)
- **Pieces:** Uppercase letters — white, lowercase — black
- **Dependencies:** Uses the sprite `src/images/Chess_Pieces_Sprite.svg.png`

### Testing

```bash
# fenlib tests
cd src/builders/fenlib/
go test -v

# markdown_builder tests
cd src/builders/markdown_builder/
go test -v
```

## Using base64_to_image.gs (Google Docs)

The `src/base64_to_image.gs` script is designed for use in Google Docs. It automatically converts base64-encoded images (generated with the `--inline-images` flag) into real images within the document.

### Setup and usage

1. Open the Google Docs document where the Markdown with base64 images was pasted
2. Go to **Extensions** → **Apps Script**
3. Copy the contents of `src/base64_to_image.gs` into the Apps Script editor
4. Click the floppy disk icon (Save) or press Ctrl+S
5. Return to the document and refresh the page
6. To run the script, go to **Extensions** → **Apps Script** → select function `convertAllBase64ToImages` → click **Run**
7. On first run, Google will request permissions (needs access to the current document) — grant them

### What the script does

- Finds all base64 strings in `data:image/...;base64,...` format within the document
- Decodes them into binary data
- Replaces the text base64 strings with real images (InlineImage)
- Processes images in reverse order (end to start) to avoid index shifting

### Note

The script only works with Google Docs documents containing base64 images generated using the `--inline-images` flag.
