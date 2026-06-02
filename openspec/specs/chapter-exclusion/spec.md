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
When `--white-except`/`--black-except`, `--white-only`/`--black-only`, and `--exam` are set, the pipeline order SHALL be: exclusion first, then inclusion, then exam.

#### Scenario: Excluded chapter never reaches inclusion filter
- **WHEN** `--white-except "Carlsen" --white-only "Carlsen" --exam "Carlsen"` is provided
- **THEN** no chapters SHALL appear in output (exclusion removes all Carlsen games before inclusion or exam can process them)

#### Scenario: Inclusion filters before exam
- **WHEN** `--white-only "Nakamura" --exam "Carlsen"` is provided and a chapter has White=`"Nakamura"`
- **THEN** the chapter SHALL be omitted from the exam filter (already removed by inclusion, so exam never sees it)

#### Scenario: Exclusion then inclusion then exam
- **WHEN** `--white-except "Carlsen" --white-only "Carlsen|Nakamura" --exam "Nakamura"` is provided and chapters have White=`["Carlsen", "Nakamura"]`
- **THEN** exclusion removes Carlsen chapter, inclusion keeps Nakamura chapter, exam matches Nakamura → Nakamura chapter SHALL appear in output
