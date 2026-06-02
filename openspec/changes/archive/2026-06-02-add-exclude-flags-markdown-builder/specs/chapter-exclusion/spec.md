## ADDED Requirements

### Requirement: Exclude chapters by White player name
The system SHALL accept a `--white-except` flag and skip any chapter whose `White` PGN field matches any of the given values.

#### Scenario: Single value excludes exact match
- **WHEN** `--white-except "Carlsen"` is provided
- **THEN** chapters where `White` is `"Carlsen"` SHALL be omitted from output

#### Scenario: Non-matching White is included
- **WHEN** `--white-except "Carlsen"` is provided
- **THEN** chapters where `White` is not `"Carlsen"` SHALL be processed normally

#### Scenario: Multiple pipe-separated values
- **WHEN** `--white-except "Carlsen|Nakamura"` is provided
- **THEN** chapters where `White` is `"Carlsen"` OR `"Nakamura"` SHALL be omitted from output

#### Scenario: Whitespace is trimmed before comparison
- **WHEN** the `White` field in PGN has leading or trailing whitespace
- **THEN** the whitespace SHALL be trimmed before comparing with `--white-except`

### Requirement: Exclude chapters by Black player name
The system SHALL accept a `--black-except` flag and skip any chapter whose `Black` PGN field matches any of the given values.

#### Scenario: Single value excludes exact match
- **WHEN** `--black-except "Nepomniachtchi"` is provided
- **THEN** chapters where `Black` is `"Nepomniachtchi"` SHALL be omitted from output

#### Scenario: Multiple pipe-separated values
- **WHEN** `--black-except "Nepomniachtchi|Karjakin"` is provided
- **THEN** chapters where `Black` is `"Nepomniachtchi"` OR `"Karjakin"` SHALL be omitted from output

### Requirement: Exclusion composes with --exam
When both `--white-except`/`--black-except` and `--exam` are set, exclusion SHALL be applied before inclusion.

#### Scenario: Excluded chapter never reaches exam filter
- **WHEN** `--white-except "Carlsen" --exam "Carlsen"` is provided
- **THEN** no chapters SHALL appear in output (exclusion removes all Carlsen games before exam tries to include them)
