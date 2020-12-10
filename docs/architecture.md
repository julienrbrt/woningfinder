# Architechture

This document defines the archtitecture of WoningFinder.

## Supported ERP

- [Itris ERP](https://www.itris.nl/#itris)
- [Dynamics Empire by cegeka-dsa](https://www.cegeka-dsa.nl/#intro)
- [WoningNet WRB](https://www.woningnet.nl)

Some housing corporation (or group of housing corporation) have their home-made system, they are independentely supported:

- De Woonplaats
- Woonkeus Stedendriehoek

## Environment Variables

The environment variables are loaded from the `.env` first. If not present, it will fallback to the system environment variables

- _NAME_ is the name of the application (WoningFinder)
- _PSQL\_\*_ contains the credentials of the PostgreSQL database

## Location

A Geolocation service is used in order to get the name of a district using coordinates.
This is used because the user can filter the house he wants to react to by city and district.
For that WoningFinder uses OpenStreetMap API.

More information about that API [here](https://nominatim.openstreetmap.org).
