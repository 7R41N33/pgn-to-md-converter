## ADDED Requirements

### Requirement: Split output by chapter count

The system SHALL accept a `--split-by-chapters <int>` flag. When provided, the output SHALL be split into multiple files, each containing at most the specified number of chapters. The split SHALL be applied after all chapter filters (`--skip`, `--chapters`, `--chapter-numbers`, `--exam`, `--white-except`, `--black-except`).

#### Scenario: Split into multiple files

- **WHEN** `--out theory.md --split-by-chapters 10` is used and 25 chapters remain after filters
- **THEN** three files SHALL be created: `theory_1.md` (chapters 1–10), `theory_2.md` (chapters 11–20), `theory_3.md` (chapters 21–25)

#### Scenario: No split when chapters ≤ threshold

- **WHEN** `--out theory.md --split-by-chapters 100` is used and 25 chapters remain after filters
- **THEN** a single file `theory.md` SHALL be created (no suffix)

#### Scenario: Split with various file extensions

- **WHEN** `--out /path/to/file.md --split-by-chapters 5` is used
- **THEN** files SHALL be named `/path/to/file_1.md`, `/path/to/file_2.md`, etc.

### Requirement: Split composes with chapter filters

When `--split-by-chapters` is combined with any chapter filter (`--skip`, `--chapters`, `--chapter-numbers`, `--exam`, `--white-except`, `--black-except`), the split SHALL be applied after all filters have been processed. Only the chapters that pass all filters SHALL be counted for the split.

#### Scenario: Split after skip

- **WHEN** `--skip 5 --split-by-chapters 10` is used and 30 chapters total, 25 after skip
- **THEN** three files SHALL be created with chapters remaining after the skip

#### Scenario: Split after exam + white-except

- **WHEN** `--exam "Target" --white-except "Carlsen" --split-by-chapters 5` is used
- **THEN** only chapters matching the exam filter (after excluding Carlsen) SHALL be split

### Requirement: Zero or negative split value ignored

The system SHALL ignore `--split-by-chapters` when the value is ≤ 0 and produce a single unsplit file.

#### Scenario: Zero value treated as no split

- **WHEN** `--split-by-chapters 0` is used
- **THEN** output SHALL be a single unsplit file

#### Scenario: Negative value treated as no split

- **WHEN** `--split-by-chapters -1` is used
- **THEN** output SHALL be a single unsplit file
