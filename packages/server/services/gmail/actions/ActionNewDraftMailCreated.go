package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"time"
)

// `Message` is a struct that contains a string `ID`, a string `ThreadID`, a struct `Payload`, and a
// string `InternalDate`.
//
// `Payload` is a struct that contains a slice of structs `Headers` and a slice of structs `Parts`.
//
// `Headers` is a slice of structs that contain a string `Name` and a string `Value`.
//
// `Parts` is a slice of structs that contain a struct `Body` and a string `MimeType`.
//
// `Body`
// @property {string} ID - The ID of the message.
// @property {string} ThreadID - The ID of the thread the message belongs to.
// @property Payload - This is the actual email message.
// @property {string} InternalDate - The date the message was received by Gmail.
type Message struct {
	ID       string `json:"id"`
	ThreadID string `json:"threadId"`
	Payload  struct {
		Headers []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"headers"`
		Parts []struct {
			Body struct {
				Data string `json:"data"`
			} `json:"body"`
			MimeType string `json:"mimeType"`
		} `json:"parts"`
	} `json:"payload"`
	InternalDate string `json:"internalDate"`
}

// `PartialMessage` is a struct with two fields, `ID` and `ThreadID`, both of which are strings.
// @property {string} ID - The ID of the message.
// @property {string} ThreadID - The ID of the thread that the message belongs to.
type PartialMessage struct {
	ID       string `json:"id"`
	ThreadID string `json:"threadId"`
}

// `GetAllDraftMailResponse` is a struct with two fields, `Drafts` and `ResultSizeEstimate`. `Drafts`
// is a slice of structs with two fields, `ID` and `Message`. `Message` is a struct with two fields,
// `ID` and `ThreadID`.
// @property {[]struct {
// 		ID      string         `json:"id"`
// 		Message PartialMessage `json:"message"`
// 	}} Drafts - An array of draft messages.
// @property {int} ResultSizeEstimate - The approximate number of messages in the mailbox.
type GetAllDraftMailResponse struct {
	Drafts []struct {
		ID      string         `json:"id"`
		Message PartialMessage `json:"message"`
	} `json:"drafts"`
	ResultSizeEstimate int `json:"resultSizeEstimate"`
}

// It checks if a new draft mail has been created since the last time it was called
func hasANewDraftMailBeenCreated(req static.AreaRequest) shared.AreaResponse {

	query := make(map[string]string)

	if (*req.Store)["ctx:time:start"] == nil {
		(*req.Store)["ctx:time:start"] = time.Now().UnixMilli()
	}

	query["maxResults"] = "1"
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

	encode, _, err := req.Service.Endpoints["GetAllDraftMailEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		(*req.Store)["ctx:user:id"],
		query,
	})
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	draftList := GetAllDraftMailResponse{}
	if errr := json.Unmarshal(encode, &draftList); errr != nil {
		return shared.AreaResponse{Error: errr}
	}

	nbDrafts := len(draftList.Drafts)
	if nbDrafts == 0 {
		return shared.AreaResponse{Success: false}
	}

	encode, _, err = req.Service.Endpoints["GetMailInfoEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		(*req.Store)["ctx:user:id"],
		draftList.Drafts[0].Message.ID,
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
			"gmail:user:id":         (*req.Store)["ctx:user:id"],
			"gmail:draft:id":        draftList.Drafts[0].ID,
			"gmail:draft:thread:id": message.ThreadID,
			"gmail:draft:from":      findHeader("From"),
			"gmail:draft:to":        findHeader("To"),
			"gmail:draft:date":      findHeader("Date"),
			"gmail:draft:subject":   findHeader("Subject"),
			"gmail:draft:body":      string(body),
			"gmail:draft:body:type": message.Payload.Parts[0].MimeType,
		},
	}
}

// It returns a static.ServiceArea object that describes the service area
func DescriptorForGmailActionNewDraftMailCreated() static.ServiceArea {
	return static.ServiceArea{
		Name:        "new_draft_mail_created",
		Description: "When a new draft mail is created",
		RequestStore: map[string]static.StoreElement{
			"req:mail:from": {
				Description: "The mail address of the sender",
				Type:        "email",
				Required:    false,
			},
			"req:mail:include:spam": {
				Priority:    2,
				Description: "If the spam mails should be included",
				Type:        "select",
				Required:    false,
				Values:      []string{"true", "false"},
			},
			"req:user:id": {
				Priority:    1,
				Description: "The id of the user to check if a new draft mail has been created",
				Type:        "string",
				Required:    false,
			},
		},
		Method: hasANewDraftMailBeenCreated,
		Components: []string{
			"gmail:draft:id",
			"gmail:draft:thread:id",
			"gmail:draft:from",
			"gmail:draft:to",
			"gmail:draft:date",
			"gmail:draft:subject",
			"gmail:draft:body",
			"gmail:draft:body:type",
		},
	}
}
