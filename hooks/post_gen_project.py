import os
import shutil
import subprocess

def clean_target_directory():
    project_root = os.getcwd()
    for item in os.listdir(project_root):
        # Исключаем временные папки, чтобы случайно не удалить файлы шаблона
        if item not in ["cookiecutter.json", "hooks"]:
            item_path = os.path.join(project_root, item)
            if os.path.isfile(item_path):
                os.remove(item_path)
            elif os.path.isdir(item_path):
                shutil.rmtree(item_path)

def copy_files(source_dir, destination_dir):
    for item in os.listdir(source_dir):
        src_path = os.path.join(source_dir, item)
        dest_path = os.path.join(destination_dir, item)
        if os.path.isdir(src_path):
            shutil.copytree(src_path, dest_path)
        else:
            shutil.copy2(src_path, dest_path)

def run_go_mod_tidy():
    try:
        subprocess.run(["go", "mod", "tidy"], check=True)
    except FileNotFoundError:
        print("Ошибка: Команда `go` не найдена. Убедитесь, что Go установлен.")
    except subprocess.CalledProcessError as e:
        print(f"Ошибка при выполнении `go mod tidy`: {e}")

def setup_project():
    project_root = os.getcwd()
    clean_target_directory()

    # Определяем, какой шаблон использовать
    if "{{ cookiecutter.include_cobra }}" == "yes":
        source_dir = os.path.join(project_root, "cobra")
    else:
        source_dir = os.path.join(project_root, "simple")

    # Копируем содержимое выбранной папки
    copy_files(source_dir, project_root)

    # Удаляем временные папки шаблона
    shutil.rmtree(os.path.join(project_root, "cobra"), ignore_errors=True)
    shutil.rmtree(os.path.join(project_root, "simple"), ignore_errors=True)

    # Выполняем go mod tidy
    run_go_mod_tidy()

setup_project()
