# Technical Specification

Data is saved and stored to individual files within the `data` folder. Configuration data for the server is stored in `config.json` within the `data` folder.

`config.json` should contain:

```javascript
{
    "envProduction": true, // true if production environment, false otherwise
    "envKey": "", // File path to TLS key. REQUIRED IF PRODUCTION.
    "envCert": "", // File path to TLS certificate. REQUIRED IF PRODUCTION.
    "srvHostname": "", // Server's hostname. This should be the domain name the server is running on. Used for security and URLs.
    "srvPort": 9874, // Sets the port that the WebSocket serber will run on. Modifying this will require changing the client to account for it.
    "googleApiKey": "", // API key for Google Cloud if using Perspective API. Not required.
    "discordWebhook": "" // URI for a Discord Webhook endpoint. The official Eos server uses Discord webhooks for notifying staff of reports; not needed on custom setups.
}
```

- Eos instances in a production environment must be operating over TLS end-to-end.
- Passwords must be stored in a secure manner. We advise hashing and salting, using a cryptographically-secure hashing algorithm (ie. not MD5 or SHA1).

Non-official instances are not affiliated with The Eos Project in any way, shape, or form, and suggesting as such is misinformation. The Eos Project will never recommend a non-official instance that it cannot guarantee follows security best-practices and follows the Eos Privacy Policy.