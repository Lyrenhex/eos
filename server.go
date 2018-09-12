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
	"fmt"
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
const VERSION = "2.0:dev"

var config = loadConfig()

// UserIDs maps email address to user id
var UserIDs map[string]uuid.UUID
var users map[uuid.UUID]*User           // @users[uuid.UUID], store a pointer to the user\"s data - avoid memory data duplication - edit the storage location directly.
var userPairs map[uuid.UUID]WaitingUser // @UserPairs[uuid.UUID], store a pointer to the other member\"s websocket connection
var wUser WaitingUser
var chatlogs map[string]([]ChatMessage)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return (strings.HasPrefix(r.Header.Get("Origin"), "https://"+config.SrvHost) || strings.HasPrefix(r.Header.Get("Origin"), "http://"+config.SrvHost))
	},
}

// Configuration stores the JSON configuration stored in `config.json` as a Go-friendly structure.
type Configuration struct {
	EnvProd bool   `json:"envProduction"`
	EnvKey  string `json:"envKey"`
	EnvCert string `json:"envCertificate"`
	SrvHost string `json:"srvHostname"`
	SrvPort int    `json:"srvPort"`
	ApiKey  string `json:"apiKey"`
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

// Payload acts as a consistent structure to interface with JSON client-server exchange data.
type Payload struct {
	Type   string `json:"type"`
	Flag   bool   `json:"flag"`
	Data   string `json:"data"`
	Email  string `json:"emailAddress"`
	Pass   string `json:"password"`
	Day    int    `json:"day"`
	Month  int    `json:"month"`
	Year   int    `json:"year"`
	Mood   int    `json:"mood"`
	MsgID  int    `json:"mid"`
	ChatID string `json:"cid"`
	User   User   `json:"user"`
}

// Data stores key information for chat sessions
type Data struct {
	UserID uuid.UUID
}

type WaitingUser struct {
	UserID     uuid.UUID
	Connection *websocket.Conn
}

type MLRequest struct {
	Comment         MLComment   `json:"comment"`
	RequestedAttrbs MLAttribute `json:"requestedAttributes"`
	DNS             bool        `json:"doNotStore"`
}
type MLComment struct {
	Text string `json:"text"`
}
type MLAttribute struct {
	Attrb MLTOXICITY `json:"SEVERE_TOXICITY"`
}
type MLTOXICITY struct{}

type MLResponse struct {
	AttrbScores AttributeScores `json:"attributeScores"`
}
type AttributeScores struct {
	Toxicity Toxicity `json:"SEVERE_TOXICITY"`
}
type Toxicity struct {
	Summary SummaryScores `json:"summaryScore"`
}
type SummaryScores struct {
	Score float64 `json:"value"`
	Type  string  `json:"type"`
}

type ChatMessage struct {
	Sent    bool
	User    uuid.UUID
	Message string
}

func loadConfig() Configuration {
	/* Load the server configuration from data/config.json, and return a Configuration structure pre-populated with the config data. If config.json is not openable, throw a Fatal error to terminate (we cannot recover; config is necessary) */
	file, err := os.Open("data/config.json")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Expected configuration values in `config.json`, got:")
			fmt.Printf("%+v\n", Configuration{})
			log.Fatal("Data folder or config.json does not exist. Please create the data folder and populate the config.json file before run.")
		} else {
			log.Println("error:", err)
		}
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal("Error reading config.json: ", err)
	}
	return config
}

func newUser(email, pass, name string) uuid.UUID {
	/* For use during new user creation, generate a new userid and a new User structure to store the new user\"s data, prepopulated with email pass and name data, with the rest zero-ed for later input. Return the userid for immediate usage, and add to the *memory* user store, and maintain the email-userid pairing in the UserIDs map. */
	uid, _ := uuid.NewV4()
	password, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	user := &User{
		UserID:    *uid, //
		EmailAddr: email,
		Password:  password,
		Name:      name,
		Moods:     Moods{},
		Positives: [20]string{},
		Neutrals:  [5]string{},
		Negatives: [5]string{},
	}
	users[*uid] = user
	UserIDs[email] = *uid
	go saveUserIDs()
	go saveUser(*uid)
	return *uid
}
func saveUser(uid uuid.UUID) {
	/* Save the user\"s data to their own file, stored according to their user id. */
	file, err := os.Create("data/userdata-" + uid.String() + ".json")
	if err == nil {
		defer file.Close()
		encoder := json.NewEncoder(file)
		err = encoder.Encode(users[uid])
		if err != nil {
			log.Println("Error saving userdata-"+uid.String()+".json: ", err)
		}
	} else {
		log.Println("Error saving userdata-"+uid.String()+".json: ", err)
	}
}
func loadUser(uid uuid.UUID) {
	/* Based on the provided user id, load the user\"s data from their own data file, add the user back to the standard users map, and return the user structure for immediate use. */
	file, err := os.Open("data/userdata-" + uid.String() + ".json")
	if err == nil {
		defer file.Close()
		user := User{}
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&user)
		if err != nil {
			log.Println("Error reading userdata-"+uid.String()+".json: ", err)
		}
		users[uid] = &user
	} else {
		log.Println("Error reading userdata-"+uid.String()+".json: ", err)
	}
}
func loginUser(email, password string) (uuid.UUID, bool) {
	/* Based on provided email and password strings, check that the login is correct and, if so, load the user\"s data and return UID to calling func. */
	success := true
	uid := UserIDs[email]
	if uid == uuid.UUID([16]byte{}) { // `uid` matches the `uuid.UUID` default value
		success = false
	} else {
		loadUser(uid)
		err := bcrypt.CompareHashAndPassword(users[uid].Password, []byte(password))
		if err != nil {
			success = false
		}
	}

	return uid, success
}

