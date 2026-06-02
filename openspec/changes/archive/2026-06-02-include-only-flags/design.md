## Context

Сейчас `markdown_builder` поддерживает `--white-except`/`--black-except` (исключение глав по игроку) и `--exam` (включение по точному совпадению). Нужны флаги для включения глав по имени игрока без режима экзамена — `--white-only`/`--black-only`.

## Goals / Non-Goals

**Goals:**
- Добавить `--white-only` — оставить только главы с указанным White
- Добавить `--black-only` — оставить только главы с указанным Black
- Pipe-разделитель для нескольких имён (как в `--white-except`)
- Pipeline: exclusion → inclusion → exam
- `--white-only` + `--black-only` вместе = логическое И (обе стороны должны совпасть)

**Non-Goals:**
- Изменение `splitNames` helper (уже есть, переиспользуем)
- Поддержка регулярных выражений или частичного совпадения
- Взаимодействие с `--split-by-chapters` (работает как и для других фильтров)

## Decisions

1. **Pipeline order:** `--white-except`/`--black-except` (remove) → `--white-only`/`--black-only` (filter in) → `--exam` (final match). Exclusion удаляет нежелательные главы, inclusion фильтрует из оставшихся, exam делает точный матч.

2. **Переиспользование `splitNames`:** Уже реализованный helper для pipe-разделителя используется и для новых флагов — без дублирования.

3. **Логическое И для одновременного использования:** Если указаны и `--white-only`, и `--black-only`, глава проходит только если совпадают оба поля.

## Risks / Trade-offs

- [Risk] Конфликт `--white-only X --white-except X` → exclusion удалит все главы с X, inclusion не найдёт совпадений. Результат: пустой вывод. Это логично — пользователь явно указал противоречивые флаги.
- [Risk] Пустой `--white-only ""` с pipe → `splitNames` вернёт пустой слайс, фильтр не применится. Безопасное поведение.
