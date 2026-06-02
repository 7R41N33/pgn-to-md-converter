## ADDED Requirements

### Requirement: Include chapters by White player name

The system SHALL accept a `--white-only` flag and include only chapters whose `White` PGN field matches any of the given values. The inclusion filter SHALL be applied after `--white-except`/`--black-except` but before `--exam`.

#### Scenario: Single value includes exact match

- **WHEN** `--white-only "Carlsen"` is provided and chapters have White values `["Carlsen", "Nakamura", "Carlsen"]`
- **THEN** only chapters where `White` is `"Carlsen"` SHALL appear in output

#### Scenario: Non-matching White is excluded

- **WHEN** `--white-only "Carlsen"` is provided
- **THEN** chapters where `White` is not `"Carlsen"` SHALL be omitted from output

#### Scenario: Multiple pipe-separated values

- **WHEN** `--white-only "Carlsen|Nakamura"` is provided
- **THEN** chapters where `White` is `"Carlsen"` OR `"Nakamura"` SHALL appear in output

### Requirement: Include chapters by Black player name

The system SHALL accept a `--black-only` flag and include only chapters whose `Black` PGN field matches any of the given values.

#### Scenario: Single value includes exact match

- **WHEN** `--black-only "Nepomniachtchi"` is provided
- **THEN** only chapters where `Black` is `"Nepomniachtchi"` SHALL appear in output

#### Scenario: Multiple pipe-separated values

- **WHEN** `--black-only "Nepomniachtchi|Karjakin"` is provided
- **THEN** chapters where `Black` is `"Nepomniachtchi"` OR `"Karjakin"` SHALL appear in output

### Requirement: Both flags together act as logical AND

When both `--white-only` and `--black-only` are set, a chapter SHALL be included only if both conditions are satisfied.

#### Scenario: Both flags match

- **WHEN** `--white-only "Carlsen" --black-only "Nakamura"` is provided and a chapter has White=`"Carlsen"` and Black=`"Nakamura"`
- **THEN** the chapter SHALL appear in output

#### Scenario: Only one flag matches

- **WHEN** `--white-only "Carlsen" --black-only "Nakamura"` is provided and a chapter has White=`"Other"` and Black=`"Nakamura"`
- **THEN** the chapter SHALL be omitted from output