func readUserIDs() {
	/* Read the email-userid pairs from users.json and store in the UserIDs map for reference when logging in (clients will send an email and password; lookup email in UserIDs to grab their userid, then load their userfile.json) */
	log.Println("Loading UserIDs from file")
	file, err := os.Open("data/users.json")
	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&UserIDs)
		if err != nil {
			log.Println("Error reading users.json: ", err)
		}
	} else {
		UserIDs = make(map[string]uuid.UUID)
	}
}
func saveUserIDs() {
	/* Save the email-userid pairs, replacing the existing file is present. */
	log.Println("Saving UserIDs to file")
	file, err := os.Create("data/users.json")
	if err == nil {
		defer file.Close()
		encoder := json.NewEncoder(file)
		err = encoder.Encode(UserIDs)
		if err != nil {
			log.Println("Error saving users.json: ", err)
		}
	} else {
		UserIDs = make(map[string]uuid.UUID)
	}
}

func redirectToHTTPS(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	log.Printf("client redirect to: %s", target)
	req.Header.Set("Strict-Transport-Security", "max-age=63072000")
	http.Redirect(w, req, target,
		// see @andreiavrammsd comment: often 307 > 301
		http.StatusTemporaryRedirect)
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
				file, err := os.Create("data/chatlogs.json")
				if err == nil {
					defer file.Close()
					encoder := json.NewEncoder(file)
					err = encoder.Encode(chatlogs)
					if err != nil {
						log.Println("Error saving chatlogs.json: ", err)
					}
				} else {
					log.Println("Error saving chatlogs.json: ", err)
				}
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
				if success {
					data.UserID = uid
				} else if uid == uuid.UUID([16]byte{}) {
					// user doesn\"t exist; make them
					uid = newUser(payload.Email, payload.Pass, "friend")
					success = true
					data.UserID = uid
				} else {
					continue
				}
				userData := *users[uid]
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
			case "chat:send":
				apiURL := "https://commentanalyzer.googleapis.com/v1alpha1/comments:analyze?key=" + config.ApiKey

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
					defer response.Body.Close()
					mlResponse := MLResponse{}
					json.Unmarshal(body, &mlResponse)
					sendMessage := true
					if mlResponse.AttrbScores.Toxicity.Summary.Score >= 0.9 {
						// reject the message
						conn.WriteJSON(&Payload{
							Type:  "chat:rejected",
							MsgID: len(chatlogs[payload.ChatID]),
						})
						sendMessage = false
					}
					log.Println(chatlogs)
					chatlogs[payload.ChatID] = append(chatlogs[payload.ChatID], ChatMessage{
						Sent:    sendMessage,
						User:    data.UserID,
						Message: payload.Data,
					})
					log.Println(chatlogs)

					conn.WriteJSON(&Payload{
						Type: "chat:message",
						Flag: false,
						Data: payload.Data,
					})
					userPairs[data.UserID].Connection.WriteJSON(&Payload{
						Type: "chat:message",
						Flag: true,
						Data: payload.Data,
					})
				}
			case "chat:verify":
				msg := chatlogs[payload.ChatID][payload.MsgID]
				if msg.User == data.UserID {
					conn.WriteJSON(&Payload{
						Type: "chat:message",
						Flag: false,
						Data: msg.Message,
					})
					userPairs[data.UserID].Connection.WriteJSON(&Payload{
						Type: "chat:message",
						Flag: true,
						Data: msg.Message,
					})
					chatlogs[payload.ChatID] = append(chatlogs[payload.ChatID], ChatMessage{
						Sent:    true,
						User:    data.UserID,
						Message: msg.Message,
					})
				}
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
