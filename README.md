# WoningFinder

## Stack

- [Go](https://golang.org)
- [Nuxt.js](https://nuxtjs.org)
- [TailwindCSS](https://tailwindcss.com)
- PostgreSQL

## Services

- Sentry
- Mapbox

## Architecture

WoningFinder is split in multiple components:

- **api**, is serving the different handlers, it serves as API for woningfinder.nl frontend so the user can register, login to a housing corporation and manage their housing preferences.
- **matcher**, is go routine listenning to a channel to match the new offers to the customer search option and react to it.
- **cron jobs**:
  - _CleanupUnconfirmedCustomer_ sends a reminder to unconfirmed email user and deletes the customers that did not confirm their email within 72 hours. Runs everyday at 08:00, 16:00.
  - _HousingFinder_ is used to query all the offers of the housing corporation. It connects them all and query them at the right time and sends its data a go channel (listened by the matcher).
  - _WeeklyUpdate_ generates and send the customer weekly updates. Runs every Friday at 18:00.
  - _CorporationCredentialsMissingReminder_ sends missing corporation credentials reminder to matching customers. Runs everyday at 08:00, 16:00.
- _[db-migrator](cmd/db-migrator)_ initializes the database with default values (housing corporations, cities, housing types, selection methods...) and run the databases migrations. It is run as a job before every deploy.
- _[city-location-updater](cmd/city-location-updater)_ updates the city location in the database. It is run as a job before every deploy.
- _[impersonate](cmd/impersonate)_ gets a JWT token for an user in order to impersonate it.
- _[customer-delete](cmd/customer-delete)_ deletes customers given their email.

More information can be found in the [docs](docs/).
