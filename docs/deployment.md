# Deployment

## Environment Variables

The environment variables are loaded from the `.env` first. If not present, it will fallback to the system environment variables.

- _APP_\*\_ defines basis information about the WoningFinder API (name & port)
- _PSQL\_\*_ contains the credentials of the PostgreSQL database
- _REDIS\_\*_ contains the crendentials of the Redis database
- _AES_SECRET_ contains the encryption key used to encrypt the user housing corporation credentials. Use `util.AESGenerateKey()` for generating a random key
