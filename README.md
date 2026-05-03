# Chess PGN to Markdown Converter

Проект для конвертации шахматных PGN-файлов в форматированный Markdown.

## Структура проекта

- `bin/` — Скомпилированные бинарные файлы (используйте их для запуска)
- `src/builders/markdown_builder/` — Исходный код скрипта
- `PGNs/` — Входные PGN-файлы
- `Markdowns/` — Выходные Markdown-файлы

## Использование markdown_builder

### Через бинарник (рекомендуется)

```bash
./bin/build_markdown -src <path_to_pgn> -out <path_to_md> [-filter-white <name>] [-skip <n>]
```

### Через исходный код

```bash
cd src/builders/markdown_builder/
go run build_markdown.go -src <path_to_pgn> -out <path_to_md> [-filter-white <name>] [-skip <n>]
```

### Параметры

- `-src` — Путь к входному PGN-файлу (обязательно)
- `-out` — Путь к выходному Markdown-файлу (обязательно)
- `-filter-white` — Фильтр по имени белого игрока (опционально)
- `-skip` — Пропустить N первых глав (опционально)

### Примеры

```bash
# Конвертация одного файла
./bin/build_markdown -src "PGNs/Shankland, Sam - Neo Catalan Part-1 (2023)/games.pgn" \
                   -out "Markdowns/Shankland, Sam - Neo Catalan Part-1 (2023)/games.md"

# Конвертация с пропуском первых 5 глав
./bin/build_markdown -src input.pgn -out output.md -skip 5

# Фильтрация по имени игрока
./bin/build_markdown -src input.pgn -out output.md -filter-white "Carlsen"
```

### Тестирование

```bash
cd src/builders/markdown_builder/
go test -v
```

## Использование fen_builder

Скрипт `fen_builder` генерирует изображения шахматных позиций из FEN-строк.

### Через бинарник

```bash
./bin/fen_to_diagram <fen_string> <output_path>
```

### Через исходный код

```bash
cd src/builders/fen_builder/
go run fen_to_diagram.go <fen_string> <output_path>
```

### Примеры

```bash
# Генерация диаграммы из FEN
./bin/fen_to_diagram "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1" board.png

# Через исходный код
cd src/builders/fen_builder/
go run fen_to_diagram.go "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1" board.png
```

### Описание

- **Вход:** FEN-строка (стандартный формат записи позиции)
- **Выход:** PNG-изображение шахматной доски
- **Требования:** Валидный FEN (при некорректном формате возвращает ошибку, не падает с panic)

### Тестирование

```bash
cd src/builders/fen_builder/
go test -v
```

## Правила конвертации PGN в Markdown

Скрипт `markdown_builder` следует правилам из `AGENTS.md`:

1. **Парная запись ходов:** `N. ХодБелых ХодЧёрных` (одна строка)
2. **Удаление `{` и `}`** из комментариев
3. **Отступ 2 пробела** для комментариев после пустой строки
4. **Конвертация результатов:** `1-0` → `1-0 (Победа белых)` и т.д.
5. **Конвертация NAG:** `$1` → `±` и т.д.
6. **Заголовки глав:** `## Chapter N: Игрок1 vs. Игрок2, Город Год`

Подробности в `AGENTS.md`.

## Сборка бинарников

Если нужно пересобрать бинарники:

```bash
# markdown_builder
cd src/builders/markdown_builder/
go build -o ../../../bin/build_markdown build_markdown.go

# fen_builder (если есть)
cd src/builders/fen_builder/
go build -o ../../../bin/fen_to_diagram fen_to_diagram.go
```

После сборки бинарники автоматически попадут в папку `bin/`.
