import json

def main():
    with open("cookiecutter.json", "r") as f:
        context = json.load(f)

    project_type = context.get("project_type")

    if project_type not in ["simple", "cobra"]:
        raise ValueError("Выберите корректный тип проекта: simple или cobra")

if __name__ == "__main__":
    main()
