# Deployment

## Environment Variables

The environment variables are loaded from the `.env` first. If not present, it will fallback to the system environment variables.

- _APP_NAME_ is the name of the application (WoningFinder)
- _PSQL\_\*_ contains the credentials of the PostgreSQL database
- _REDIS\_\*_ contains the crendentials of the Redis database
