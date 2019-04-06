package user

import (
	"encoding/json"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	uuid "github.com/nu7hatch/gouuid"
	"golang.org/x/crypto/bcrypt"
)

// UserIDs is a map of string email addresses returning their associated UserID (of type uuid.UUID)
var UserIDs map[string]uuid.UUID

// Users is a map of uuid.UUID userIDs returning the associated pointer to the intended user's User object.
var Users map[uuid.UUID]*User

// PendingUsers is a map of unique account-generation IDs to the associated email address, for verifying emails.
var PendingUsers map[string]string

// ResetKeys is a map of unique password-reset tokens which are mapped to the associated email address, for bypassing the standard password when necessary.
var ResetKeys map[string]string

// User represents a specific user account, and stores the following data:
//
// - UserID, of type uuid.UUID, a unique identifier for this account used internally
//
// - EmailAddr, a string keeping track of the user's current email address, regardless of verification status (THIS MAY NOT BE THE USER'S LOGIN EMAIL ADDRESS)
//
// - Password, stored as an array of bytes as the result of the bcrypt hashing routine
//
// - Name, a string of the user's provided name (THIS MAY NOT BE THE USER'S FULL / LEGAL NAME)
//
// - Moods, a Moods structure
//
// - Positives, an array of 20 strings storing the user's last 20 positive comments
//
// - Neutrals, an array of 5 strings storing the user's last 5 unemottional comments
//
// - Negatives, an array of 5 strings storing the user's last 5 negative comments
//
// - Admin, a boolean flag representing whether the user has administrator powers or not
//
// - Banned, a boolean flag representing whether the user has had their chatting rights revoked or not
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

// Moods - a sub-structure of User - stores arrays of moods for each day of the week, month of the year, and a Year object
type Moods struct {
	Day   [7]Mood // array of 7 moods, one for each day of week.
	Month [12]Mood
	Years [2]Year // only keep specific data on the past two years. we cannot overload the server.
}

// Mood stores information for a particular time unit, as a Mood integer which acts as a total sum of the user's moods, and a Num integer, serving as a running tally of the number of Moods recorded, for division to form the mean average
type Mood struct {
	Mood int
	Num  int
}

// Year structure to create a copy of a year's Moods structure
type Year struct {
	Year  int
	Month [12]Mood
}

// New creates a new User object with the email address, password, and name  specified, returning a pointer to the created user for later use.
// The User is added to the Users list automatically and the UserID is added to the UserIDs list, which are subsequently saved.
func New(email, pass, name string) *User {
	/* For use during new user creation, generate a new userid and a new User structure to store the new user's data, prepopulated with email pass and name data, with the rest zero-ed for later input. Return the userid for immediate usage, and add to the *memory* user store, and maintain the email-userid pairing in the UserIDs map. */
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
	Users[*uid] = user
	UserIDs[email] = *uid
	go SaveIDs()
	go user.Save()
	return user
}

// Save encodes the User object into basic JSON notation, which is then written to a file specific to that user's UserID.
func (u *User) Save() {
	/* Save the user\"s data to their own file, stored according to their user id. */
	home, err := homedir.Dir()
	if err != nil {
		home = "."
	}
	file, err := os.Create(home + "/eos/data/userdata-" + u.UserID.String() + ".json")
	if err == nil {
		defer file.Close()
		encoder := json.NewEncoder(file)
		err = encoder.Encode(u)
		if err != nil {
			log.Println("Error saving userdata-"+u.UserID.String()+".json: ", err)
		}
	} else {
		log.Println("Error saving userdata-"+u.UserID.String()+".json: ", err)
	}
}

// Load reads the data stored in the file relevant to the provided UserID, uid, and decodes the resulting JSON into the current User object.
func (u *User) Load(uid uuid.UUID) {
	/* Based on the provided user id, load the user\"s data from their own data file, add the user back to the standard users map, and return the user structure for immediate use. */
	home, err := homedir.Dir()
	if err != nil {
		home = "."
	}
	file, err := os.Open(home + "/eos/data/userdata-" + uid.String() + ".json")
	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		err = decoder.Decode(u)
		if err != nil {
			log.Println("Error reading userdata-"+uid.String()+".json: ", err)
		}
		// users[uid] = &user
	} else {
		log.Println("Error reading userdata-"+uid.String()+".json: ", err)
	}
}

