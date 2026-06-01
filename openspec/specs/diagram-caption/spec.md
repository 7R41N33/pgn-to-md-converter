## ADDED Requirements

### Requirement: Render caption below chess diagram
The system SHALL accept an optional caption string and render it as center-aligned, word-wrapped text directly below the chess board diagram.

#### Scenario: Caption renders below the board
- **WHEN** a non-empty caption string is provided to GenerateBoard
- **THEN** the output image SHALL be taller than the uncaptioned version by an amount proportional to the number of rendered lines
- **THEN** the caption text SHALL be visible in the added vertical area below the board

#### Scenario: Empty caption produces identical output
- **WHEN** no caption or an empty caption string is provided
- **THEN** the output image SHALL be identical in dimensions and content to the current uncaptioned output

### Requirement: Text is center-aligned
The caption text SHALL be horizontally centered within the board image width.

#### Scenario: Single-line caption is centered
- **WHEN** a single-line caption is rendered
- **THEN** the text SHALL be positioned so that its horizontal center matches the image center

#### Scenario: Multi-line caption lines are centered
- **WHEN** a multi-line caption wraps across multiple lines
- **THEN** each line SHALL be independently centered horizontally

### Requirement: Text word-wraps when exceeding diagram width
If a line of caption text would exceed the board image width, the system SHALL split it at word boundaries onto multiple lines.

#### Scenario: Long caption wraps to multiple lines
- **WHEN** the caption text exceeds the available width
- **THEN** the text SHALL be broken at word boundaries (spaces) onto the next line

#### Scenario: Single word longer than width
- **WHEN** a single word is longer than the diagram width
- **THEN** the word SHALL be rendered on its own line, potentially overflowing the right edge

### Requirement: Support Latin and Cyrillic characters
The system SHALL support Latin letters (A–Z, a–z), digits (0–9), common punctuation, and Russian Cyrillic characters (А–Я, а–я, Ё, ё) in caption text.

#### Scenario: Russian caption renders correctly
- **WHEN** the caption contains Russian Cyrillic text
- **THEN** the text SHALL be rendered legibly using the Go Mono TrueType font

#### Scenario: Unsupported glyph is replaced
- **WHEN** the caption contains a character not defined in the Go Mono WGL4 charset
- **THEN** the system SHALL render the character as the font's `.notdef` glyph (usually a box or `?`)
