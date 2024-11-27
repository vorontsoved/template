import os
import shutil
import json

def move_files(project_type, project_slug):
    """Копирует файлы из папки simple или cobra в корень"""
    source_dir = os.path.join(project_type, "{{cookiecutter.project_slug}}")
    target_dir = "."

    if not os.path.exists(source_dir):
        raise ValueError(f"Не найдена папка-источник: {source_dir}")

    for item in os.listdir(source_dir):
        s = os.path.join(source_dir, item)
        d = os.path.join(target_dir, item)
        if os.path.isdir(s):
            shutil.copytree(s, d, dirs_exist_ok=True)
        else:
            shutil.copy2(s, d)

    # Удаляем временные файлы
    shutil.rmtree(source_dir)

def main():
    # Читаем настройки проекта из cookiecutter.json
    with open("cookiecutter.json", "r") as f:
        context = json.load(f)

    project_type = context.get("project_type")
    project_slug = context.get("project_slug")

    if project_type not in ["simple", "cobra"]:
        raise ValueError("Некорректный тип проекта. Выберите simple или cobra.")

    move_files(project_type, project_slug)

if __name__ == "__main__":
    main()
