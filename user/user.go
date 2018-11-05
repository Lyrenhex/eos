package user

import (
	"encoding/json"
	"log"
	"os"

	uuid "github.com/nu7hatch/gouuid"
	"golang.org/x/crypto/bcrypt"
)

// UserIDs maps email address to user id
var UserIDs map[string]uuid.UUID // @UserIDs["email address"] => uuid.UUID
var Users map[uuid.UUID]*User    // @Users[uuid.UUID], store a pointer to the user\"s data - avoid memory data duplication - edit the storage location directly.

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

func New(email, pass, name string) *User {
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
	Users[*uid] = user
	UserIDs[email] = *uid
	go SaveIDs()
	go user.Save()
	return user
}

func (u *User) Save() {
	/* Save the user\"s data to their own file, stored according to their user id. */
	file, err := os.Create("data/userdata-" + u.UserID.String() + ".json")
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

func (u *User) Load(uid uuid.UUID) {
	/* Based on the provided user id, load the user\"s data from their own data file, add the user back to the standard users map, and return the user structure for immediate use. */
	file, err := os.Open("data/userdata-" + uid.String() + ".json")
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
			success = false
		}
	}

	return success, exists
}

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

func ReadIDs() {
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

func SaveIDs() {
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
