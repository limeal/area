package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/textproto"

	message "github.com/sloonz/go-mime-message"
	"github.com/sloonz/go-qprintable"
)

// ProfileResponse is a struct that contains a field called Email that is a string.
// @property {string} Email - The email address of the user.
type ProfileResponse struct {
	Email string `json:"emailAddress"`
}

// It creates a new mail and sends it to the specified recipient
func createNewMail(req static.AreaRequest) shared.AreaResponse {
	if (*req.Store)["ctx:user:id"] == nil {
		if (*req.Store)["req:user:id"] != nil {
			(*req.Store)["ctx:user:id"] = (*req.Store)["req:user:id"]
		} else {
			(*req.Store)["ctx:user:id"] = "me"
		}
	}

	if (*req.Store)["ctx:mail:from"] == nil {
		encode, _, err := req.Service.Endpoints["GetProfileEndpoint"].CallEncode([]interface{}{
			req.Authorization,
			(*req.Store)["ctx:user:id"],
		})
		if err != nil {
			return shared.AreaResponse{Error: err}
		}

		profile := ProfileResponse{}
		errr := json.Unmarshal(encode, &profile)
		if errr != nil {
			return shared.AreaResponse{Error: errr}
		}

		(*req.Store)["ctx:mail:from"] = profile.Email
	}

	if (*req.Store)["ctx:mail:to"] == nil {
		(*req.Store)["ctx:mail:to"] = utils.GenerateFinalComponent((*req.Store)["req:mail:to"].(string), req.ExternalData, []string{
			"gmail:mail:author",
		})
	}

	if (*req.Store)["ctx:mail:subject"] == nil {
		(*req.Store)["ctx:mail:subject"] = utils.GenerateFinalComponent((*req.Store)["req:mail:subject"].(string), req.ExternalData, []string{})
	}

	if (*req.Store)["ctx:mail:body"] == nil {
		(*req.Store)["ctx:mail:body"] = utils.GenerateFinalComponent((*req.Store)["req:mail:body"].(string), req.ExternalData, []string{})
	}

	mail := message.NewMultipartMessage("alternative", "")
	mail.SetHeader("From", "<"+(*req.Store)["ctx:mail:from"].(string)+">")
	mail.SetHeader("To", "<"+(*req.Store)["ctx:mail:to"].(string)+">")
	mail.SetHeader("Subject", message.EncodeWord((*req.Store)["ctx:mail:subject"].(string)))
	m2 := message.NewTextMessage(qprintable.UnixTextEncoding, bytes.NewBufferString((*req.Store)["ctx:mail:body"].(string)))
	m2.SetHeader("Content-Type", "text/plain")
	mail.AddPart(m2)

	mbuf := bufio.NewReader(mail)
	if mbuf == nil {
		return shared.AreaResponse{Error: errors.New("Error while creating mime message")}
	}

	tp := textproto.NewReader(mbuf)
	if tp == nil {
		return shared.AreaResponse{Error: errors.New("Error while creating mime message")}
	}

	raw := ""
	for content, err := tp.ReadLine(); err == nil; content, err = tp.ReadLine() {
		raw += content + "\r\n"
	}

	raw = base64.URLEncoding.EncodeToString([]byte(raw))
	_, _, err := req.Service.Endpoints["SendMailEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		(*req.Store)["ctx:user:id"],
		raw,
	})
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	return shared.AreaResponse{
		Error: nil,
	}
}

// It returns a static.ServiceArea object that describes the service area
func DescriptorForGmailReactionSendMail() static.ServiceArea {
	return static.ServiceArea{
		Name:        "create_new_mail",
		Description: "Create a new mail",
		Method:      createNewMail,
		RequestStore: map[string]static.StoreElement{
			"req:mail:to": {
				Priority:    1,
				Description: "The email address of the recipient",
				Type:        "email",
				Required:    true,
			},
			"req:mail:subject": {
				Priority:    2,
				Description: "The subject of the email",
				Type:        "string",
				Required:    true,
			},
			"req:mail:body": {
				Priority:    3,
				Description: "The body of the email",
				Type:        "text",
				Required:    true,
			},
			"req:user:id": {
				Description: "The id of the user that will send the mail",
				Type:        "string",
				Required:    false,
			},
		},
	}
}
