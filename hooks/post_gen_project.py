import os

def create_service_structure(project_slug, services, is_cobra):
    cmd_dir = os.path.join(project_slug, "cmd")
    os.makedirs(cmd_dir, exist_ok=True)

    for service in services:
        service_dir = os.path.join(cmd_dir, service)
        os.makedirs(os.path.join(service_dir, "config"), exist_ok=True)
        if is_cobra:
            os.makedirs(os.path.join(service_dir, "loggerbf"), exist_ok=True)

            with open(os.path.join(service_dir, "loggerbf", "logger.go"), "w") as f:
                f.write(f"""package loggerbf

import "fmt"

func InitLogger() {{
    fmt.Println("Initializing logger for {service}")
}}
""")

            with open(os.path.join(service_dir, "loggerbf", "app.go"), "w") as f:
                f.write(f"""package loggerbf

import "fmt"

func StartApp() {{
    fmt.Println("{service} app is running")
}}
""")

        with open(os.path.join(service_dir, "config", "config.go"), "w") as f:
            f.write(f"""package config

import "fmt"

func InitConfig() {{
    fmt.Println("Initializing config for {service}")
}}
""")

def main():
    import json

    with open("cookiecutter.json", "r") as f:
        context = json.load(f)

    project_slug = context.get("project_slug")
    services = context.get("services", [])
    project_type = context.get("project_type")

    is_cobra = project_type == "cobra"
    create_service_structure(project_slug, services, is_cobra)

if __name__ == "__main__":
    main()
