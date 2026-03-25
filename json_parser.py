#!/usr/bin/env python3
"""
json_to_html.py — извлекает HTML из поля "content" в JSON-файле
и сохраняет его как .html файл.

Использование:
    python json_to_html.py input.json [output.html]

Если output.html не указан, файл сохраняется рядом с input.json
с тем же именем, но расширением .html.
"""

import json
import sys
import os
import argparse


def find_content(obj, path="root"):
    """
    Рекурсивно обходит JSON-объект и возвращает список всех найденных
    значений поля с ключом "content".
    Возвращает [(путь, значение), ...]
    """
    results = []

    if isinstance(obj, dict):
        for key, value in obj.items():
            current_path = f"{path}.{key}"
            if key == "content":
                results.append((current_path, value))
            # Продолжаем искать глубже, даже если нашли content
            results.extend(find_content(value, current_path))
    elif isinstance(obj, list):
        for i, item in enumerate(obj):
            results.extend(find_content(item, f"{path}[{i}]"))

    return results


def looks_like_html(value):
    """Проверяет, похоже ли строковое значение на HTML."""
    if not isinstance(value, str):
        return False
    stripped = value.strip()
    return stripped.startswith("<") and ">" in stripped


def main():
    parser = argparse.ArgumentParser(
        description="Извлекает HTML из поля 'content' в JSON-файле."
    )
    parser.add_argument("input", help="Путь к входному JSON-файлу")
    parser.add_argument(
        "output",
        nargs="?",
        help="Путь к выходному HTML-файлу (необязательно)",
    )
    parser.add_argument(
        "--all",
        action="store_true",
        help="Если найдено несколько полей content, сохранить все (добавит суффикс _1, _2, ...)",
    )
    parser.add_argument(
        "--index",
        type=int,
        default=0,
        help="Индекс (с 0) нужного поля content, если их несколько (по умолчанию 0 — первое)",
    )
    args = parser.parse_args()

    # --- Читаем JSON ---
    if not os.path.isfile(args.input):
        print(f"[Ошибка] Файл не найден: {args.input}", file=sys.stderr)
        sys.exit(1)

    try:
        with open(args.input, "r", encoding="utf-8") as f:
            data = json.load(f)
    except json.JSONDecodeError as e:
        print(f"[Ошибка] Невалидный JSON: {e}", file=sys.stderr)
        sys.exit(1)

    # --- Ищем все поля content ---
    found = find_content(data)

    if not found:
        print("[Ошибка] Поле 'content' не найдено в JSON.", file=sys.stderr)
        sys.exit(1)

    # Оставляем только те, что похожи на HTML
    html_candidates = [(path, val) for path, val in found if looks_like_html(val)]

    if not html_candidates:
        # Если HTML-подобных не нашли — берём все content как строки
        html_candidates = [(path, val) for path, val in found if isinstance(val, str)]

    if not html_candidates:
        print(
            "[Ошибка] Поле 'content' найдено, но его значение не является строкой/HTML.",
            file=sys.stderr,
        )
        print("Найденные поля:", file=sys.stderr)
        for path, val in found:
            print(f"  {path}: {type(val).__name__}", file=sys.stderr)
        sys.exit(1)

    # --- Определяем базовое имя выходного файла ---
    base_output = args.output or os.path.splitext(args.input)[0] + ".html"

    def save_html(filepath, html_content):
        with open(filepath, "w", encoding="utf-8") as f:
            f.write(html_content)
        print(f"[OK] HTML сохранён в: {filepath}")

    # --- Сохраняем ---
    if args.all and len(html_candidates) > 1:
        base, ext = os.path.splitext(base_output)
        ext = ext or ".html"
        for i, (path, val) in enumerate(html_candidates):
            out_path = f"{base}_{i + 1}{ext}"
            print(f"  Поле: {path}")
            save_html(out_path, val)
    else:
        if len(html_candidates) > 1:
            print(f"[Внимание] Найдено {len(html_candidates)} полей 'content'. Используется индекс {args.index}.")
            print("Все найденные пути:")
            for i, (path, _) in enumerate(html_candidates):
                marker = " <-- выбрано" if i == args.index else ""
                print(f"  [{i}] {path}{marker}")
            print("Используйте --index N для выбора нужного или --all для сохранения всех.")

        try:
            chosen_path, chosen_val = html_candidates[args.index]
        except IndexError:
            print(
                f"[Ошибка] Индекс {args.index} вне диапазона (найдено {len(html_candidates)} полей).",
                file=sys.stderr,
            )
            sys.exit(1)

        print(f"  Поле: {chosen_path}")
        base, ext = os.path.splitext(base_output)
        out_path = base_output if ext else base_output + ".html"
        save_html(out_path, chosen_val)


if __name__ == "__main__":
    main()