## MODIFIED Requirements

### Requirement: Exclusion composes with --exam

**Original:** When both `--white-except`/`--black-except` and `--exam` are set, exclusion SHALL be applied before inclusion.

**Updated:** When `--white-except`/`--black-except`, `--white-only`/`--black-only`, and `--exam` are set, the pipeline order SHALL be: exclusion first, then inclusion, then exam.

#### Scenario: Excluded chapter never reaches inclusion filter

- **WHEN** `--white-except "Carlsen" --white-only "Carlsen" --exam "Carlsen"` is provided
- **THEN** no chapters SHALL appear in output (exclusion removes all Carlsen games before inclusion or exam can process them)

#### Scenario: Inclusion filters before exam

- **WHEN** `--white-only "Nakamura" --exam "Carlsen"` is provided and a chapter has White=`"Nakamura"`
- **THEN** the chapter SHALL be omitted from the exam filter (already removed by inclusion, so exam never sees it)

#### Scenario: Exclusion then inclusion then exam

- **WHEN** `--white-except "Carlsen" --white-only "Carlsen|Nakamura" --exam "Nakamura"` is provided and chapters have White=`["Carlsen", "Nakamura"]`
- **THEN** exclusion removes Carlsen chapter, inclusion keeps Nakamura chapter, exam matches Nakamura → Nakamura chapter SHALL appear in output
