
## Housing-Finder

### Geocoding

A geolocation service is used in order to get the name of a district using coordinates.
This is used because the user can filter the house he wants to react to by city and district.
For that WoningFinder uses Mapbox Geocoding API.

[More information about that API](https://docs.mapbox.com/api/search/geocoding/).

### Adding Housing Corporation

- Add new housing corporation in database (check documentation [here](https://github.com/WoningFinder/woningfinder/blob/main/docs/architecture.md))
  - Write database migration
  - Update client_provider with its client
  - If a new city is supported add all suggested city districts for that city and add the city in the city table


**Cities**
Suggested city districts are useful when Mapbox does not get the main neighbourhood but only its descendent (example what people call in Enschede "Centrum" is composed of different districts, the suggested districts let the user select them all instead of one by one).

City coordinates can be found via https://developer.mapquest.com/documentation/tools/latitude-longitude-finder.

### Update Housing Corporation

- New cities supported by the housing corporation are added automatically. Some work must be done however to update their cities in the codebase directly too (see Sentry warning).
- Any other update must made directly in the database and then synchronize with the corporation struct.

### Supported Housing Corporation

Housing Corporations are supported in WoningFinder by _connectors_. Those connector are based in the ERPs that the housing corporations use. This permits to make a connector be compabible with multiples housing corporations. WoningFinder supports for now the following:

- [x] [Itris ERP](https://www.itris.nl/#itris) - Itris Connector
- [x] [Zig](https://zig.nl) - Zig Connector
- [x] [WoningNet WRB](https://www.woningnet.nl)

Some housing corporation (or group of housing corporation) sometimes have their homemade system, the _connectors_ supports then that system and no specific ERPs:

- [De Woonplaats](http://www.dewoonplaats.nl) - JSON API
- [Woonbureau Almelo](http://www.woonburoalmelo.nl) - JSON API
- Woonkeus Stedendriehoek - JSON API

The definition of corporation is something done offline, once a corporation is supported and a client is created.
The mapping of the corporation and the client is made in the `client_provider` using the name of the housing corporation.
