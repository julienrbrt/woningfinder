# Architechture

## Environment Variables

The environment variables are loaded from the `.env` first. If not present, it will fallback to the system environment variables

- _PSQL\_\*_ contains the credentials of the PostgreSQL database

## How Is Supported Each Housing Corporation

In Enschede:

- _De Woonplaats_ is supported by reverse engineering their API. We make requests to their API as if we are in a browser.
- _OnsHuis_ is supported by using website scrapping, they do not seems to have any API, so we scrape their website in order to get their offers.
- _Domijn_ is as well supported by website scrapping.
