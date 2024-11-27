import re
import sys

def validate_project_slug(project_slug):
    if not re.match(r'^[a-z][a-z0-9_]+$', project_slug):
        print(f"ERROR: Project slug '{project_slug}' is invalid. It must start with a letter and contain only lowercase letters, numbers, and underscores.")
        sys.exit(1)

project_slug = "{{ cookiecutter.project_slug }}"
validate_project_slug(project_slug)
