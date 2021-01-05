# Housing Corporation

Housing Corporation rarely develop their own systems. They use ERPs to achieve the publication and reaction of their housing offers.
WoningFinder is implementing these systems permitting to interact with the housing corporations.

## Supported ERP

- [x] [Itris ERP](https://www.itris.nl/#itris)
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
