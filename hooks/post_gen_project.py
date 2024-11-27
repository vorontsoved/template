import os
import shutil

def setup_project():
    project_root = os.getcwd()

    if "{{ cookiecutter.include_cobra }}" == "yes":
        source_dir = os.path.join(project_root, "cobra")
    else:
        source_dir = os.path.join(project_root, "simple")

    # Удалить содержимое корня проекта
    for item in os.listdir(project_root):
        item_path = os.path.join(project_root, item)
        if os.path.isfile(item_path) or os.path.islink(item_path):
            os.unlink(item_path)  # Удалить файл или символическую ссылку
        elif os.path.isdir(item_path):
            shutil.rmtree(item_path)  # Удалить директорию

    # Переместить только содержимое из выбранного шаблона
    for item in os.listdir(source_dir):
        shutil.move(os.path.join(source_dir, item), project_root)

    # Удалить временные папки
    shutil.rmtree(os.path.join(project_root, "cobra"), ignore_errors=True)
    shutil.rmtree(os.path.join(project_root, "simple"), ignore_errors=True)

setup_project()
