## Why

При конвертации больших PGN-файлов (сотни глав) в единый Markdown файл становится неудобно навигировать по содержимому. Нужна возможность автоматически разбивать выходной файл на несколько частей, каждая из которых содержит заданное количество глав.

## What Changes

- Добавить новый флаг `--split-by-chapters <int>` в `markdown_builder`
- При указании этого флага выходной Markdown разбивается на N файлов, каждый содержит не более указанного числа глав
- Имя файла получает числовой суффикс: `output_1.md`, `output_2.md`, `output_3.md` и т.д.
- Если число глав меньше или равно `--split-by-chapters`, разбивка не выполняется (выходной файл остаётся единым)
- Флаг работает совместно с существующими `--skip`, `--chapters`, `--chapter-numbers`, `--exam`, `--white-except`/`--black-except` — разбивка применяется после всех фильтров

## Capabilities

### New Capabilities
- `output-split`: Разбивка выходного Markdown-файла на несколько частей по заданному числу глав

### Modified Capabilities
<!-- None — existing specs are unchanged -->

## Impact

- `src/builders/markdown_builder/build_markdown.go`: новая логика записи результата в несколько файлов
- `src/builders/markdown_builder/build_markdown_test.go`: тесты для флага и сценариев разбивки
