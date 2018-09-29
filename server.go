/*
 * Eos Backend Server
 *
 * Copyright (c) Damian Heaton 2017 All rights reserved.
 *
 * Server operates on port 9874 by default -- please see config.json
 */

package main

import (
	"bytes"
	"encoding/json"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/nu7hatch/gouuid"
	"golang.org/x/crypto/bcrypt"
)

// VERSION stores hardcoded constant storing the server version. AN X VERSION SERVER SHOULD DISTRIBUTE WEBAPP FILES COMPATIBLE WITH IT
const VERSION = "2.0:stage"

var config = loadConfig()

// UserIDs maps email address to user id
var UserIDs map[string]uuid.UUID        // @UserIDs["email address"] => uuid.UUID
var users map[uuid.UUID]*User           // @users[uuid.UUID], store a pointer to the user\"s data - avoid memory data duplication - edit the storage location directly.
var userPairs map[uuid.UUID]WaitingUser // @UserPairs[uuid.UUID], store a pointer to the other member\"s websocket connection
var wUser WaitingUser                   // store a WaitingUser object to represent the connection data and UserID of the currently waiting user for Eos chat.
var chatlogs map[string]([]ChatMessage) // store a 2d array of ChatMessages - []ChatMessage per chat.

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return (strings.HasPrefix(r.Header.Get("Origin"), "https://"+config.SrvHost) || strings.HasPrefix(r.Header.Get("Origin"), "http://"+config.SrvHost))
	},
}

