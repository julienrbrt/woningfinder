# Emails

The major interaction that WoningFinder has with its customers is via email.
The distribution of the email must thus be perfect and their design as well.

## Usage

Email are sent from `contact@woningfinder.nl` via SMTP using the following enviroment variable:

```
EMAIL_ADDRESS=''
EMAIL_PASSWORD=''
EMAIL_SMTP_ADDRESS=''
EMAIL_SMTP_PORT=
```

## Templates

The templates can be found in the [notifications service](internal/services/notifications) as well as below:

**Welcome**

```
## Welkom bij WoningFinder!

**Je zoekopdracht is ingesteld.**

Om voor jou te kunnen reageren, hoef je alleen maar in te loggen bij de woningcorporaties waar je wilt reageren.

We hebben voor jou alleen de woningcorporaties geselecteerd die met jouw zoekopdracht matchen.

We raden je aan bij elke in te loggen zodat je sneller een huis kunt vinden.

<Mijn woningcorporaties>

Dan kan je ontspannen, we reageren voor jou. Je hoeft verder niets meer te doen.

Hulp nodig of vragen? Antwoord deze e-mail, we helpen je graag.

Groetjes,
Team WoningFinder

Als de _Mijn woningcorporaties_ knop niet voor jou werkt, je kan de URL hieronder klikken of kopiëren.

_https://woningfinder.com/woningcorporaties?jwt=_
```

**Weekly update**

```
## Wekelijkse update

Hallo <naam>,

We hebben **goed nieuws**: in de afgelopen week hebben we op <number> woningen gereageerd:

| reactie datum | adres | woningcorporatie |
| ------------- | ----- | ---------------- |
| XXX           | XXX   | XXX              |
| XXX           | XXX   | XXX              |

Voor meer informatie, kun je altijd kijken op de website van de woningcorporaties waar we hebben gereageerd.
Je huis staat tussen jouw reacties.

We hopen dat je word gekozen voor een van deze woningen!

Hulp nodig of vragen? Antwoord deze e-mail, we helpen je graag.

Groetjes,
Team WoningFinder

---

## Wekelijkse update

Hoi <naam>,

We hebben elke dag gekeken, maar hebben deze week niets voor jou kunnen vinden.
Maak je geen zorgen, we blijven zoeken!

Hulp nodig of vragen? Antwoord deze e-mail, we helpen je graag.
Tot volgende week.

Groetjes,
Team WoningFinder
```

**Corporation Credentials**

```
## Er is iets misgegaan met je inloggegevens

Hoi <naam>!

We hebben geprobeerd om in te loggen in <woningcorporatie>, maar het lijkt erop dat je inloggegevens niet meer kloppen (je hebt waarschijnlijk je wachtwoord veranderd).
Ze zijn nu dus verwijderd van ons systeem.

Als je nog steeds wilt dat we reageren op het <woningcorporatie> aanbod, log dan in via <link> om je inloggegevens voor deze woningcorporatie
opnieuw in te stellen.

Hulp nodig of vragen? Antwoord deze e-mail, we helpen je graag.

Groetjes,
Team WoningFinder
```

**Bye**

```
## Hoera je hebt een huis gevonden!

Hallo <naam>,

We hebben gezien dat je een huis hebt gevonden. Van harte gefeliciteerd!

Omdat je ons niet meer nodig hebt, zijn al je gegevens van ons systeem verwijderd (we zijn pricacy-freaks, weet je nog).

Bedankt voor je vertrouwen in ons en geniet van jouw nieuwe woning.

Groetjes,
Team WoningFinder
```
