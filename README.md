# Chess PGN to Markdown Converter

Проект для конвертации шахматных PGN-файлов в форматированный Markdown.

## Структура проекта

- `bin/` — Скомпилированные бинарные файлы (используйте их для запуска)
- `src/builders/fenlib/` — Библиотека для генерации шахматных диаграмм (используется fen_builder и markdown_builder)
- `src/builders/fen_builder/` — Исходный код скрипта для генерации отдельных диаграмм
- `src/builders/markdown_builder/` — Исходный код скрипта для конвертации PGN в Markdown
- `PGNs/` — Входные PGN-файлы
- `Markdowns/` — Выходные Markdown-файлы

## Использование markdown_builder

Скрипт `markdown_builder` конвертирует PGN-файлы в форматированный Markdown. Поддерживает вставку диаграмм в виде base64-изображений.

### Через бинарник (рекомендуется)

```bash
./bin/build_markdown -src <path_to_pgn> -out <path_to_md> [-filter-white <name>] [-skip <n>] [-chapters <n>] [-chapter-numbers <list>] [--inline-images]
```

### Через исходный код

```bash
cd src/builders/markdown_builder/
go run build_markdown.go -src <path_to_pgn> -out <path_to_md> [-filter-white <name>] [-skip <n>] [-chapters <n>] [-chapter-numbers <list>] [--inline-images]
```

### Параметры

- `-src` — Путь к входному PGN-файлу (обязательно)
- `-out` — Путь к выходному Markdown-файлу (обязательно)
- `-filter-white` — Фильтр по имени белого игрока (опционально)
- `-skip` — Пропустить N первых глав (опционально)
- `-chapters` — Обработать только первые N глав (после skip) (опционально)
- `-chapter-numbers` — Обработать только указанные номера глав (1‑based, через запятую) (опционально)
- `--inline-images` — Заменять FEN-строки на data URI формат: `data:image/png;base64,...` (опционально)
- `-numbered-lists` — Форматировать нумерованные списки (по умолчанию: true) (опционально)
- `-dash-lists` — Форматировать списки с тире (по умолчанию: false) (опционально)

### Примеры

```bash
# Конвертация с встроенными base64 диаграммами
./bin/build_markdown -src input.pgn -out output.md --inline-images

# Отключить форматирование нумерованных списков
./bin/build_markdown -src input.pgn -out output.md -numbered-lists=false

# Включить форматирование списков с тире
./bin/build_markdown -src input.pgn -out output.md -dash-lists=true

# Комбинация параметров
./bin/build_markdown -src input.pgn -out output.md -skip 2 -chapters 8 -numbered-lists=true -dash-lists=true
```

### Описание

- **Вход:** PGN-файл с шахматными партиями
- **Выход:** Markdown-файл с отформатированными ходами
- **FEN-строки:** По умолчанию сохраняются как текст. С флагом `--inline-images` заменяются на `data:image/png;base64,...` (без markdown-обертки)
- **Заголовки глав:** Автоматически форматирует заголовки глав: `White` → `##`, `Black` → `###`. Если `White` повторяется, пропускает `##` и использует только `###`
- **Нумерованные списки:** По умолчанию форматируются с новой строки и отступом. Можно отключить через `-numbered-lists=false`
- **Списки с тире:** По умолчанию не форматируются. Можно включить через `-dash-lists=true`
- **Зависимости:** Использует пакет `fenlib` для генерации диаграмм

### Тестирование

```bash
cd src/builders/markdown_builder/
go test -v
```
