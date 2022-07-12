# WoningFinder

> WoningFinder is a project that I had during the COVID lockdown, while trying to find a house in the Netherlands.
> It automates the process of finding a rent house by searching and reacting automatically to houses matching your preferences (think of number of bedrooms, city and neighborhood).
> More information about the project can be found on [woningfinder.nl](https://woningfinder.nl) (in Dutch).
>
> Having found my house thanks to this tool, I have decided to make WoningFinder backend open-source.
> If you want it to support more websites feel free to submit a PR that adds a connector for the website.
> I'll integrate it on the [woningfinder.nl](https://woningfinder.nl) as soon as possible.
> Please don't hesitate to submit a PR to this repository if you have any other suggestions or improvements.
> Keep in mind that this project is in its early stage, so the code is not yet optimized nor pretty.

More information can be found in the [docs](docs/).

## Stack

- [Go](https://golang.org)
- [Nuxt.js](https://nuxtjs.org)
- [TailwindCSS](https://tailwindcss.com)
- PostgreSQL
- Redis

## Services

- Sentry
- Mapbox

## Architecture

WoningFinder is split in multiple components:

- _[woningfinder-api](cmd/woningfinder-api)_, is serving the different handlers, it serves as API for woningfinder.nl frontend so the user can register, login to a housing corporation and manage their housing preferences.
- _[housing-matcher](cmd/housing-matcher)_, is triggered by _HousingFinder_ via a queue (redis list). It will match the new offers to the customer search option and react to it.
- _[orchestrator](cmd/orchestrator)_, orchestrates the different jobs that needs to be often ran by WoningFinder.
  - _CleanupUnconfirmedCustomer_ sends a reminder to unconfirmed email user and deletes the customers that did not confirm their email within 72 hours. Runs everyday at 08:00, 16:00.
  - _HousingFinder_ is used to query all the offers of the housing corporation. It connects them all and query them at the right time and sends its data to a redis queue (triggering _HousingMatcher_).
  - _WeeklyUpdate_ generates and send the customer weekly updates. Runs every Friday at 18:00.
  - _CorporationCredentialsMissingReminder_ sends missing corporation credentials reminder to matching customers. Runs everyday at 08:00, 16:00.
- _[housing-finder](cmd/housing-finder)_ replicates the _HousingFinder_ job for a given corporation.
- _[db-migrator](cmd/db-migrator)_ initializes the database with default values (housing corporations, cities, housing types, selection methods...) and run the databases migrations. It is run as a job before every deploy.
- _[city-location-updater](cmd/city-location-updater)_ updates the city location in the database. It is run as a job before every deploy.
- _[impersonate](cmd/impersonate)_ gets a JWT token for an user in order to impersonate it.
- _[customer-delete](cmd/customer-delete)_ deletes customers given their email.

## Issue names

- feature and bug: `Implement [Issue Name] (closes #issue)`
