import os
import shutil

def setup_project():
    project_root = os.getcwd()

    if "{{ cookiecutter.include_cobra }}" == "yes":
        source_dir = os.path.join(project_root, "cobra")
    else:
        source_dir = os.path.join(project_root, "simple")

    for item in os.listdir(source_dir):
        shutil.move(os.path.join(source_dir, item), project_root)

    shutil.rmtree(os.path.join(project_root, "cobra"), ignore_errors=True)
    shutil.rmtree(os.path.join(project_root, "simple"), ignore_errors=True)
    shutil.rmtree(os.path.join(project_root, "{{cookiecutter.project_slug}}"), ignore_errors=True)

setup_project()
