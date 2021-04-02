# Housing Corporation

Housing Corporation rarely develop their own systems. They use ERPs to achieve the publication and reaction of their housing offers.
WoningFinder is implementing these systems permitting to interact with the housing corporations.

## Adding Housing Corporation

- Add new housing corporation in database (check documentation [here](https://github.com/WoningFinder/woningfinder/blob/main/docs/architecture.md))
  - Write database migration
  - Update client_provider with its client
  - If a new city is supported add all city districts in the database for that city

## Update Housing Corporation

- New cities supported by the housing corporation are added automatically. Some work must be done however to update their cities in the codebase directly too (see Sentry warning).
- Any other update must made directly in the database and then synchronize with the corporation struct.

## Supported ERP

- [x] [Itris ERP](https://www.itris.nl/#itris) - Itris Connector
- [ ] [Embrace Cloud](https://www.embracecloud.nl/woningcorporaties/wat-kan-het-allemaal/)
- [ ] [Zig](https://zig.nl)
- [ ] [WoningNet WRB](https://www.woningnet.nl)

Some housing corporation (or group of housing corporation) have their home-made system, they are independentely supported:

- [De Woonplaats](http://www.dewoonplaats.nl) - JSON API
- [Woonbureau Almelo](http://www.woonburoalmelo.nl) - JSON API
- Woonkeus Stedendriehoek - JSON API

The definition of corporation are something done offline, once a corporation is supported and a client is created.
The mapping of the corporation and the client is made in the `client_provider`. The matching is done using the name and the url of the housing corporation.

## Supported Housing Corporation

| Housing Corporation | Selection Method                | Running Cities                                                | Implemented |
| ------------------- | ------------------------------- | ------------------------------------------------------------- | ----------- |
| De Woonplaats       | Random, First-Come First-Served | Enschede, Aatlen, Winterswijk, Dinxperlo, Neede, Wehl, Zwolle | ✅          |
| OnsHuis             | Random                          | Enschede, Hengelo                                             |             |
| Domijn              | Random                          | Enschede, Haaksbergen, Losser, Overdinkel                     |             |
