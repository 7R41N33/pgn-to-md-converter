## ADDED Requirements

### Requirement: Include chapters by White player name using regex

The system SHALL accept a `--white-only-regex` flag and include only chapters whose `White` PGN field matches any of the given regular expression patterns. Each pattern SHALL be compiled via `regexp.Compile` and matched via `MatchString`. The inclusion filter SHALL be applied after `--white-only` but before `--black-only` in the filter chain.

#### Scenario: Regex matches substring

- **WHEN** `--white-only-regex "Carl"` is provided and chapters have White values `["Carlsen, Magnus", "Nakamura, Hikaru"]`
- **THEN** only the chapter with White `"Carlsen, Magnus"` SHALL appear in output

#### Scenario: Regex with character class

- **WHEN** `--white-only-regex "[CK]arl"` is provided and chapters have White values `["Carlsen", "Karpov", "Kasparov"]`
- **THEN** chapters with White `"Carlsen"` and `"Karpov"` SHALL appear; `"Kasparov"` SHALL NOT

#### Scenario: No regex match excludes chapter

- **WHEN** `--white-only-regex "Nakamura"` is provided and a chapter has White `"Carlsen, Magnus"`
- **THEN** the chapter SHALL be omitted from output

#### Scenario: Multiple pipe-separated patterns

- **WHEN** `--white-only-regex "Carl|Nakam"` is provided and chapters have White values `["Carlsen", "Nakamura", "Anand"]`
- **THEN** chapters `"Carlsen"` and `"Nakamura"` SHALL appear; `"Anand"` SHALL NOT

#### Scenario: Invalid regex pattern causes error

- **WHEN** `--white-only-regex "[invalid"` is provided
- **THEN** the system SHALL exit with a non-zero status and print an error message describing the invalid pattern

#### Scenario: Case-sensitive by default

- **WHEN** `--white-only-regex "carlsen"` is provided and a chapter has White `"Carlsen, Magnus"`
- **THEN** the chapter SHALL be omitted from output (case mismatch)

#### Scenario: Case-insensitive flag works

- **WHEN** `--white-only-regex "(?i)carlsen"` is provided and a chapter has White `"Carlsen, Magnus"`
- **THEN** the chapter SHALL appear in output

#### Scenario: Combined with --white-only (AND logic)

- **WHEN** `--white-only "Carlsen, Magnus" --white-only-regex "Magnus"` is provided and a chapter has White `"Carlsen, Magnus"`
- **THEN** the chapter SHALL appear in output (passes both filters)

#### Scenario: Combined with --white-only (no regex match)

- **WHEN** `--white-only "Carlsen, Magnus" --white-only-regex "Nakamura"` is provided and a chapter has White `"Carlsen, Magnus"`
- **THEN** the chapter SHALL be omitted from output (fails regex filter)

### Requirement: --white-only-regex applied after --white-except and --black-except

When `--white-except`/`--black-except` and `--white-only-regex` are set together, the exclusion filters SHALL run first. A chapter removed by exclusion SHALL never reach the regex inclusion filter.

#### Scenario: Exclusion removes before regex inclusion

- **WHEN** `--white-except "Carlsen" --white-only-regex "Carl"` is provided and chapters have White values `["Carlsen", "Carlsen, Magnus"]`
- **THEN** both chapters SHALL be omitted (exclusion removes all Carlsen matches before regex can include them)
