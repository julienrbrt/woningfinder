# Architecture

This document defines the archtitecture of WoningFinder. Its data schema can be found [here](db_schema.png).

WoningFinder is split in 3 components: _WoningFinder-API_, _HousingFinder_ and _HousingMatcher_ and a landing page.

- _[WoningFinder](../cmd/woningfinder-api)_, is serving the different handlers, it serves as API for woningfinder.nl frontend so the user can register, login to a housing corporation and manage their housing preferences.
- _[HousingMatcher](../cmd/housing-matcher)_, is triggered by _HousingFinder_ via a queue (redis list). It will match the new offers to the customer search option and react to it.
- _[Orchestrator](../cmd/orchestrator)_, permits to orchestrate the different jobs that needs to be often ran by WoningFinder.
  - _HousingFinder_ is used to query all the offers of the housing corporation. It connects them all and query them at the right time and sends its data to a redis queue (triggering _HousingMatcher_).
  - _WeeklyUpdate_ generates and send the customer weekly updates.

There is as well small **tools** that are run for special reasons:

- _[db-migrator](../cmd/tools/db-migrator)_ permits to initialize the database with default values (housing corporations, cities, housing types, selection methods...) and run the databases migrations. It is run as a job before every deploy.
- _[customer-delete](../cmd/tools/customer-delete)_ permits to delete customers given their email.
- _[impersonate](../cmd/tools/impersonate)_ permits to get a JWT token for an user in order to impersonate it.

## Landing Page

The landing page is available at https://woningfinder.nl.
The source code is available in another [repository](https://github.com/woningfinder/woningfinder.nl).

## WoningFinder-API

Following is a list of endpoint supported by WoningFinder-API. The API works exclusively with JSON. Validation is obviously performed in the frontend and the backend.

| Endpoint Name               | Method     | Description                                                                       |
| --------------------------- | ---------- | --------------------------------------------------------------------------------- |
| /offering                   | GET        | Gets all supported plans and type housing and cities                              |
| /login                      | POST       | Sends a link to the user in order to log him. The link is valid 6h                |
| /signup                     | POST       | Handles the registration flow                                                     |
| /contact                    | POST       | Handles the contact form to send an email to _contact@woningfinder.nl_            |
| /me                         | GET        | Get all the user information                                                      |
| /me/housing-preferences     | PUT        | Updates the user housing preferences                                              |
| /me/corporation-credentials | GET + POST | Manages the user the different housing credentials for the supported corporation. |
| /stripe-webhook             | POST       | Endpoint where Stripe sends its webhook events (used for validating user payment) |

WoningFinder's API is available at https://woningfinder.nl/api.

### Authentication

The authentication works with JWT. The token are generated in the sent mail and valid for 6h.
One can use the token as header (`Authorization`) and as query parameter (`jwt`).
More information on how built the token in the [code](../internal/auth/jwt.go).

### Payment

The payment is managed by Stripe. We use Stripe Checkout Session in order to redirect the user after signup to a payment page.
Stripe will then confirms that an user has paid via a webhook (_/stripe-webhook_).
The information returned by Stripe must be the user email address and the payment amount.
Our webhook then update the paying information (user and plan) to the concerned user.

More [documentation on how to test the webhook](https://stripe.com/docs/webhooks/test).

## Housing-Finder

### Geocoding

A geolocation service is used in order to get the name of a district using coordinates.
This is used because the user can filter the house he wants to react to by city and district.
For that WoningFinder uses Mapbox Geocoding API.

[More information about that API](https://docs.mapbox.com/api/search/geocoding/).

## Housing-Matcher

### Matching

We use redis in order to check if we already try to match a user with an offer. We create an uuid of the user and the address and only check if it does not exists.
This permits to do not have to re-check multiple times an offer as offers stay published for multiple days. Once there is a match, the match is added in the `HousingPreferencesMatch` table of the database.

People having registered their credentials for the longest get reaction priority, regardless of their plans (probably should be changed on day).

### Corporation Credentials

For reacting to an offer, WoningFinder must authenticate itself as the customer. This means that WoningFinder stores the consumer credentials in the database (`CorporationCredentials`).
WoningFinder supports privacy and security of its customers. We use AES encryption to encrypt and store the user password in the datababse. The password is only decrypted to login to the housing corporation with a private key. No plaintext password is ever stored.
