## API

Following is a list of endpoint supported the WoningFinder API. The API works exclusively with JSON. Validation is obviously performed in the frontend and the backend.

| Endpoint Name               | Method     | Description                                                                               |
| --------------------------- | ---------- | ----------------------------------------------------------------------------------------- |
| /offering                   | GET        | Gets all supported plans and type housing and cities                                      |
| /register                   | POST       | Handles the registration flow                                                             |
| /payment                    | POST       | Permits to complete a payment registration (after trial or cancellation)                  |
| /stripe-webhook             | POST       | Endpoint where Stripe sends its webhook events (used for validating user payment)         |
| /crypto-webhook             | POST       | Endpoint where Crypto.com Pay sends its webhook events (used for validating user payment) |
| /login                      | POST       | Sends a link to the user in order to log him. The link is valid 6h                        |
| /me                         | GET + POST | Get and update all the user information. Confirms user account the first time requested.  |
| /me/corporation-credentials | GET + POST | Manages the user the different housing credentials for the supported corporation.         |
| /me/delete                  | POST       | Let user delete its account                                                               |
| /contact                    | POST       | Handles the contact form to send an email to _contact@woningfinder.nl_                    |
| /waitinglist                | POST       | Handles the city waiting list                                                             |

WoningFinder's API is available at https://woningfinder.nl/api.

### Authentication

The authentication works with JWT. The token are generated in the sent mail and valid for 6h.
One can use the token as header (`Authorization`) and as query parameter (`jwt`).
More information on how built the token in the [code](../internal/auth/jwt.go).

### Payment

The payment is managed by Stripe and by Crypto.com. We use Stripe Checkout Session in order to redirect the user after free trial ended (the user is informed via mail or via its interface).
The PSP will then confirms that an user has paid via a webhook (_/stripe-webhook_ or _/crypto-webhook_).

The information returned by Stripe must be the user email address and the payment amount.
Our webhook then update the plan information of the concerned user.

More [documentation on Stripe webhook](https://stripe.com/docs/webhooks/test).
More [documentation on Crypto.com webhook](https://pay-docs.crypto.com/#integration-guides-checkout-page-redirection).
