# Architechture

This document defines the archtitecture of WoningFinder.

WoningFinder is splits in two components: _HousingFinder_ and _HousingMatcher_.

- HousingFinder, is used to query all the offers of the housing corporation. It connects them all and query them at the right time.
- HousingMatcher, is trigged after HousingFinder via a messaging broker. It will match the new offer to the customer research option and apply to the right one.

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
