## Why

Пользователям нужно выбирать главы только конкретных игроков без переключения в режим `--exam`. Флаги `--white-only` и `--black-only` работают как инверсная пара к `--white-except`/`--black-except` — включают только те главы, где имя игрока совпадает.

## What Changes

- Добавить флаг `--white-only` — оставить в выводе только главы, где `White` совпадает с указанным значением (pipe-разделитель для нескольких имён)
- Добавить флаг `--black-only` — оставить только главы, где `Black` совпадает с указанным значением
- Порядок фильтрации: `--white-except`/`--black-except` → `--white-only`/`--black-only` → `--exam`
- `--white-only` и `--black-only` могут использоваться одновременно (логическое И)
- Флаги используют `splitNames` helper (pipe-разделитель), уже реализованный для `--white-except`

## Capabilities

### New Capabilities
- `chapter-inclusion`: Включение глав по имени игрока (противоположность chapter-exclusion)

### Modified Capabilities
- `chapter-exclusion`: расширяется правилом композиции с inclusion-флагами (exclusion → inclusion → exam)

## Impact

- `src/builders/markdown_builder/build_markdown.go`: два новых флага, логика фильтрации
- `src/builders/markdown_builder/build_markdown_test.go`: тесты
- `AGENTS.md`, `README.md`: документация
