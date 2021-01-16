# db-migration

db-migration is a tool that permits us to migrate our PostgreSQL database at WoningFinder.
We use the [go-pg](https://github.com/go-pg/pg) for interaction with our database and its migration tool: [migrations](https://github.com/go-pg/migrations).

## Usage

- Write a migration in the format `X_migration_name.go`. Migrations can be as well written in SQL directly.
- Run `go run . up` in order to run your migration
