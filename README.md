# WoningFinder

More information can be found here [docs](docs/).

## Architecture

WoningFinder is split in multiple components:

- _[woningfinder-api](cmd/woningfinder-api)_, is serving the different handlers, it serves as API for woningfinder.nl frontend so the user can register, login to a housing corporation and manage their housing preferences.
- _[housing-matcher](cmd/housing-matcher)_, is triggered by _HousingFinder_ via a queue (redis list). It will match the new offers to the customer search option and react to it.
- _[orchestrator](cmd/orchestrator)_, permits to orchestrate the different jobs that needs to be often ran by WoningFinder.
  - _CleanupUnconfirmedCustomer_ sends a reminder to unconfirmed email user and deletes the customers that did not confirm their email within 72 hours. Runs everyday at 08:00, 16:00.
  - _WeeklyUpdate_ generates and send the customer weekly updates. Runs every Friday at 18:00.
- _[impersonate](cmd/impersonate)_ permits to get a JWT token for an user in order to impersonate it.
- _[customer-delete](cmd/customer-delete)_ permits to delete customers given their email.
- _[db-migrator](cmd/db-migrator)_ permits to initialize the database with default values (housing corporations, cities, housing types, selection methods...) and run the databases migrations. It is run as a job before every deploy.
- _[city-location-updater](cmd/city-location-updater)_ permits to update the city location in the database.

## Issue names

**feature and bug**

- `Implement [Issue Name] (closes #issue)`

**bug on sentry**

- `Implement [Issue Name] (fixes #sentry-name and closes #issue)`
