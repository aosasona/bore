import os
import re
import sys

MIGRATION_DIR = "internal/database/migrations"


class Migration:
    def __init__(self, name):
        name = re.sub(r"[^a-zA-Z0-9_]", "_", name)
        name = name.lower()
        self.name = name

    def get_migration_name(self):
        migration_files = self.get_migration_files()
        last_migration = migration_files[-1] if len(migration_files) > 0 else 0
        last_idx = int(last_migration.split("_")[0]) if last_migration else 0
        new_idx = str(last_idx + 1).zfill(4)
        file_name = f"{new_idx}_{self.name}"
        return file_name

    def create(self):
        migration_name = self.get_migration_name()
        files_to_create = [f"{migration_name}.up.sql"]

        # Check if the user wanted a down migration
        if len(sys.argv) > 2 and sys.argv[2] == "--down":
            files_to_create.append(f"{migration_name}.down.sql")

        for file in files_to_create:
            if os.path.exists(f"{MIGRATION_DIR}/{file}"):
                print(f"Migration {file} already exists.")
                continue
            with open(f"{MIGRATION_DIR}/{file}", "w") as f:
                f.write("-- Write your migration here")
                print(f"Created migration {file}")

    def get_migration_files(self):
        # Get all migration files from the migrations directory
        files = os.listdir(path=MIGRATION_DIR)
        files = [file for file in files if file.endswith(".up.sql")]
        files = [file.replace(".up.sql", "") for file in files]
        files = sorted(files)
        return files if files else []


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python create_migration.py <migration_name>")
        sys.exit(1)
    migration_name = sys.argv[1]
    migration = Migration(migration_name)
    migration.create()
