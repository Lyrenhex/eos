# Server API documentation

The Eos server is written in Go, a modern, free, compiled language developed by Google. As a result of this, however, there are specific quirks to the Eos server's handling of certain API calls, owing to the precise nature of Go's `structure` handling.

The Eos server actually operates 2 (in a non-production environment) to 3 (production) servers, specifically:

In a non-production environment:

- A HTTP server operating on port 80 to serve the PWA files (in the `/webclient` folder)
- A WebSocket server operating on a specified port, serving as the back-end of the Eos service.

In a production environment:

- A HTTP server operating on port 80 to redirect requests to the HTTPS server
- A HTTPS server (with valid TLS certificate) operating on port 443 to serve PWA files in an encrypted fashion
- A WebSocket server secured with TLS operating on a specified port, serving as the back-end of the Eos service.

API calls should be sent via a WebSocket packet to the WebSocket server, from the HTTP(S) server operating **on the same host** (for security purposes).

## API Methods

### Connection

Connection to the WebSocket server follows standard WebSocket specification. For example, a JavaScript connection to the WebSocket server (`new WebSocket('address here')`).

Upon successfully connecting to the WebSocket server, the server **must** respond with the following data packet:

```json
{
    "type": "version",
    "data": // version number of the server
}
```

### Logging in (login)

```json
{
    "type": "login",
    "emailAddress": "", // email address goes here
    "password": "", // password goes here
}
```

This method **MUST** be called before calling any subsequent method within the server's functionality; non-logged-in method calls will be discarded by the server automatically.

### Submitting new Mood data (mood)

```json
{
    "type": "mood",
    "day": 0-6, // day of the week, 0 (Sunday) to 6 (Saturday)
    "month": "MM", // month of the year, 0 (January) to 11 (December)
    "year": "yyyy", // current year in UTC time
    "mood": -2-2, // relative mood, -2 to 2
}
```

## Submitting a new comment for a mood (comment)

```json
{
    "type": "comment",
    "mood": -1-1, // integer, either negative (-1), neutral (0), or positive (1)
    "data": "" // string containing comment data
}
```

## Updating user details

```json
{
    "type": "details",
    "emailAddress": "", // new email address
    "password": "", // new password
    "data": "" // new name
}
```

Any values which remain unchanged **MUST** still be sent, but should be sent as an empty string (`""`), to indicate no change is needed.

The only details required for an account is the **email address and password**, the user's name is not required and will default to "friend".

