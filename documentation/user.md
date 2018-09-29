# Eos Server Response: User object

In response to a successful login, a filled User object is returned. If login is unsuccessful, or if the method does not return a User, then this object will have empty (default value) fields. Refer to the Server API documentation for more information.

## The Structure of a User Object

User objects are sent as a direct, JSON representation of the server-side User object from Go. Therefore, they follow a strict structure (as with all server responses.)

First, a skeletal structure will be shown. Each field will be described in more detail below.

```json
{
    "UserID" : [],
    "EmailAddr": "",
    "Password": [],
    "Name": "",
    "Moods": {},
    "Positives": [],
    "Neutrals": [],
    "Negatives": [],
    "Admin": false,
    "Banned": false
}
```

## User ID `UserID`

The UserID is a unique, server time-based identifier for each individual user, generated upon initial account creation. No two accounts can have the same UserID. The UserID can, internally, be either a string or a byte array; in responses, the UserID is provided as a byte array.

The UserID is provided for marginal use-cases only. It's highly improbable it will ever be needed.

## Email address `EmailAddr`

The `EmailAddr` property stores the current email address for the user that we have on file. When the user logs in, we perform a reverse-lookup of the email address to discover the account's UserID for subsequent operations. Two accounts cannot have the same email address.

This is sent as a string value.

## Password `Password`

The `Password` property represents the server-internal password property of the user structure, and **is always empty**. It is sent as a byte array containing all null values. This is because we do not see a use-case in which the user or client will require the stored hash, and to be even safer, we prefer to not disclose it.

## Name `Name`

Sent as a string, this is the name currently associated with the account. There is no unique behaviour associated with this; this should simply be used to communicate with the user. (For example, "Hey, `name`!" in your app.)

## Moods `Moods`

The Moods parameter is sent as a unique object, and is an exhaustive copy of all stored mood data within the Eos server. It follows the following structure:

```json
{
    "Day": [
        {
            "Mood": 0,
            "Num": 0
        },
        {
            "Mood": 0,
            "Num": 0
        },
        {
            "Mood": 0,
            "Num": 0
        },
        {
            "Mood": 0,
            "Num": 0
        },
        {
            "Mood": 0,
            "Num": 0
        },
        {
            "Mood": 0,
            "Num": 0
        },
        {
            "Mood": 0,
            "Num": 0
        }
    ],
    "Month": [
        {
            "Mood": 0,
            "Num": 0
        },
        {
            "Mood": 0,
            "Num": 0
        },
        {
            "Mood": 0,
            "Num": 0
        },
        {
            "Mood": 0,
            "Num": 0
        },
        {
            "Mood": 0,
            "Num": 0
        },
        {
            "Mood": 0,
            "Num": 0
        },
        {
            "Mood": 0,
            "Num": 0
        },
        {
            "Mood": 0,
            "Num": 0
        },
        {
            "Mood": 0,
            "Num": 0
        },
        {
            "Mood": 0,
            "Num": 0
        },
        {
            "Mood": 0,
            "Num": 0
        },
        {
            "Mood": 0,
            "Num": 0
        }
    ],
    "Years": [
        {
            "Year": 2018,
            "Month": [
                {
                    "Mood": 0,
                    "Num": 0
                },
                {
                    "Mood": 0,
                    "Num": 0
                },
                {
                    "Mood": 0,
                    "Num": 0
                },
                {
                    "Mood": 0,
                    "Num": 0
                },
                {
                    "Mood": 0,
                    "Num": 0
                },
                {
                    "Mood": 0,
                    "Num": 0
                },
                {
                    "Mood": 0,
                    "Num": 0
                },
                {
                    "Mood": 0,
                    "Num": 0
                },
                {
                    "Mood": 0,
                    "Num": 0
                },
                {
                    "Mood": 0,
                    "Num": 0
                },
                {
                    "Mood": 0,
                    "Num": 0
                },
                {
                    "Mood": 0,
                    "Num": 0
                }
            ]
        },
        {
            "Year": 2018,
            "Month": [
                {
                    "Mood": 0,
                    "Num": 0
                },
                {
                    "Mood": 0,
                    "Num": 0
                },
                {
                    "Mood": 0,
                    "Num": 0
                },
                {
                    "Mood": 0,
                    "Num": 0
                },
                {
                    "Mood": 0,
                    "Num": 0
                },
                {
                    "Mood": 0,
                    "Num": 0
                },
                {
                    "Mood": 0,
                    "Num": 0
                },
                {
                    "Mood": 0,
                    "Num": 0
                },
                {
                    "Mood": 0,
                    "Num": 0
                },
                {
                    "Mood": 0,
                    "Num": 0
                },
                {
                    "Mood": 0,
                    "Num": 0
                },
                {
                    "Mood": 0,
                    "Num": 0
                }
            ]
        }
    ]
}
```

Where:

- `Day` is an array of 7 `Mood` objects. Eos begins counting weeks from Sunday, and thus object 0 is Sunday, with object 6 being Saturday.
- `Month` is an array of 12 `Mood` objects, beginning at January [0] and ending at December [11].
- `Years` is an array of 2 `Year` objects, each containing a `Year` number (which will be the name of the current year - for example, 2018) and another array of 12 `Mood` objects specific to that year, beginning at January [0] and ending at December [11].

Each `Mood` object contains:

- A `Mood` property, which is a running total of the user's moods. For example, if a user selected -1, 2, and 0, then the `Mood` property would be (-1 + 2 + 0 = ) 1.
- A `Num` property, which is a running tally of the number of `Mood` changes (that is, number of moods recorded by the user). For our previous example, this would be 3.
- These are used to calculate the average. For example, 1/3 = ~0.33.

## Positive comments `Positives`

This is an array of 20 strings, containing the latest 20 positive comments the user has entered. Ideally, these would be displayed on a negative response page, to help lighten the user's mood. Oldest comments are "pushed out" when new comments are added. Unfilled comments will be an empty string (`""`).

## Neutral comments `Neutrals`

Provided for statistical or report purposes, this contains the latest 5 'neutral' comments (that is, comments entered with a relative mood of `0`). This behaves the same otherwise as `Positives`.

## Negative comments `Neutrals`

Provided for statistical or report purposes, this contains the latest 5 negative comments. This behaves the same otherwise as `Positives`.

## Administrator flag `Admin`

This is a boolean value, and is provided to allow clients to customise their UI as appropriate. This property will be `true` if the user is marked as an administrator account, or `false` otherwise. If the value is `false`, the account will not be able to perform administrator actions, so any UI elements making use of this should be hidden. The server will reject unauthorised methods silently.

## Ban flag `Banned`

This is a boolean value, provided to provide early notice that the user's account is banned from using the Eos Chat functionality. Clients may customise UI based on this. If ignored, the server will reject any `chat:start` requests by the user.