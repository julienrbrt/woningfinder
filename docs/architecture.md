# Architecture

This document defines the archtitecture of WoningFinder. Its data schema is found [here](db_schema.png).

WoningFinder is split in 3 components: _WoningFinder-API_, _HousingFinder_ and _HousingMatcher_ and a landing page.

- _[WoningFinder](../cmd/woningfinder-api)_, is serving the different handlers, it serves as API for WoningFinder.nl frontend so the user can register, login to a housing corporation and manage their housing preferences.
- _[HousingFinder](../cmd/housing-finder)_, is used to query all the offers of the housing corporation. It connects them all and query them at the right time and sends its data to redis pub/sub.
- _[HousingMatcher](../cmd/housing-matcher)_, is trigged by _HousingFinder_ via a messaging broker (redis pub/sub). It will match the new offers to the customer search option and react it.

There is as well small tools that are run for special reasons:

- _[db-migration](../cmd/tools/db-initiator)_ permits to initialize the database with default values (housing corporations, cities, housing types, selection methods...) and run the databases migrations.
- _[customer-delete](../cmd/tools/customer-delete)_ permits to delete customers given their email.

## Landing Page

The landing page is available at https://woningfinder.nl.
The source code is available in another [repository](https://github.com/woningfinder/woningfinder.nl).

## WoningFinder-API

Following is a list of endpoint supported by WoningFinder-API. The API works exclusively with JSON. Validation is obviously performed in the frontend and the backend.

| Endpoint Name                    | Method     | Description                                                                          |
| -------------------------------- | ---------- | ------------------------------------------------------------------------------------ |
| /signup                          | POST       | Handles the registration flow                                                        |
| /login [more](###Authentication) | POST       | Login sends an email to the user in order to login                                   |
| /token [more](###Authentication) | POST       | Verifies users token                                                                 |
| /cities                          | GET        | Gets all supported cities                                                            |
| /housing-preferences             | PUT        | Updates the housing preferences of a given user                                      |
| /corporation-credentials         | GET + POST | Manages the different housing credentials for the supported corporation of the user. |

WoningFinder's API is available at https://woningfinder.nl/api.

### Authentication

The authentication works with email one time usage login token.

## Housing-Finder

### Geocoding

A geolocation service is used in order to get the name of a district using coordinates.
This is used because the user can filter the house he wants to react to by city and district.
For that WoningFinder uses Mapbox Geocoding API.

[More information about that API](https://docs.mapbox.com/api/search/geocoding/).

## Housing-Matcher

### Matching

We use redis in order to check if we already try to match a user with an offer. We create an uuid of the user and the address and only check if it does not exists.
This permits to do not have to re-check multiple times an offer as offers stay published for multiple days.

### Security

For reacting to an offer, WoningFinder must authenticate itself as the customer. This means that WoningFinder stores the consumer credentials in the database (`CorporationCredentials`).
Storing it plaintext is obviously not allowed. WoningFinder supports privacy and security of its customers. We use AES encryption to encrypt and store the user password in the datababse. The password is only decrypted to login to the housing corporation with a private key. No plaintext password is ever stored.