func main() {
	users = make(map[uuid.UUID]*User)
	chatlogs = make(map[string][]ChatMessage)
	userPairs = make(map[uuid.UUID]WaitingUser)

	f, err := os.OpenFile("data/server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	log.SetOutput(f)

	readUserIDs()

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
		data := &Data{}
		alive := true
		conn.SetCloseHandler(func(code int, text string) error {
			log.Println("connection closed; breaking inf loop")
			if partner, activeConv := userPairs[data.UserID]; activeConv {
				partner.Connection.WriteJSON(&Payload{
					Type: "chat:closed",
				})
			} else if wUser.UserID == data.UserID {
				defaultWUser := WaitingUser{}
				wUser = defaultWUser
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
				uid, success := loginUser(payload.Email, payload.Pass)
				userData := User{}
				if success {
					data.UserID = uid
					userData = *users[uid]
				} else if uid == uuid.UUID([16]byte{}) {
					// user doesn\"t exist; make them
					uid = newUser(payload.Email, payload.Pass, "friend")
					success = true
					data.UserID = uid
					userData = *users[uid]
				} else {
					continue
				}
				userData.Password = []byte("")
				conn.WriteJSON(&Payload{
					Type: "login",
					Flag: success,
					User: userData,
				})
			case "mood":
				users[data.UserID].Moods.Day[payload.Day].Mood += payload.Mood
				users[data.UserID].Moods.Day[payload.Day].Num++
				users[data.UserID].Moods.Month[payload.Month].Mood += payload.Mood
				users[data.UserID].Moods.Month[payload.Month].Num++

				yearRecorded := false
				for i, year := range users[data.UserID].Moods.Years {
					if year.Year == payload.Year {
						users[data.UserID].Moods.Years[i].Month[payload.Month].Mood += payload.Mood
						users[data.UserID].Moods.Years[i].Month[payload.Month].Num++
						yearRecorded = true
					}
				}
				if !yearRecorded {
					newYear := Year{
						Year:  payload.Year,
						Month: [12]Mood{},
					}
					users[data.UserID].Moods.Years = [2]Year{
						users[data.UserID].Moods.Years[1],
						newYear,
					}
					users[data.UserID].Moods.Years[1].Month[payload.Month].Mood += payload.Mood
					users[data.UserID].Moods.Years[1].Month[payload.Month].Num++
				}
				saveUser(data.UserID)
			case "comment":
				log.Println(payload)
				mood := payload.Mood
				comment := payload.Data
				if comment != "" {
					switch mood {
					case 1:
						users[data.UserID].Positives = [20]string{
							users[data.UserID].Positives[1],
							users[data.UserID].Positives[2],
							users[data.UserID].Positives[3],
							users[data.UserID].Positives[4],
							users[data.UserID].Positives[5],
							users[data.UserID].Positives[6],
							users[data.UserID].Positives[7],
							users[data.UserID].Positives[8],
							users[data.UserID].Positives[9],
							users[data.UserID].Positives[10],
							users[data.UserID].Positives[11],
							users[data.UserID].Positives[12],
							users[data.UserID].Positives[13],
							users[data.UserID].Positives[14],
							users[data.UserID].Positives[15],
							users[data.UserID].Positives[16],
							users[data.UserID].Positives[17],
							users[data.UserID].Positives[18],
							users[data.UserID].Positives[19],
							comment,
						}
					case 0:
						users[data.UserID].Neutrals = [5]string{
							users[data.UserID].Neutrals[1],
							users[data.UserID].Neutrals[2],
							users[data.UserID].Neutrals[3],
							users[data.UserID].Neutrals[4],
							comment,
						}
					case -1:
						users[data.UserID].Negatives = [5]string{
							users[data.UserID].Negatives[1],
							users[data.UserID].Negatives[2],
							users[data.UserID].Negatives[3],
							users[data.UserID].Negatives[4],
							comment,
						}
					}
					saveUser(data.UserID)
				}
			case "details":
				newEmail := payload.Email
				newPass := payload.Pass
				newName := payload.Data
				if newEmail != "" {
					users[data.UserID].EmailAddr = newEmail
				}
				if newPass != "" {
					newPass, _ := bcrypt.GenerateFromPassword([]byte(payload.Pass), bcrypt.DefaultCost)
					users[data.UserID].Password = newPass
				}
				if newName != "" {
					users[data.UserID].Name = newName
				}
				saveUser(data.UserID)
			case "delete":
				delete(UserIDs, users[data.UserID].EmailAddr)
				err := os.Remove("data/userdata-" + data.UserID.String() + ".json")
				if err != nil {
					log.Println("Error deleting userdata-"+data.UserID.String()+".json: ", err)
				}
				saveUserIDs()
			case "chat:start":
				if !users[data.UserID].Banned {
					userWUser := WaitingUser{
						UserID:     data.UserID,
						Connection: conn,
					}
					defaultWUser := WaitingUser{}
					if wUser != defaultWUser {
						// generate new chat ID
						cid, _ := uuid.NewV4()
						strCid := cid.String()
						log.Println(chatlogs)
						chatlogs[strCid] = make([]ChatMessage, 0)
						log.Println(chatlogs)

						userPairs[data.UserID] = wUser
						userPairs[wUser.UserID] = userWUser

						conn.WriteJSON(&Payload{
							Type:   "chat:ready",
							Flag:   true,
							ChatID: strCid,
						})
						wUser.Connection.WriteJSON(&Payload{
							Type:   "chat:ready",
							Flag:   true,
							ChatID: strCid,
						})
						wUser = defaultWUser
					} else {
						wUser = userWUser
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

					request := &MLRequest{
						Comment:         MLComment{Text: payload.Data},
						RequestedAttrbs: MLAttribute{MLTOXICITY{}},
						DNS:             true,
					}
					jsonValue, _ := json.Marshal(request)
					resp, _ := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonValue))
					body, err := ioutil.ReadAll(resp.Body)
					response := *resp
					if response.StatusCode == 200 && err == nil { // request went through - huzzah
						mlResponse := MLResponse{}
						json.Unmarshal(body, &mlResponse)
						response.Body.Close()
						sendMessage := true
						if mlResponse.AttrbScores.Toxicity.Summary.Score >= 0.9 {
							// reject the message
							conn.WriteJSON(&Payload{
								Type:  "chat:rejected",
								MsgID: len(chatlogs[payload.ChatID]),
							})
							sendMessage = false
						}
						chatlogs[payload.ChatID] = append(chatlogs[payload.ChatID], ChatMessage{
							Sent:    sendMessage,
							User:    data.UserID.String(),
							Message: html.EscapeString(payload.Data),
						})

						conn.WriteJSON(&Payload{
							Type: "chat:message",
							Flag: false,
							Data: html.EscapeString(payload.Data),
						})
						userPairs[data.UserID].Connection.WriteJSON(&Payload{
							Type: "chat:message",
							Flag: true,
							Data: html.EscapeString(payload.Data),
						})
					} else {
						log.Println("Error occurred in NEURAL NETWORK: ", response.StatusCode, err)
						log.Println("Bypassing filter, sending message.")

						chatlogs[payload.ChatID] = append(chatlogs[payload.ChatID], ChatMessage{
							Sent:    true,
							User:    data.UserID.String(),
							Message: html.EscapeString(payload.Data),
						})

						conn.WriteJSON(&Payload{
							Type: "chat:message",
							Flag: false,
							Data: html.EscapeString(payload.Data),
						})
						userPairs[data.UserID].Connection.WriteJSON(&Payload{
							Type: "chat:message",
							Flag: true,
							Data: html.EscapeString(payload.Data),
						})
					}
				}
			case "chat:verify":
				msg := chatlogs[payload.ChatID][payload.MsgID]
				if msg.User == data.UserID.String() {
					conn.WriteJSON(&Payload{
						Type: "chat:message",
						Flag: false,
						Data: html.EscapeString(payload.Data),
					})
					userPairs[data.UserID].Connection.WriteJSON(&Payload{
						Type: "chat:message",
						Flag: true,
						Data: html.EscapeString(payload.Data),
					})
					chatlogs[payload.ChatID] = append(chatlogs[payload.ChatID], ChatMessage{
						Sent:    true,
						User:    data.UserID.String(),
						Message: html.EscapeString(payload.Data),
					})
				}
			case "chat:report":
				file, err := os.Create("data/reportlog-" + payload.ChatID + ".json")
				if err == nil {
					encoder := json.NewEncoder(file)
					err = encoder.Encode(ChatLog{
						ChatLog: chatlogs[payload.ChatID],
					})
					if err != nil {
						log.Println("Error saving reportlog-"+payload.ChatID+".json: ", err)
					}
					file.Close()

					request := &DiscordWebhookRequest{
						Content: [1]DiscordWebhookEmbed{DiscordWebhookEmbed{
							ReportID:    payload.ChatID,
							Description: "New reported chat log. Please click the link to access the page with which to handle this report log. This link will expire after the report has been addressed, and requires a valid administrator login. In cases where the chat log includes illegal content, please refer to Lyrenhex for escalation and referral to the local law enforcement authorities.",
							ReportUri:   "http://" + config.SrvHost + "/admin.html?id=" + payload.ChatID,
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
			case "admin:access":
				if users[data.UserID].Admin {
					file, err := os.Open("data/reportlog-" + payload.ChatID + ".json")
					if err != nil {
						if os.IsNotExist(err) {
							log.Println("Request to access nonexistent report " + payload.ChatID + ".")
						} else {
							log.Println("error:", err)
						}
					} else {
						decoder := json.NewDecoder(file)
						reportlog := ChatLog{}
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
				if users[data.UserID].Admin {
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
							users[*bannedID].Banned = true
							saveUser(*bannedID)
							log.Println("User " + bannedID.String() + " banned.")
						}
					}
					conn.WriteJSON(&Payload{
						Type: "admin:success",
					})
				}
			case "admin:flag":
				if users[data.UserID].Admin {
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
