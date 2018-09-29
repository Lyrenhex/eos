/*
 * Eos Backend Server
 *
 * Copyright (c) Damian Heaton 2017 All rights reserved.
 *
 * Server Structures
 */

package main

import (
	"github.com/gorilla/websocket"
	"github.com/nu7hatch/gouuid"
)

// Configuration stores the JSON configuration stored in `config.json` as a Go-friendly structure.
type Configuration struct {
	EnvProd  bool   `json:"envProduction"`
	EnvKey   string `json:"envKey"`
	EnvCert  string `json:"envCertificate"`
	SrvHost  string `json:"srvHostname"`
	SrvPort  int    `json:"srvPort"`
	GApiKey  string `json:"googleApiKey"`
	DWebhook string `json:"discordWebhook"`
}

// Payload acts as a consistent structure to interface with JSON client-server exchange data.
type Payload struct {
	Type   string        `json:"type"`
	Flag   bool          `json:"flag"`
	Data   string        `json:"data"`
	Email  string        `json:"emailAddress"`
	Pass   string        `json:"password"`
	Day    int           `json:"day"`
	Month  int           `json:"month"`
	Year   int           `json:"year"`
	Mood   int           `json:"mood"`
	MsgID  int           `json:"mid"`
	ChatID string        `json:"cid"`
	User   User          `json:"user"`
	Log    []ChatMessage `json:"chatlog"`
}

// User data type, built on numerous structures.
type User struct {
	UserID    uuid.UUID
	EmailAddr string
	Password  []byte
	Name      string
	Moods     Moods
	Positives [20]string // allow more positives, but cap at 20 comments before replacing existing ones.
	Neutrals  [5]string  // less emphasis on storing non-positive comments. Keep 5 for reports before replacement.
	Negatives [5]string
	Admin     bool
	Banned    bool
}

// Moods acts as a reusable structure to store mood data - sub structure of User
type Moods struct {
	Day   [7]Mood // array of 7 moods, one for each day of week.
	Month [12]Mood
	Years [2]Year // only keep specific data on the past two years. we cannot overload the server. (not sure if we should decrease this to 1 year?)
}

// Mood stores information for a particular time unit
type Mood struct {
	Mood int
	Num  int
}

// Year structure to create a copy of a year\"s Moods structure
type Year struct {
	Year  int
	Month [12]Mood
}

// Data stores key information for chat sessions
type Data struct {
	UserID uuid.UUID
}

// WaitingUser stores a userID and pointer to the user's WebSocket connection for when the chat is instantiated
type WaitingUser struct {
	UserID     uuid.UUID
	Connection *websocket.Conn
}

// MLRequest is a template for generating Google Perspective API requests
type MLRequest struct {
	Comment         MLComment   `json:"comment"`
	RequestedAttrbs MLAttribute `json:"requestedAttributes"`
	DNS             bool        `json:"doNotStore"`
}

// MLComment : struct, part of MLRequest to serve chat message content
type MLComment struct {
	Text string `json:"text"`
}

// MLAttribute : struct, part of MLRequest to request SEVERE_TOXICITY model results
type MLAttribute struct {
	Attrb MLTOXICITY `json:"SEVERE_TOXICITY"`
}

// MLTOXICITY : struct, part of MLAttribute; empty
type MLTOXICITY struct{}

// MLResponse is a template for parsing Google Perspective API responses
type MLResponse struct {
	AttrbScores AttributeScores `json:"attributeScores"`
}

// AttributeScores : struct, part of MLResponse to parse the SEVERE_TOXICITY results
type AttributeScores struct {
	Toxicity Toxicity `json:"SEVERE_TOXICITY"`
}

// Toxicity : struct, part of Toxicity to parse the summaryScore of the SEVERE_TOXICITY results
type Toxicity struct {
	Summary SummaryScores `json:"summaryScore"`
}

// SummaryScores : struct to parse the summaryScore
type SummaryScores struct {
	Score float64 `json:"value"`
	Type  string  `json:"type"`
}

// DiscordWebhookRequest is a structure template to generate JSON objects for a Discord API request to send via webhook URI
type DiscordWebhookRequest struct {
	Content [1]DiscordWebhookEmbed `json:"embeds"`
}

// DiscordWebhookEmbed : struct, content of DiscordWebhookRequest
type DiscordWebhookEmbed struct {
	ReportID    string `json:"title"`
	Description string `json:"description"`
	ReportUri   string `json:"url"`
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
