## Housing-Matcher

### Matching

We use redis in order to check if we already try to match a user with an offer. We create an uuid of the user and the address and only check if it does not exists.
This permits to do not have to re-check multiple times an offer as offers stay published for multiple days. Once there is a match, the match is added in the `HousingPreferencesMatch` table of the database.

People having registered their credentials for the longest get reaction priority.

### Corporation credentials

For reacting to an offer, WoningFinder must authenticate itself as the customer. This means that WoningFinder stores the consumer credentials in the database (`CorporationCredentials`).
WoningFinder supports privacy and security of its customers. We use AES encryption to encrypt and store the user password in the datababse. The password is only decrypted to login to the housing corporation with a private key. No plaintext password is ever stored.
