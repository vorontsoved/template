import os
import shutil

def setup_project():
    project_root = os.getcwd()

    if "{{ cookiecutter.include_cobra }}" == "yes":
        source_dir = os.path.join(project_root, "cobra")
    else:
        source_dir = os.path.join(project_root, "simple")

    # Переместить содержимое выбранной папки в корень
    for item in os.listdir(source_dir):
        shutil.move(os.path.join(source_dir, item), project_root)

    # Удалить временные папки
    shutil.rmtree(os.path.join(project_root, "cobra"), ignore_errors=True)
    shutil.rmtree(os.path.join(project_root, "simple"), ignore_errors=True)

setup_project()
