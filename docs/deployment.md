# Deployment

## GitHub Actions

We use GitHub Actions in order to test our code. The repository is tested and build every push.
Moreover, in order to support properly the housing corporations, we are testing our implementation everyday at 4 in the morning.

## Environment Variables

The environment variables are loaded from the `.env` first. If not present, it will fallback to the system environment variables.

- _APP_\*\_ defines basis information about the WoningFinder API (name & port)
- _PSQL\_\*_ contains the credentials of the PostgreSQL database
- _REDIS\_\*_ contains the crendentials of the Redis database
- _AES_SECRET_ contains the encryption key used to encrypt the user housing corporation credentials. Use `util.AESGenerateKey()` for generating a random key
