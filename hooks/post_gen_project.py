import os
import shutil

def setup_project():
    project_root = os.getcwd()

    if "{{ cookiecutter.include_cobra }}" == "yes":
        source_dir = os.path.join(project_root, "cobra")
    else:
        source_dir = os.path.join(project_root, "simple")

    # Переместить только содержимое выбранного шаблона
    for item in os.listdir(source_dir):
        source_item = os.path.join(source_dir, item)
        target_item = os.path.join(project_root, item)

        if os.path.exists(target_item):
            # Удалить существующий файл/папку
            if os.path.isfile(target_item) or os.path.islink(target_item):
                os.unlink(target_item)
            elif os.path.isdir(target_item):
                shutil.rmtree(target_item)

        # Переместить содержимое в корень проекта
        shutil.move(source_item, target_item)

    # Удалить временные папки
    shutil.rmtree(os.path.join(project_root, "cobra"), ignore_errors=True)
    shutil.rmtree(os.path.join(project_root, "simple"), ignore_errors=True)

setup_project()
