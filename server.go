/*
 * Eos Backend Server
 *
 * Copyright (c) Damian Heaton 2017 All rights reserved.
 *
 * Server operates on port 9874 by default -- please see config.json
 */

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/go-homedir"
	uuid "github.com/nu7hatch/gouuid"
	"gitlab.com/lyrenhex/eos-v2/chat"
	"gitlab.com/lyrenhex/eos-v2/mail"
	"gitlab.com/lyrenhex/eos-v2/perspectiveapi"
	"gitlab.com/lyrenhex/eos-v2/user"
	"golang.org/x/crypto/bcrypt"
)

// VERSION stores hardcoded constant storing the server version. AN X VERSION SERVER SHOULD DISTRIBUTE WEBAPP FILES COMPATIBLE WITH IT
const VERSION = "2.0:live"

// Configuration stores the JSON configuration stored in `config.json` as a Go-friendly structure.
type Configuration struct {
	EnvProd      bool   `json:"envProduction"`
	EnvKey       string `json:"envKey"`
	EnvCert      string `json:"envCertificate"`
	SrvHost      string `json:"srvHostname"`
	SrvPort      int    `json:"srvPort"`
	GApiKey      string `json:"googleApiKey"`
	DWebhook     string `json:"discordWebhook"`
	MailAPIKey   string `json:"sendgridApiKey"`
	MailAPIAuth  string `json:"sendgridApiAuth"`
	MailAPIReset string `json:"sendgridApiReset"`
	MailAddress  string `json:"sendgridAddress"`
}

func (c *Configuration) load() {
	/* Load the server configuration from ~/eos/data/config.json, and return a Configuration structure pre-populated with the config data. If config.json is not openable, throw a Fatal error to terminate (we cannot recover; config is necessary) */
	home, err := homedir.Dir()
	if err != nil {
		home = "."
	}
	file, err := os.Open(home + "/eos/data/config.json")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Expected configuration values in `config.json`, got:")
			fmt.Printf("%+v\n", c)
			log.Fatal("Data folder or config.json does not exist. Please create the data folder and populate the config.json file before run.")
		} else {
			log.Println("error:", err)
		}
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		log.Fatal("Error reading config.json: ", err)
	}
}

// input accepts a string prompt and optionally a default value, and will pose this as an input and return user response as a string.
// If `d` is an empty string, then it will be treated as having no default.
// If `r` is true, then the input is required and will be repeated until successful.
func input(p, d string, r bool) string {
	prompt := p + ": "
	if d != "" {
		prompt += "[" + d + "] "
	}
	var result string
	for {
		fmt.Print(prompt)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		result = scanner.Text()
		if len(result) != 0 {
			return result
		} else if d != "" {
			return d
		} else if !r {
			return ""
		}
		fmt.Println("No input supplied and no default is provided. Please supply an input.")
	}
}

func (c *Configuration) setup() {
	fmt.Println("** CONFIGURATION FILE EMPTY OR NOT FOUND **")
}

// Payload acts as a consistent structure to interface with JSON client-server exchange data.
type Payload struct {
	Type   string             `json:"type"`
	Flag   bool               `json:"flag"`
	Data   string             `json:"data"`
	Email  string             `json:"emailAddress"`
	Pass   string             `json:"password"`
	Day    int                `json:"day"`
	Month  int                `json:"month"`
	Year   int                `json:"year"`
	Mood   int                `json:"mood"`
	MsgID  int                `json:"mid"`
	ChatID string             `json:"cid"`
	User   user.User          `json:"user"`
	Log    []chat.ChatMessage `json:"chatlog"`
}

func redirectToHTTPS(w http.ResponseWriter, req *http.Request) {
	/* For production environments, forcefully redirect the user to HTTPS using HSTS */
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	log.Printf("client redirect to: %s", target)
	req.Header.Set("Strict-Transport-Security", "max-age=63072000")
	http.Redirect(w, req, target,
		http.StatusTemporaryRedirect)
}

var config = Configuration{}
var mailService = mail.SendGrid{}

func init() {
	config.load()

	mailService.APIKey = config.MailAPIKey
	mailService.APIAuth = config.MailAPIAuth
	mailService.APIReset = config.MailAPIReset
	mailService.Email = config.MailAddress

	user.Users = make(map[uuid.UUID]*user.User)
	user.PendingUsers = make(map[string]string)
	user.ResetKeys = make(map[string]string)
	chat.Chatlogs = make(map[string][]chat.ChatMessage)
	chat.UserPairs = make(map[uuid.UUID]chat.WaitingUser)

	user.ReadIDs()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return (strings.HasPrefix(r.Header.Get("Origin"), "https://"+config.SrvHost) || strings.HasPrefix(r.Header.Get("Origin"), "http://"+config.SrvHost))
	},
}

