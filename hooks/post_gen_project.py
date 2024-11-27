import os
import shutil

def setup_project():
    project_root = os.getcwd()

    if "{{ cookiecutter.include_cobra }}" == "yes":
        # Переместить содержимое with_cobra в корень
        source_dir = os.path.join(project_root, "cobra")
    else:
        # Переместить содержимое without_cobra в корень
        source_dir = os.path.join(project_root, "simple")

    # Переместить файлы из выбранного шаблона
    for item in os.listdir(source_dir):
        shutil.move(os.path.join(source_dir, item), project_root)

    # Удалить временные папки
    shutil.rmtree(os.path.join(project_root, "cobra"), ignore_errors=True)
    shutil.rmtree(os.path.join(project_root, "simple"), ignore_errors=True)

setup_project()
