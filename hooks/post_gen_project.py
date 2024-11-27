import os
import shutil
import json

def move_files(project_type, project_slug):
    source_dir = os.path.join(project_type, "{{cookiecutter.project_slug}}")
    target_dir = project_slug

    for item in os.listdir(source_dir):
        s = os.path.join(source_dir, item)
        d = os.path.join(target_dir, item)
        if os.path.isdir(s):
            shutil.copytree(s, d, dirs_exist_ok=True)
        else:
            shutil.copy2(s, d)

    shutil.rmtree(source_dir)

def main():
    with open("cookiecutter.json", "r") as f:
        context = json.load(f)

    project_type = context.get("project_type")
    project_slug = context.get("project_slug")

    move_files(project_type, project_slug)

if __name__ == "__main__":
    main()
