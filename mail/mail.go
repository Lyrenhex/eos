package mail

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// SendGrid is created in program with an API key and Transactional Template ID
// for SendGrid.com. This structure handles email-related tasks -- sending
// authorisation emails to users to verify new / modified account email
// addresses.
type SendGrid struct {
	APIKey string
	APIID  string
	Email  string
}

type sendgridRequest struct {
	From      sendgridRequestEmail       `json:"from"`
	Personals []sendgridRequestPersonals `json:"personalizations"`
	ID        string                     `json:"template_id"`
}

type sendgridRequestEmail struct {
	Email string `json:"email"`
}

type sendgridRequestPersonals struct {
	To   []sendgridRequestEmail `json:"to"`
	Data sendgridRequestData    `json:"dynamic_template_data"`
}

type sendgridRequestData struct {
	ID string `json:"verifyID"`
}

// SendAuth accepts the recipient's email and the email verification code as
// string parameters. This function sends the user a verification email through
// the Transactional Template specified during SendGrid creation.
func (s *SendGrid) SendAuth(emailAddr, authID string) {
	APIURI := "https://api.sendgrid.com/v3/mail/send"

	client := &http.Client{}

	request := &sendgridRequest{
		From: sendgridRequestEmail{Email: s.Email},
		Personals: []sendgridRequestPersonals{
			sendgridRequestPersonals{
				To: []sendgridRequestEmail{
					sendgridRequestEmail{Email: emailAddr},
				},
				Data: sendgridRequestData{ID: authID},
			},
		},
		ID: s.APIID,
	}

	jsonValue, _ := json.Marshal(request)

	req, _ := http.NewRequest("POST", APIURI, bytes.NewBuffer(jsonValue))
	req.Header.Add("Authorization", "Bearer "+s.APIKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "sendgrid/v3;go")

	client.Do(req)
}
