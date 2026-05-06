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

## Использование fen_builder

Скрипт `fen_builder` генерирует изображения шахматных позиций из FEN-строк. Поддерживает сохранение в PNG-файл или вывод в формате base64.

### Через бинарник (рекомендуется)

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
- `[output_path]` — Путь для сохранения PNG-файла (по умолчанию: `tmp/images/diagram.png`)
- `--base64` — Вывод изображения в формате base64 в stdout вместо сохранения в файл (опционально)

### Примеры

```bash
# Генерация диаграммы и сохранение в файл
./bin/fen_builder "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBKQBNR w KQkq - 0 1" output.png

# Генерация и вывод в base64 (для встраивания в HTML/Markdown)
./bin/fen_builder "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBKQBNR w KQkq - 0 1" --base64
```

### Описание

- **Вход:** FEN-строка (6 полей: позиция, активный цвет, рокировка, en passant, полуход, полный ход)
- **Выход:** PNG-изображение 540x540 пикселей (8 клеток по 60px + границы по 30px)
- **Формат FEN:** Начинается с 8-го ранга (ближайшего к чёрным), заканчивается 1-м рангом (ближайшим к белым)
- **Цвета клеток:** Светлые (#F0D9B5) и тёмные (#B58863)
- **Фигуры:** Заглавные буквы — белые, строчные — чёрные
- **Зависимости:** Использует спрайт `src/images/Chess_Pieces_Sprite.svg.png`

### Тестирование

```bash
cd src/builders/fenlib/
go test -v
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

## Использование base64_to_image.gs (Google Docs)

Скрипт `src/base64_to_image.gs` предназначен для использования в Google Docs. Он автоматически преобразует изображения в формате base64 (сгенерированные с флагом `--inline-images`) в реальные изображения в документе.

### Установка и использование

1. Откройте документ Google Docs, куда был вставлен Markdown с base64-изображениями
2. Перейдите в меню **Extensions** → **Apps Script**
3. Скопируйте содержимое файла `src/base64_to_image.gs` в редактор Apps Script
4. Нажмите иконку дискеты (Save) или Ctrl+S
5. Вернитесь в документ и обновите страницу
6. Для запуска скрипта перейдите в **Extensions** → **Apps Script** → выберите функцию `convertAllBase64ToImages` → нажмите **Run**

### Что делает скрипт

- Находит все base64-строки в формате `data:image/...;base64,...` в документе
- Декодирует их в бинарные данные
- Заменяет текстовые base64-строки на реальные изображения (InlineImage)
- Обрабатывает изображения в обратном порядке (от конца к началу), чтобы не сбить индексы

### Примечание

Скрипт работает только с документами Google Docs, содержащими base64-изображения, сгенерированные с помощью флага `--inline-images`.
