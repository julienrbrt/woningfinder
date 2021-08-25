# Email

The major interaction that WoningFinder has with its customers is via email.
The distribution of the email must thus be perfect and their design as well.

The templates are written in MJML and are generated to html with `mjml file.mjml -o file.html`.
Note that when were this `range` or `if` condition in the mjml, it has to be re-added manually in the html after each re-generation.

## Lists

### Action needed

- Welcome (activate account via connecting)
- Login
- Invalid credentials warning
- No credentials available
- Free trial ended reminder
- Email confirmation reminder

## Information

- Weekly update
- Weekly update empty
- Payment confirmation
- Goodbye email

## Configuration

Email are sent from `contact@woningfinder.nl` via SMTP using the following enviroment variable:

```
EMAIL_ADDRESS=''
EMAIL_PASSWORD=''
EMAIL_SMTP_ADDRESS=''
EMAIL_SMTP_PORT=
```
