import os
import shutil
import json

def move_files(project_type, project_slug):
    # Определяем исходную директорию (simple или cobra)
    source_dir = os.path.join(project_type, "{{cookiecutter.project_slug}}")
    target_dir = project_slug  # Итоговая директория

    # Копируем файлы из шаблонной директории в итоговую директорию
    for item in os.listdir(source_dir):
        s = os.path.join(source_dir, item)
        d = os.path.join(".", item)  # Копируем в текущую папку проекта
        if os.path.isdir(s):
            shutil.copytree(s, d, dirs_exist_ok=True)
        else:
            shutil.copy2(s, d)

    # Удаляем временную директорию после копирования
    shutil.rmtree(source_dir)

def main():
    # Читаем контекст Cookiecutter (данные из cookiecutter.json)
    with open("cookiecutter.json", "r") as f:
        context = json.load(f)

    project_type = context.get("project_type")
    project_slug = context.get("project_slug")

    # Копируем файлы в зависимости от типа проекта
    move_files(project_type, project_slug)

if __name__ == "__main__":
    main()