func main() {
	f, err := os.OpenFile("data/server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	log.SetOutput(f)

	// Concurrently run a simple static webserver on port 80 or port 443 if in Production environment, for serving the online webapp from the `webclient` directory.
	go func() {
		if config.EnvProd {
			log.Println("Running redirectToHTTPS server on port 80 and TLS FS on port 443")
			go http.ListenAndServe(":80", http.HandlerFunc(redirectToHTTPS))
			panic(http.ListenAndServeTLS(":443", config.EnvCert, config.EnvKey, http.FileServer(http.Dir("webclient"))))
		} else {
			log.Println("Running FS on port 80")
			panic(http.ListenAndServe(":80", http.FileServer(http.Dir("webclient"))))
		}
	}()

	// Run the main websocket server on the chosen port.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("CONNECT: WS")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		// connection established! inform the client of the server version
		conn.WriteJSON(&Payload{
			Type: "version",
			Data: VERSION,
		})
		u := &user.User{}
		alive := true
		conn.SetCloseHandler(func(code int, text string) error {
			log.Println("connection closed; breaking inf loop")
			if partner, activeConv := chat.UserPairs[u.UserID]; activeConv {
				partner.Connection.WriteJSON(&Payload{
					Type: "chat:closed",
				})
				defaultWUser := chat.WaitingUser{}
				chat.UserPairs[u.UserID] = defaultWUser
			} else if chat.QueuedUser.UserID == u.UserID {
				defaultWUser := chat.WaitingUser{}
				chat.QueuedUser = defaultWUser
			}
			alive = false
			return nil
		})
		for alive {
			// Read data from client
			payload := &Payload{}
			err = conn.ReadJSON(payload)
			if err != nil {
				log.Println(err)
			}
			switch payload.Type {
			case "login":
				payload.Email = strings.ToLower(payload.Email)
				success, _ := u.Login(payload.Email, payload.Pass)
				conn.WriteJSON(&Payload{
					Type: "login",
					Flag: success,
					User: *u,
				})
			case "resetPassword":
				payload.Email = strings.ToLower(payload.Email)
				success, authToken := user.ResetPassword(payload.Email)
				if success {
					mailService.SendToken(payload.Email, authToken)
					conn.WriteJSON(&Payload{
						Type: "resetPassword",
						Flag: true,
					})
				} else {
					conn.WriteJSON(&Payload{
						Type: "resetPassword",
						Flag: false,
					})
				}
			case "signup":
				payload.Email = strings.ToLower(payload.Email)
				_, exists := u.Login(payload.Email, "")
				if !exists {
					newID, _ := uuid.NewV4()
					emailID := newID.String()
					user.PendingUsers[emailID] = payload.Email
					mailService.SendAuth(payload.Email, emailID)
					conn.WriteJSON(&Payload{
						Type: "signup",
						Flag: true,
					})
				} else {
					conn.WriteJSON(&Payload{
						Type: "signup",
						Flag: false,
					})
				}
			case "verifyEmail":
				emailID := payload.Data
				conn.WriteJSON(&Payload{
					Type: "verifyEmail",
					Data: user.PendingUsers[emailID],
				})
			case "createAccount":
				emailID := payload.Data
				if user.UserIDs[user.PendingUsers[emailID]] == uuid.UUID([16]byte{}) {
					u = user.New(user.PendingUsers[emailID], payload.Pass, "friend")
					delete(user.PendingUsers, emailID)
					conn.WriteJSON(&Payload{
						Type: "login",
						Flag: true,
						User: *u,
					})
				} else {
					conn.WriteJSON(&Payload{
						Type: "login",
						Flag: false,
					})
				}
			case "mood":
				u.AddMood(payload.Day, payload.Month, payload.Year, payload.Mood)
			case "comment":
				u.AddComment(payload.Mood, payload.Data)
			case "details":
				newEmail := payload.Email
				newPass := payload.Pass
				newName := payload.Data
				if newEmail != "" {
					newID, _ := uuid.NewV4()
					emailID := newID.String()
					user.PendingUsers[emailID] = payload.Email
					mailService.SendAuth(newEmail, emailID)
					conn.WriteJSON(&Payload{
						Type: "changeEmailVerification",
					})
				}
				if newPass != "" {
					newPass, _ := bcrypt.GenerateFromPassword([]byte(payload.Pass), bcrypt.DefaultCost)
					u.Password = newPass
				}
				if newName != "" {
					u.Name = newName
				}
				u.Save()
			case "changeEmail":
				emailID := payload.Data
				newEmail := user.PendingUsers[emailID]
				if newEmail != "" {
					u.EmailAddr = newEmail
					conn.WriteJSON(&Payload{
						Type:  "changeEmail",
						Email: newEmail,
					})
				}
			case "delete":
				delete(user.UserIDs, u.EmailAddr)
				err := os.Remove("data/userdata-" + u.UserID.String() + ".json")
				if err != nil {
					log.Println("Error deleting userdata-"+u.UserID.String()+".json: ", err)
				}
				user.SaveIDs()
			case "chat:start":
				if !u.Banned {
					userWUser := chat.WaitingUser{
						UserID:     u.UserID,
						Connection: conn,
					}
					defaultWUser := chat.WaitingUser{}
					if chat.QueuedUser != defaultWUser {
						// generate new chat ID
						cid, _ := uuid.NewV4()
						strCid := cid.String()
						log.Println(chat.Chatlogs)
						chat.Chatlogs[strCid] = make([]chat.ChatMessage, 0)
						log.Println(chat.Chatlogs)

						chat.UserPairs[u.UserID] = chat.QueuedUser
						chat.UserPairs[chat.QueuedUser.UserID] = userWUser

						conn.WriteJSON(&Payload{
							Type:   "chat:ready",
							Flag:   true,
							ChatID: strCid,
						})
						chat.QueuedUser.Connection.WriteJSON(&Payload{
							Type:   "chat:ready",
							Flag:   true,
							ChatID: strCid,
						})
						chat.QueuedUser = defaultWUser
					} else {
						chat.QueuedUser = userWUser
						conn.WriteJSON(&Payload{
							Type: "chat:ready",
							Flag: false,
						})
					}
				} else {
					conn.WriteJSON(&Payload{
						Type: "chat:banned",
					})
				}
			case "chat:send":
				if payload.Data != "" {
					apiURL := "https://commentanalyzer.googleapis.com/v1alpha1/comments:analyze?key=" + config.GApiKey

					request := &perspectiveapi.MLRequest{
						Comment:         perspectiveapi.MLComment{Text: payload.Data},
						RequestedAttrbs: perspectiveapi.MLAttribute{Attrb: perspectiveapi.MLTOXICITY{}},
						DNS:             true,
					}
					jsonValue, _ := json.Marshal(request)
					resp, _ := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonValue))
					body, err := ioutil.ReadAll(resp.Body)
					response := *resp
					if response.StatusCode == 200 && err == nil { // request went through - huzzah
						mlResponse := perspectiveapi.MLResponse{}
						json.Unmarshal(body, &mlResponse)
						response.Body.Close()
						sendMessage := true
						if mlResponse.AttrbScores.Toxicity.Summary.Score >= 0.9 {
							// reject the message
							conn.WriteJSON(&Payload{
								Type:  "chat:rejected",
								MsgID: len(chat.Chatlogs[payload.ChatID]),
							})
							sendMessage = false
						}
						chat.Chatlogs[payload.ChatID] = append(chat.Chatlogs[payload.ChatID], chat.ChatMessage{
							Sent:    sendMessage,
							User:    u.UserID.String(),
							Message: html.EscapeString(payload.Data),
						})

						conn.WriteJSON(&Payload{
							Type: "chat:message",
							Flag: false,
							Data: html.EscapeString(payload.Data),
						})
						chat.UserPairs[u.UserID].Connection.WriteJSON(&Payload{
							Type: "chat:message",
							Flag: true,
							Data: html.EscapeString(payload.Data),
						})
					} else {
						log.Println("Error occurred in NEURAL NETWORK: ", response.StatusCode, err)
						log.Println("Bypassing filter, sending message.")

						chat.Chatlogs[payload.ChatID] = append(chat.Chatlogs[payload.ChatID], chat.ChatMessage{
							Sent:    true,
							User:    u.UserID.String(),
							Message: html.EscapeString(payload.Data),
						})

						conn.WriteJSON(&Payload{
							Type: "chat:message",
							Flag: false,
							Data: html.EscapeString(payload.Data),
						})
						chat.UserPairs[u.UserID].Connection.WriteJSON(&Payload{
							Type: "chat:message",
							Flag: true,
							Data: html.EscapeString(payload.Data),
						})
					}
				}
			case "chat:verify":
				msg := chat.Chatlogs[payload.ChatID][payload.MsgID]
				if msg.User == u.UserID.String() {
					conn.WriteJSON(&Payload{
						Type: "chat:message",
						Flag: false,
						Data: html.EscapeString(payload.Data),
					})
					chat.UserPairs[u.UserID].Connection.WriteJSON(&Payload{
						Type: "chat:message",
						Flag: true,
						Data: html.EscapeString(payload.Data),
					})
					chat.Chatlogs[payload.ChatID] = append(chat.Chatlogs[payload.ChatID], chat.ChatMessage{
						Sent:    true,
						User:    u.UserID.String(),
						Message: html.EscapeString(payload.Data),
					})
				}
			case "chat:report":
				file, err := os.Create("data/reportlog-" + payload.ChatID + ".json")
				if err == nil {
					encoder := json.NewEncoder(file)
					err = encoder.Encode(chat.ChatLog{
						ChatLog: chat.Chatlogs[payload.ChatID],
					})
					if err != nil {
						log.Println("Error saving reportlog-"+payload.ChatID+".json: ", err)
					}
					file.Close()

					request := &chat.DiscordWebhookRequest{
						Content: [1]chat.DiscordWebhookEmbed{chat.DiscordWebhookEmbed{
							ReportID:    payload.ChatID,
							Description: "New reported chat log. Please click the link to access the page with which to handle this report log. This link will expire after the report has been addressed, and requires a valid administrator login. In cases where the chat log includes illegal content, please refer to Lyrenhex for escalation and referral to the local law enforcement authorities.",
							ReportURI:   "https://" + config.SrvHost + "/app/admin.html?id=" + payload.ChatID,
						}},
					}
					jsonValue, _ := json.Marshal(request)
					log.Println(string(jsonValue))
					resp, err := http.Post(config.DWebhook, "application/json", bytes.NewBuffer(jsonValue))
					if resp.StatusCode != 204 || err != nil {
						log.Println("Discord Webhook error: ", resp, err)
					}
				} else {
					log.Println("Error saving reportlog-"+payload.ChatID+".json: ", err)
				}
			case "chat:close":
				if partner, activeConv := chat.UserPairs[u.UserID]; activeConv {
					defaultWUser := chat.WaitingUser{}
					partner.Connection.WriteJSON(&Payload{
						Type: "chat:closed",
					})
					conn.WriteJSON(&Payload{
						Type: "chat:closed",
					})
					chat.UserPairs[u.UserID] = defaultWUser
				}
			case "admin:access":
				if u.Admin {
					file, err := os.Open("data/reportlog-" + payload.ChatID + ".json")
					if err != nil {
						if os.IsNotExist(err) {
							log.Println("Request to access nonexistent report " + payload.ChatID + ".")
						} else {
							log.Println("error:", err)
						}
					} else {
						decoder := json.NewDecoder(file)
						reportlog := chat.ChatLog{}
						err = decoder.Decode(&reportlog)
						file.Close()
						if err != nil {
							log.Fatal("Error reading reportlog-"+payload.ChatID+".json: ", err)
						}
						log.Println(reportlog)
						conn.WriteJSON(&Payload{
							Type: "admin:chatlog",
							Log:  reportlog.ChatLog,
						})
					}
				}
			case "admin:decision":
				if u.Admin {
					err := os.Remove("data/reportlog-" + payload.ChatID + ".json")
					if err != nil {
						if os.IsNotExist(err) {
							log.Println("Request to decide nonexistent report " + payload.ChatID + ".")
						} else {
							log.Println("error:", err)
						}
					} else {
						log.Println("Decision rendered on Report " + payload.ChatID)
						if payload.Data != "" {
							bannedID, _ := uuid.ParseHex(payload.Data)
							user.Users[*bannedID].Banned = true
							user.Users[*bannedID].Save()
							log.Println("User " + bannedID.String() + " banned.")
						}
					}
					conn.WriteJSON(&Payload{
						Type: "admin:success",
					})
				}
			case "admin:flag":
				if u.Admin {
					err := os.Rename("data/reportlog-"+payload.ChatID+".json", "data/FLAGGED-reportlog-"+payload.ChatID+".json")
					if err != nil {
						if os.IsNotExist(err) {
							log.Println("Request to flag nonexistent report " + payload.ChatID + ".")
						} else {
							log.Println("error:", err)
						}
					} else {
						log.Println("Report " + payload.ChatID + " flagged.")
					}
				}
				conn.WriteJSON(&Payload{
					Type: "admin:success",
				})
			}
		}
	})
	if config.EnvProd {
		log.Println("Running TLS WS on port ", config.SrvPort)
		panic(http.ListenAndServeTLS(":"+strconv.Itoa(config.SrvPort), config.EnvCert, config.EnvKey, nil))
	} else {
		log.Println("Running WS on port", config.SrvPort)
		panic(http.ListenAndServe(":"+strconv.Itoa(config.SrvPort), nil))
	}
}
