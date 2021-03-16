# db-migrator

db-migrator is a tool that permits us to migrate our PostgreSQL database at WoningFinder.
We use the [go-pg](https://github.com/go-pg/pg) for interaction with our database and its migration tool: [migrations](https://github.com/go-pg/migrations).

## Usage

- Write a migration in the format `X_migration_name.go`.
