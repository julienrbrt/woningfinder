## API

Following is a list of endpoint supported the WoningFinder API. The API works exclusively with JSON. Validation is obviously performed in the frontend and the backend.

| Endpoint Name               | Method     | Description                                                                              |
| --------------------------- | ---------- | ---------------------------------------------------------------------------------------- |
| /offering                   | GET        | Gets all supported housing types and cities                                              |
| /register                   | POST       | Handles the registration flow                                                            |
| /login                      | POST       | Sends a link to the user in order to log him. The link is valid 6h                       |
| /me                         | GET + POST | Get and update all the user information. Confirms user account the first time requested. |
| /me/housing-preferences     | POST       | Update the user husing preferences                                                       |
| /me/corporation-credentials | GET + POST | Manages the user the different housing credentials for the supported corporation.        |
| /me/delete                  | POST       | Let user delete its account                                                              |
| /contact                    | POST       | Handles the contact form to send an email to _contact@woningfinder.nl_                   |
| /waitinglist                | POST       | Handles the city waiting list                                                            |

WoningFinder's API is available at https://woningfinder.nl/api.

### Authentication

The authentication works with JWT. The token are generated in the sent mail and valid for 6h.
One can use the token as header (`Authorization`) and as query parameter (`jwt`).
More information on how built the token in the [code](../internal/auth/jwt.go).