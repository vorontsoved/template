import json

def prompt_for_simple_service():
    print("Вы выбрали тип проекта 'simple'. Введите название единственного сервиса:")
    while True:
        service_name = input("Название сервиса: ").strip()
        if service_name:
            return [service_name]
        print("Название сервиса не может быть пустым. Попробуйте снова.")

def prompt_for_cobra_services():
    print("Вы выбрали тип проекта 'cobra'. Введите названия сервисов по одному. Нажмите Enter без ввода, чтобы закончить.")
    services = []
    while True:
        service_name = input("Название сервиса: ").strip()
        if not service_name:
            break
        if service_name in services:
            print(f"Сервис '{service_name}' уже добавлен. Попробуйте другое имя.")
        else:
            services.append(service_name)
    if not services:
        print("Вы не ввели ни одного сервиса. Минимум один сервис обязателен.")
        return prompt_for_cobra_services()
    return services

def main():
    with open("cookiecutter.json", "r") as f:
        context = json.load(f)

    project_type = context.get("project_type")
    
    if project_type == "simple":
        context["services"] = prompt_for_simple_service()
    elif project_type == "cobra":
        context["services"] = prompt_for_cobra_services()

    with open("cookiecutter.json", "w") as f:
        json.dump(context, f, indent=4)

if __name__ == "__main__":
    main()
