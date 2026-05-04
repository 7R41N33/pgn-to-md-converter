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

Скрипт `fen_builder` генерирует изображения шахматных позиций из FEN-строк, используя спрайт `src/images/Chess_Pieces_Sprite.svg.png`.

### Через бинарник

```bash
./bin/fen_builder <fen_string> [output_path] [--base64]
```

### Через исходный код

```bash
cd src/builders/fen_builder/
go run fen_to_diagram.go <fen_string> [output_path] [--base64]
```

### Параметры

- `<fen_string>` — FEN-строка (обязательно)
- `[output_path]` — Путь для сохранения PNG-изображения (опционально, по умолчанию `tmp/images/diagram.png`)
- `--base64` — Вывод base64-кодированного PNG в stdout вместо сохранения файла (опционально)

### Примеры

```bash
# Генерация диаграммы с указанием пути
./bin/fen_builder "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1" board.png

# Генерация с путём по умолчанию (tmp/images/diagram.png)
./bin/fen_builder "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

# Пример из пользовательского запроса
./bin/fen_builder "8/pp3ppp/8/2p5/8/3P2P1/PP2P2P/8 w - - 0 1"

# Вывод в base64 (для вставки в Markdown)
./bin/fen_builder "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1" --base64

# Использование base64 вывода в Markdown
./bin/fen_builder "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1" --base64 > /tmp/board.b64
echo "![Позиция](data:image/png;base64,$(cat /tmp/board.b64))" >> output.md

# Через исходный код
cd src/builders/fen_builder/
go run fen_to_diagram.go "8/8/8/8/8/8/8 w - - 0 1" empty_board.png
```

### Описание

- **Вход:** FEN-строка (стандартный формат записи позиции)
- **Выход (файл):** PNG-изображение шахматной доски (64x64 пикселей на клетку, всего 540x540 с учётом белой границы)
- **Выход (base64):** Base64-кодированная строка PNG-изображения для вставки в HTML/Markdown
- **Требования:** Валидный FEN (при некорректном формате возвращает ошибку, не падает с panic)
- **Спрайт:** Используется `src/images/Chess_Pieces_Sprite.svg.png` (1280x427, 2 строки x 6 столбцов)
- **Границы:** Белая граница вокруг доски с координатами (a-h, 1-8)

### Тестирование

```bash
cd src/builders/fen_builder/
go test -v
```

### Использование base64 вывода

Для вставки диаграмм напрямую в Markdown-документы без сохранения файлов:

```bash
# Получение base64 строки
BASE64_IMAGE=$(./bin/fen_builder "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1" --base64)

# Вставка в Markdown
echo "## Позиция после 1.e4" > output.md
echo "![Диаграмма](data:image/png;base64,${BASE64_IMAGE})" >> output.md
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

# fen_builder
cd src/builders/fen_builder/
go build -o ../../../bin/fen_builder fen_to_diagram.go
```

После сборки бинарники автоматически попадут в папку `bin/`.
