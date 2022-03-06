## DevOps

### GitHub Actions

We use GitHub Actions in order to test our code. The repository is tested and build every push.
Moreover, in order to support properly the housing corporations, we are testing our implementation everyday at 4 in the morning.

### Environment Variables

The environment variables are loaded from the `.env` first. If not present, it will fallback to the system environment variables.

- _APP_\*\_ defines basis information about the WoningFinder API (name & port)
- _PSQL\_\*_ contains the credentials of the PostgreSQL database
- _REDIS\_\*_ contains the crendentials of the Redis database
- _AES\_SECRET_ contains the encryption key used to make the encrypting of housing corporation credentials more random.

### Local execution

- Fill in the `.env` variables with random data
- **Redis** and **PostgreSQL** must be used via locally environement.

### Useful queries

Find cities to add in WoningFinder

```sql
select *
from corporation_cities cc 
where lower(cc.city_name) in (select lower(name) from cities where cities.latitude is NULL)
```

```sql
select *
from cities c 
where lower(c.name) not in (select lower(cc.city_name) from corporation_cities cc)
```
