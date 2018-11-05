package chat

import (
	"github.com/gorilla/websocket"
	uuid "github.com/nu7hatch/gouuid"
)

var UserPairs map[uuid.UUID]WaitingUser // @UserPairs[uuid.UUID], store a pointer to the other member\"s websocket connection
var QueuedUser WaitingUser              // store a WaitingUser object to represent the connection data and UserID of the currently waiting user for Eos chat.
var Chatlogs map[string]([]ChatMessage) // store a 2d array of ChatMessages - []ChatMessage per chat.

// WaitingUser stores a userID and pointer to the user's WebSocket connection for when the chat is instantiated
type WaitingUser struct {
	UserID     uuid.UUID
	Connection *websocket.Conn
}

// DiscordWebhookRequest is a structure template to generate JSON objects for a Discord API request to send via webhook URI
type DiscordWebhookRequest struct {
	Content [1]DiscordWebhookEmbed `json:"embeds"`
}

// DiscordWebhookEmbed : struct, content of DiscordWebhookRequest
type DiscordWebhookEmbed struct {
	ReportID    string `json:"title"`
	Description string `json:"description"`
	ReportURI   string `json:"url"`
}

// ChatLog stores an array of ChatMessages for simple log parsing
type ChatLog struct {
	ChatLog []ChatMessage `json:"chatlog"`
}

// ChatMessage stores information on: whether message was sent or blocked; the string userID of the sender; and the message content
type ChatMessage struct {
	Sent    bool   `json:"aiDecision"`
	User    string `json:"sender"`
	Message string `json:"message"`
}
