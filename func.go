/*
 * Eos Backend Server
 *
 * Copyright (c) Damian Heaton 2017 All rights reserved.
 *
 * Server Functions
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nu7hatch/gouuid"
	"golang.org/x/crypto/bcrypt"
)

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
