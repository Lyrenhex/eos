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

```javascript
{
    "type": "version",
    "data": // version number of the server
}
```

### Logging in (login)

```javascript
{
    "type": "login",
    "emailAddress": "", // email address goes here
    "password": "", // password goes here
}
```

This method **MUST** be called before calling any subsequent method within the server's functionality; non-logged-in method calls will be discarded by the server automatically.

### Submitting new Mood data (mood)

```javascript
{
    "type": "mood",
    "day": 0-6, // day of the week, 0 (Sunday) to 6 (Saturday)
    "month": "MM", // month of the year, 0 (January) to 11 (December)
    "year": "yyyy", // current year in UTC time
    "mood": -2-2, // relative mood, -2 to 2
}
```

## Submitting a new comment for a mood (comment)

```javascript
{
    "type": "comment",
    "mood": -1-1, // integer, either negative (-1), neutral (0), or positive (1)
    "data": "" // string containing comment data
}
```

## Updating user details

```javascript
{
    "type": "details",
    "emailAddress": "", // new email address
    "password": "", // new password
    "data": "" // new name
}
```

Any values which remain unchanged **MUST** still be sent, but should be sent as an empty string (`""`), to indicate no change is needed.

The only details required for an account is the **email address and password**, the user's name is not required and will default to "friend".


## Sending a deletion request

```javascript
{
    "type": "delete"
}
```

This action, depending on how the server's setup, will usually result in instantaneous data deletion. This should remove the user file, and purge the in-memory user data (nullify the variable). **Report data will be unaffected by this request, accounting for the fact that report data may be critical to a law enforcement investigation.**

Handling of this request should begin immediately; finality warnings should be addressed by the client.

## Chat API: Start a new chat

This method must be called, and a ChatID must be provided, before any other chat API methods can be called.

```javascript
{
    "type": "chat:start"
}
```

Handling of setting up chat details, and connecting users together (usually via a first-come-first-serve queue system) is handled by the server. If the user is awaiting a partner, the server should send this as response:

```javascript
{
    "type": "chat:ready",
    "flag": false
}
```

Otherwise, if a chat has been initiated, this should be sent:

```javascript
{
    "type": "chat:ready",
    "flag": true,
    "cid": "" // chat ID goes here (see below)
}
```

The `cid` field should contain a string **unique chat identification**, which SHOULD NOT be an incremental counter. Treat this like a user ID; this will be used by the clients to address the chat which they are in to the server. If an invalid ChatID is sent, just ignore the request; it's malformed and likely a spoof attack.

**UPON CHAT TERMINATION BY EITHER PARTY, THE SERVER WILL SEND A `"chat:closed"` NOTICE TO THE OTHER PARTY.**

## Chat API: Send a message

```javascript
{
    "type": "chat:send",
    "data": "", // Message content goes here, in string format
    "cid": "" // the ChatID provided by the server goes here
}
```

This will be handled by the server and sent to the other party of the conversation. If the server has a provided Google Cloud API Key, then the message will first be sent through the Perspective API Neural Network. A toxicity probability above 0.9 will block the message (see below), whereas a lower probability **or an unsuccessful (non-200 return value) request** will allow the message to go straight through.

If a message is rejected, the server will send the client a JSON request of the following format:

```javascript
{
    "type": "chat:rejected",
    "mid": "" // incremental message ID goes here relating to the message's index position in the RAM log.
}
```

To authorise a message, see the `verify` Chat API method below.

Upon a successful message send, the sender will receive the following JSON message from the server:

```javascript
{
    "type": "chat:message",
    "flag": false,
    "data": "" // chat message
}
```

Whereas the recipient will receive:

```javascript
{
    "type": "chat:message",
    "flag": true,
    "data": "" // chat message
}
```

## Chat API: Verify a rejected message

Message rejections are not final, acknowledging the unreliability of AI. Thus, if a message is rejected, the client should ask the user. If the user insists, this method can be used to 'force' the message through.

```javascript
{
    "type": "chat:verify",
    "cid": "", // the chat id goes here
    "mid": "" // the message id goes here
}
```

Sending will then occur as if the message had been allowed through the filter automatically.

## Chat API: Report a chat

```javascript
{
    "type": "chat:report",
    "cid": "" // the chat id to report
}
```

Different servers may handle chat reports differently. This is not addressed by this specification.

## Admin API: Access a report

**ADMIN API METHODS WILL ONLY RETURN A RESPONSE IF THE LOGGED-IN EOS ACCOUNT'S "ADMIN" PARAMETER IS SET TO TRUE.**

```javascript
{
    "type": "admin:access",
    "cid": "" // ID of the chatlog to access (aka chat ID or report ID)
}
```

=> This will return the following response:
```javascript
{
    "type": "admin:chatlog",
    "chatlog": [
        {
            "aiDecision": true, // true if the message was delivered, false otherwise. Messages may have two objects here, one true and one false, if a rejected message was forced through. This aids moderation.
            "sender": "", // user ID of the message's sender.
            "message": "" // message content.
        }, ...
    ]
}
```

## Admin API: Decide a report

**ADMIN API METHODS WILL ONLY RETURN A RESPONSE IF THE LOGGED-IN EOS ACCOUNT'S "ADMIN" PARAMETER IS SET TO TRUE.**

```javascript
{
    "type": "admin:decision",
    "cid": "", // ID of the chatlog you are deciding on.
    "flag": true // true to accept the report as valid and ban the user; false to reject the report as invalid and trigger no action.
}
```

Upon a decision being made, the chatlog should therefore be deleted immediately. **It is important to note that, if a chatlog contains illegal content or is part of an ongoing investigation, the matter should be escalated to the server admin and a copy of all relevant data should be made by them. The report should NOT be decided upon.**