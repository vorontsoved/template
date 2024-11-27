import os
import subprocess

def initialize_git():
    try:
        subprocess.run(["git", "init"], check=True)
        subprocess.run(["git", "add", "."], check=True)
        subprocess.run(["git", "commit", "-m", "Initial commit"], check=True)
        print("Git repository initialized successfully.")
    except subprocess.CalledProcessError as e:
        print(f"ERROR: Git initialization failed: {e}")
        sys.exit(1)

def install_dependencies():
    try:
        subprocess.run(["go", "mod", "tidy"], check=True)
        print("Dependencies installed successfully.")
    except subprocess.CalledProcessError as e:
        print(f"ERROR: Failed to install Go dependencies: {e}")
        sys.exit(1)

if os.path.exists(".git"):
    print("Git repository already exists. Skipping initialization.")
else:
    initialize_git()

install_dependencies()
