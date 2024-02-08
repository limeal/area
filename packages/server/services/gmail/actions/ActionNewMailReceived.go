package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"time"
)

// `NormalMailResponse` is a struct that contains a slice of structs that contain a string, a string,
// and a struct that contains a slice of structs that contain a string and a string, a struct that
// contains a struct that contains a string, and a string.
// @property {[]struct {
// 		ID       string `json:"id"`
// 		ThreadID string `json:"threadId"`
// 		Payload  struct {
// 			Headers []struct {
// 				Name  string `json:"name"`
// 				Value string `json:"value"`
// 			} `json:"headers"`
// 			Body struct {
// 				Data string `json:"data"`
// 			} `json:"body"`
// 			MimeType string `json:"mimeType"`
// 		} `json:"payload"`
// 	}} Messages - This is an array of messages.
type NormalMailResponse struct {
	Messages []struct {
		ID       string `json:"id"`
		ThreadID string `json:"threadId"`
		Payload  struct {
			Headers []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"headers"`
			Body struct {
				Data string `json:"data"`
			} `json:"body"`
			MimeType string `json:"mimeType"`
		} `json:"payload"`
	} `json:"messages"`
}

// It checks if the user has received a new mail since the last time it was checked
func hasReceivedANewMail(req static.AreaRequest) shared.AreaResponse {
	query := make(map[string]string)

	if (*req.Store)["ctx:time:start"] == nil {
		(*req.Store)["ctx:time:start"] = time.Now().UnixMilli()
	}

	if (*req.Store)["req:mail:from"] != nil {
		query["q"] = "from:" + (*req.Store)["req:mail:from"].(string)
	}

	if (*req.Store)["req:mail:include:spam"] != nil {
		query["includeSpamTrash"] = (*req.Store)["req:mail:include:spam"].(string)
	}

	if (*req.Store)["req:user:id"] != nil {
		(*req.Store)["ctx:user:id"] = (*req.Store)["req:user:id"]
	} else {
		(*req.Store)["ctx:user:id"] = "me"
	}

	encode, _, err := req.Service.Endpoints["GetAllMailEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		(*req.Store)["ctx:user:id"],
		query,
	})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	mailList := NormalMailResponse{}
	errr := json.Unmarshal(encode, &mailList)
	if errr != nil {
		return shared.AreaResponse{Error: errr}
	}

	nbMails := len(mailList.Messages)
	if nbMails == 0 {
		return shared.AreaResponse{Success: false}
	}

	encode, _, err = req.Service.Endpoints["GetMailInfoEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		(*req.Store)["ctx:user:id"],
		mailList.Messages[0].ID,
	})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	var message Message
	if err = json.Unmarshal(encode, &message); err != nil {
		return shared.AreaResponse{Error: err}
	}

	timeMessage, errP := strconv.ParseInt(message.InternalDate, 10, 64)
	if errP != nil {
		return shared.AreaResponse{Error: errP}
	}

	if timeMessage-(*req.Store)["ctx:time:start"].(int64) < 0 {
		return shared.AreaResponse{Success: false}
	}

	findHeader := func(name string) string {
		for _, header := range message.Payload.Headers {
			if header.Name == name {
				return header.Value
			}
		}
		return ""
	}

	body, err := base64.StdEncoding.DecodeString(message.Payload.Parts[0].Body.Data)
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	(*req.Store)["ctx:time:start"] = time.Now().UnixMilli()
	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"gmail:user:id":        (*req.Store)["ctx:user:id"],
			"gmail:mail:id":        mailList.Messages[0].ID,
			"gmail:mail:thread:id": mailList.Messages[0].ThreadID,
			"gmail:mail:from":      findHeader("From"),
			"gmail:mail:subject":   findHeader("Subject"),
			"gmail:mail:date":      findHeader("Date"),
			"gmail:mail:body":      string(body),
			"gmail:mail:body:type": mailList.Messages[0].Payload.MimeType,
		},
	}
}

// It returns a `static.ServiceArea` struct that describes the service area `new_mail_received` for the
// `gmail` service
func DescriptorForGmailActionNewMailReceived() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_mail_received",
		Description: "When a new mail is received",
		Method:      hasReceivedANewMail,
		RequestStore: map[string]static.StoreElement{
			"req:mail:from": {
				Description: "The mail address of the sender",
				Type:        "email",
				Required:    false,
			},
			"req:mail:include:spam": {
				Priority:    2,
				Description: "If the spam mails should be included",
				Type:        "select", // true or false
				Required:    false,
				Values:      []string{"true", "false"},
			},
			"req:user:id": {
				Priority:    1,
				Description: "The user id to check the mail for",
				Type:        "string",
				Required:    false,
			},
		},
		Components: []string{
			"gmail:mail:id",
			"gmail:mail:thread:id",
			"gmail:mail:from",
			"gmail:mail:subject",
			"gmail:mail:date",
			"gmail:mail:body",
			"gmail:mail:body:type",
		},
	}
}