// Login runs standard security checks (specifically, comparing the existence of the requested user account and confirming whether the password is correct) before either running Load() on the current User or not, and returning whether the login was successful (account exists and password is correct) and if the account exists (for if the login was unsuccessful). Login does not automatically run New() if the account does not exist.
func (u *User) Login(email, password string) (bool, bool) {
	/* Based on provided email and password strings, check that the login is correct and, if so, load the user\"s data and return UID to calling func. */
	success := true
	exists := true
	uid := UserIDs[email]
	if uid == uuid.UUID([16]byte{}) { // `uid` matches the `uuid.UUID` default value
		success = false
		exists = false
	} else {
		u.Load(uid)
		err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))
		if err != nil {
			if password != ResetKeys[email] {
				success = false
			} else {
				success = true
			}
		}
	}

	return success, exists
}

// ResetPassword checks that the user exists and, if so, generates an auth token for the account, adding it to the ResetKeys map.
func ResetPassword(email string) (bool, string) {
	success := true
	authToken := ""
	uid := UserIDs[email]
	if uid == uuid.UUID([16]byte{}) {
		success = false
	} else {
		token, err := uuid.NewV4()
		if err != nil {
			success = false
		} else {
			authToken = token.String()
			ResetKeys[email] = authToken
		}
	}

	return success, authToken
}

// AddMood updates the User's Mood data, adding the mood onto each Mood object as per the provided day, month, and year.
func (u *User) AddMood(day, month, year, mood int) {
	u.Moods.Day[day].Mood += mood
	u.Moods.Day[day].Num++
	u.Moods.Month[month].Mood += mood
	u.Moods.Month[month].Num++

	yearRecorded := false
	for i, _year := range u.Moods.Years {
		if _year.Year == year {
			u.Moods.Years[i].Month[month].Mood += mood
			u.Moods.Years[i].Month[month].Num++
			yearRecorded = true
		}
	}
	if !yearRecorded {
		newYear := Year{
			Year:  year,
			Month: [12]Mood{},
		}
		u.Moods.Years = [2]Year{
			u.Moods.Years[1],
			newYear,
		}
		u.Moods.Years[1].Month[month].Mood += mood
		u.Moods.Years[1].Month[month].Num++
	}
	u.Save()
}

// AddComment appends a new comment to the end of the current User's comment array based on the supplied mood (which should be an integer betwen -1 and 1, as both negatives and both positives are treated identically)
func (u *User) AddComment(mood int, comment string) {
	if comment != "" {
		switch mood {
		case 1:
			u.Positives = [20]string{
				u.Positives[1],
				u.Positives[2],
				u.Positives[3],
				u.Positives[4],
				u.Positives[5],
				u.Positives[6],
				u.Positives[7],
				u.Positives[8],
				u.Positives[9],
				u.Positives[10],
				u.Positives[11],
				u.Positives[12],
				u.Positives[13],
				u.Positives[14],
				u.Positives[15],
				u.Positives[16],
				u.Positives[17],
				u.Positives[18],
				u.Positives[19],
				comment,
			}
		case 0:
			u.Neutrals = [5]string{
				u.Neutrals[1],
				u.Neutrals[2],
				u.Neutrals[3],
				u.Neutrals[4],
				comment,
			}
		case -1:
			u.Negatives = [5]string{
				u.Negatives[1],
				u.Negatives[2],
				u.Negatives[3],
				u.Negatives[4],
				comment,
			}
		}
		u.Save()
	}
}

// ReadIDs loads the UserIDs stored in persistant memory into the UserIDs array, creating the array if nonexistent.
func ReadIDs() {
	/* Read the email-userid pairs from users.json and store in the UserIDs map for reference when logging in (clients will send an email and password; lookup email in UserIDs to grab their userid, then load their userfile.json) */
	home, err := homedir.Dir()
	if err != nil {
		home = "."
	}
	log.Println("Loading UserIDs from file")
	file, err := os.Open(home + "/eos/data/users.json")
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

// SaveIDs writes the contents of the UserIDs map into persistant memory.
func SaveIDs() {
	/* Save the email-userid pairs, replacing the existing file is present. */
	home, err := homedir.Dir()
	if err != nil {
		home = "."
	}
	log.Println("Saving UserIDs to file")
	file, err := os.Create(home + "/eos/data/users.json")
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
