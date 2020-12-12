# Architechture

This document defines the archtitecture of WoningFinder. Its data schema is found [here](db_schema.png).

WoningFinder is splits in two components: _HousingFinder_ and _HousingMatcher_.

- HousingFinder, is used to query all the offers of the housing corporation. It connects them all and query them at the right time.
- HousingMatcher, is trigged after HousingFinder via a messaging broker. It will match the new offer to the customer research option and apply to the right one.

There is as well small scrips that are run for special reasons:

- _db-initiator_ permits to fills the default values in the database (housing corporations, cities, housing types, selection methods...).
- _add-customer_ permits to create a customer that looks for a house in WoningFinder. This is a temporary commands, during the time WoningFinder is with limited availability, where the users are created manually via this cmd.

## HousingFinder

### Supported ERP

- [Itris ERP](https://www.itris.nl/#itris)
- [Dynamics Empire by cegeka-dsa](https://www.cegeka-dsa.nl/#intro)
- [WoningNet WRB](https://www.woningnet.nl) (JSON API)

Some housing corporation (or group of housing corporation) have their home-made system, they are independentely supported:

- De Woonplaats (JSON API)
- Woonkeus Stedendriehoek (JSON API)

## HousingMatcher

### Location Provider

A Geolocation service is used in order to get the name of a district using coordinates.
This is used because the user can filter the house he wants to react to by city and district.
For that WoningFinder uses OpenStreetMap API.

More information about that API [here](https://nominatim.openstreetmap.org).

### Security

For applying to an offer, WoningFinder must authenticate itself as the customer. This means that WoningFinder stores the consumer credentials in the database (`CorporationCredentials`).
Storing it plaintext is obviously not allowed. WoningFinder supports privacy and security of its customers. We use AES encryption to encrypt and store the user password in the datababse. The password is only decrypted to login to the housing corporation with a private key. No plaintext password is ever stored.
