# Housing Corporation

Housing Corporation rarely develop their own systems. They use ERPs to achieve the publication and reaction of their housing offers.
WoningFinder is implementing these systems permitting to interact with the housing corporations.

## How Housing Corporation are supported?

### Supported ERP

- [x] [Itris ERP](https://www.itris.nl/#itris) (via web parsing)
- [ ] [Embrace Cloud](https://www.embracecloud.nl/woningcorporaties/wat-kan-het-allemaal/)
- [ ] [Dynamics Empire by cegeka-dsa](https://www.cegeka-dsa.nl/#intro)
- [ ] [WoningNet WRB](https://www.woningnet.nl) (JSON API)
- [ ] [Zig](https://zig.nl)

Some housing corporation (or group of housing corporation) have their home-made system, they are independentely supported:

- De Woonplaats (JSON API)
- Woonkeus Stedendriehoek (JSON API)

The definition of corporation are something done offline, once a corporation is supported and a client is created.
The mapping of the corporation and the client is made in the `client_provider`. The matching is done using the name and the url of the housing corporation.

## Which Housing Corporation is supported?

| Housing Corporation     | Selection Method                | Running Cities                                                | Implemented |
| ----------------------- | ------------------------------- | ------------------------------------------------------------- | ----------- |
| De Woonplaats           | Random, First-Come First-Served | Enschede, Aatlen, Winterswijk, Dinxperlo, Neede, Wehl, Zwolle | ✅          |
| OnsHuis                 | Random                          | Enschede, Hengelo                                             |             |
| Domijn                  | Random                          | Enschede, Haaksbergen, Losser, Overdinkel                     |             |
| WoningNet               |                                 |                                                               |             |
| Woonkeus Stedendriehoek |                                 | Apeldoorn, Zutphen, Deventer, Twello, Eerbeek, Schalkhaar     |             |
