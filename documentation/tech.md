# Technical Specification

Data is saved and stored to individual files within the `data` folder. Configuration data for the server is stored in `config.json` within the `data` folder.

`config.json` should contain:

```javascript
{
    "envProduction": true, // true if production environment, false otherwise
    "envKey": "", // File path to TLS key. REQUIRED IF PRODUCTION.
    "envCertificate": "", // File path to TLS certificate. REQUIRED IF PRODUCTION.
    "srvHostname": "", // Server's hostname. This should be the domain name the server is running on. Used for security and URLs.
    "srvPort": 9874, // Sets the port that the WebSocket serber will run on. Modifying this will require changing the client to account for it.
    "googleApiKey": "", // API key for Google Cloud if using Perspective API. Not required.
    "discordWebhook": "", // URI for a Discord Webhook endpoint. The official Eos server uses Discord webhooks for notifying staff of reports; not needed on custom setups.
    "sendgridApiKey": "", // API key for SendGrid.com for email-related tasks (address verification, password resets, etc).
    "sendgridApiAuth": "", // ID for the transactional template for the email verification email.
    "sendgridApiReset": "", // ID for the transactional template for the password reset email.
    "sendgridAddress": "", // the email address to send the email from.
}
```

- Eos instances in a production environment must be operating over TLS end-to-end.
- Passwords must be stored in a secure manner. We advise hashing and salting, using a cryptographically-secure hashing algorithm (ie. not MD5 or SHA1).

## Notices for server administrators

Eos does not offer built-in API calls to modify the `Admin` flag on an account, as this introduces *potential insecurities*. This action is expected to be completed manually, through modifying the user's `json` file and setting `admin` to `true`. A map of email addresses and user IDs is available in `users.json` if required. A manual modification does **NOT** require a server reboot to take effect, but will take effect on the affected user's next login (existing sessions are not affected).

Non-official instances are not affiliated with The Eos Project in any way, shape, or form, and suggesting as such is misinformation. The Eos Project will never recommend a non-official instance that it cannot guarantee follows security best-practices and follows the Eos Privacy Policy.