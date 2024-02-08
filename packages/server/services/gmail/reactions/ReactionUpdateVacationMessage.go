package reactions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"encoding/json"
	"errors"
	"strconv"
)

// It updates the vacation message of a user
func updateVacationMessage(req static.AreaRequest) shared.AreaResponse {

	if (*req.Store)["ctx:user:id"] == nil {
		if (*req.Store)["req:user:id"] != nil {
			(*req.Store)["ctx:user:id"] = (*req.Store)["req:user:id"]
		} else {
			(*req.Store)["ctx:user:id"] = "me"
		}
	}

	(*req.Store)["ctx:auto:reply"] = true
	if (*req.Store)["req:auto:reply"] != nil {
		var err error
		(*req.Store)["ctx:auto:reply"], err = strconv.ParseBool((*req.Store)["req:auto:reply"].(string))
		if err != nil {
			return shared.AreaResponse{Error: err}
		}
	}

	if (*req.Store)["ctx:response:subject"] == nil {
		(*req.Store)["ctx:response:subject"] = utils.GenerateFinalComponent((*req.Store)["req:response:subject"].(string), req.ExternalData, []string{})
	}

	if (*req.Store)["ctx:response:body"] == nil {
		(*req.Store)["ctx:response:body"] = utils.GenerateFinalComponent((*req.Store)["req:response:body"].(string), req.ExternalData, []string{})
	}

	body, err := json.Marshal(map[string]interface{}{
		"enableAutoReply":  (*req.Store)["ctx:auto:reply"],
		"responseSubject":  (*req.Store)["ctx:response:subject"],
		"responseBodyHtml": (*req.Store)["ctx:response:body"],
	})

	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	_, httpResp, errr := req.Service.Endpoints["UpdateVacationMessageEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		(*req.Store)["ctx:user:id"],
		string(body),
	})

	if httpResp != nil && httpResp.StatusCode == 403 {
		return shared.AreaResponse{Error: errors.New("You don't have the permission to update the vacation message")}
	}

	return shared.AreaResponse{Error: errr}
}

// `DescriptorForGmailReactionUpdateVacationMessage` returns a `static.ServiceArea` that describes the
// `updateVacationMessage` function
func DescriptorForGmailReactionUpdateVacationMessage() static.ServiceArea {
	return static.ServiceArea{
		Name:        "update_vacation_message",
		Description: "Update the vacation message for the user",
		RequestStore: map[string]static.StoreElement{
			"req:response:subject": {
				Priority:    1,
				Description: "The subject line of the auto-reply message.",
				Required:    true,
			},
			"req:response:body": {
				Priority:    2,
				Description: "The response body in HTML format. Gmail will sanitize the HTML before storing it.",
				Required:    true,
			},
			"req:auto:reply": {
				Priority:    3,
				Type:        "select",
				Description: "Enables auto-replying to messages. (default: true)",
				Required:    false,
				Values:      []string{"true", "false"},
			},
			"req:user:id": {
				Type:        "string",
				Description: "The user id to update the vacation message for",
				Required:    false,
			},
		},
		Method: updateVacationMessage,
	}
}
