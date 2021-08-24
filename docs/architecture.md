# Architecture

This document defines the archtitecture of WoningFinder.

WoningFinder is split in 3 components: _WoningFinder-API_, _HousingFinder_ and _HousingMatcher_ and a landing page.

- _[WoningFinder](../cmd/woningfinder-api)_, is serving the different handlers, it serves as API for woningfinder.nl frontend so the user can register, login to a housing corporation and manage their housing preferences.
- _[HousingMatcher](../cmd/housing-matcher)_, is triggered by _HousingFinder_ via a queue (redis list). It will match the new offers to the customer search option and react to it.
- _[Orchestrator](../cmd/orchestrator)_, permits to orchestrate the different jobs that needs to be often ran by WoningFinder.
  - _CustomerUnconfirmedCleanup_ sends a reminder to unconfirmed email user and deletes the customers that did not confirm their email within 48 hours. Runs everyday at 08:00, 16:00.
  - _CustomerEndFreeTrialReminder_ reminds a free trial customer to pay the plan. Runs everyday at 08:00, 14:00 and 20:00.
  - _HousingFinder_ is used to query all the offers of the housing corporation. It connects them all and query them at the right time and sends its data to a redis queue (triggering _HousingMatcher_).
  - _WeeklyUpdate_ generates and send the customer weekly updates. Runs every Friday at 20:00.

There is as well small **tools** that are run for special reasons:

- _[db-migrator](../cmd/tools/db-migrator)_ permits to initialize the database with default values (housing corporations, cities, housing types, selection methods...) and run the databases migrations. It is run as a job before every deploy.
- _[customer-delete](../cmd/tools/customer-delete)_ permits to delete customers given their email.
- _[impersonate](../cmd/tools/impersonate)_ permits to get a JWT token for an user in order to impersonate it.
- _[housing-finder](../cmd/tools/housing-finder)_ replicates _HousingFinder_ as a command line tool.

## Landing Page

The landing page is available at https://woningfinder.nl.
The source code is available in another [repository](https://github.com/woningfinder/woningfinder.nl).

## WoningFinder-API

Following is a list of endpoint supported by WoningFinder-API. The API works exclusively with JSON. Validation is obviously performed in the frontend and the backend.

| Endpoint Name               | Method     | Description                                                                               |
| --------------------------- | ---------- | ----------------------------------------------------------------------------------------- |
| /offering                   | GET        | Gets all supported plans and type housing and cities                                      |
| /register                   | POST       | Handles the registration flow                                                             |
| /payment                    | POST       | Permits to complete a payment registration (after trial or cancellation)                  |
| /stripe-webhook             | POST       | Endpoint where Stripe sends its webhook events (used for validating user payment)         |
| /crypto-webhook             | POST       | Endpoint where Crypto.com Pay sends its webhook events (used for validating user payment) |
| /login                      | POST       | Sends a link to the user in order to log him. The link is valid 6h                        |
| /me                         | GET        | Get all the user information and confirm user account the first time requested            |
| /me/corporation-credentials | GET + POST | Manages the user the different housing credentials for the supported corporation.         |
| /contact                    | POST       | Handles the contact form to send an email to _contact@woningfinder.nl_                    |
| /waitinglist                | POST       | Handles the city waiting list                                                             |

WoningFinder's API is available at https://woningfinder.nl/api.

### Authentication

The authentication works with JWT. The token are generated in the sent mail and valid for 6h.
One can use the token as header (`Authorization`) and as query parameter (`jwt`).
More information on how built the token in the [code](../internal/auth/jwt.go).

### Payment

The payment is managed by Stripe and by Crypto.com. We use Stripe Checkout Session in order to redirect the user after free trial ended (the user is informed via mail or via its interface).
The PSP will then confirms that an user has paid via a webhook (_/stripe-webhook_ or _/crypto-webhook_).

The information returned by Stripe must be the user email address and the payment amount.
Our webhook then update the plan information of the concerned user.

More [documentation on Stripe webhook](https://stripe.com/docs/webhooks/test).
More [documentation on Crypto.com webhook](https://pay-docs.crypto.com).

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

People having registered their credentials for the longest get reaction priority.

### Corporation Credentials

For reacting to an offer, WoningFinder must authenticate itself as the customer. This means that WoningFinder stores the consumer credentials in the database (`CorporationCredentials`).
WoningFinder supports privacy and security of its customers. We use AES encryption to encrypt and store the user password in the datababse. The password is only decrypted to login to the housing corporation with a private key. No plaintext password is ever stored.
