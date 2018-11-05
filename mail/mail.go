package mail

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func sendAuthEmail(authID string) {
	apiURL := "https://api.sendgrid.com/v3/mail/send" + config.MailAPIKey

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
	}
}
