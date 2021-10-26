## Email

The major interaction that WoningFinder has with its customers is via email.
The distribution of the email must thus be perfect and their design as well.

The templates are written in MJML and are generated to html with `mjml file.mjml -o file.html` or using the script `generate-html.sh` in the `email/templates` folder.
Note that when were this `range` or `if` condition in the mjml, it has to be re-added manually in the html after each re-generation.

### Configuration

Email are sent using [Postmark](https://postmarkapp.com) email service.