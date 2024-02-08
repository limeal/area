package actions

import (
	"area-server/classes/shared"
	"area-server/classes/static"
	"area-server/utils"
	"encoding/json"
)

// `ListFileMembersResponse` is a struct with three fields: `Groups`, `Invitees`, and `Users`.
//
// The `Groups` field is an array of structs, each of which has three fields: `AccessType`, `Group`,
// and `IsInherited`.
//
// The `AccessType` field is a struct with one field: `Tag`.
//
// The `Tag` field is a string.
//
// The `Group` field is a struct with seven fields: `GroupID`, `GroupManagementType`, `GroupName`,
// `Group
// @property {[]struct {
// 		AccessType struct {
// 			Tag string `json:".tag"`
// 		} `json:"access_type"`
// 		Group struct {
// 			GroupID             string `json:"group_id"`
// 			GroupManagementType struct {
// 				Tag string `json:".tag"`
// 			} `json:"group_management_type"`
// 			GroupName string `json:"group_name"`
// 			GroupType struct {
// 				Tag string `json:".tag"`
// 			} `json:"group_type"`
// 			IsMember    bool `json:"is_member"`
// 			IsOwner     bool `json:"is_owner"`
// 			MemberCount int  `json:"member_count"`
// 			SameTeam    bool `json:"same_team"`
// 		} `json:"group"`
// 		IsInherited bool `json:"is_inherited"`
// 	}} Groups - A list of groups that have access to the file.
// @property {[]struct {
// 		AccessType struct {
// 			Tag string `json:".tag"`
// 		} `json:"access_type"`
// 		Invitee struct {
// 			Tag   string `json:".tag"`
// 			Email string `json:"email"`
// 		} `json:"invitee"`
// 		IsInherited bool `json:"is_inherited"`
// 	}} Invitees - A list of users who have been invited to this file.
// @property {[]struct {
// 		AccessType struct {
// 			Tag string `json:".tag"`
// 		} `json:"access_type"`
// 		User struct {
// 			AccountID    string `json:"account_id"`
// 			DisplayName  string `json:"display_name"`
// 			Email        string `json:"email"`
// 			SameTeam     bool   `json:"same_team"`
// 			TeamMemberID string `json:"team_member_id"`
// 		} `json:"user"`
// 		IsInherited bool `json:"is_inherited"`
// 	}} Users - A list of users who have access to the file.
type ListFileMembersResponse struct {
	Groups []struct {
		AccessType struct {
			Tag string `json:".tag"`
		} `json:"access_type"`
		Group struct {
			GroupID             string `json:"group_id"`
			GroupManagementType struct {
				Tag string `json:".tag"`
			} `json:"group_management_type"`
			GroupName string `json:"group_name"`
			GroupType struct {
				Tag string `json:".tag"`
			} `json:"group_type"`
			IsMember    bool `json:"is_member"`
			IsOwner     bool `json:"is_owner"`
			MemberCount int  `json:"member_count"`
			SameTeam    bool `json:"same_team"`
		} `json:"group"`
		IsInherited bool `json:"is_inherited"`
	} `json:"groups"`
	Invitees []struct {
		AccessType struct {
			Tag string `json:".tag"`
		} `json:"access_type"`
		Invitee struct {
			Tag   string `json:".tag"`
			Email string `json:"email"`
		} `json:"invitee"`
		IsInherited bool `json:"is_inherited"`
	} `json:"invitees"`
	Users []struct {
		AccessType struct {
			Tag string `json:".tag"`
		} `json:"access_type"`
		User struct {
			AccountID    string `json:"account_id"`
			DisplayName  string `json:"display_name"`
			Email        string `json:"email"`
			SameTeam     bool   `json:"same_team"`
			TeamMemberID string `json:"team_member_id"`
		} `json:"user"`
		IsInherited bool `json:"is_inherited"`
	} `json:"users"`
}

type CommonEntity struct {
	ID           string `json:"id"`             // (group -> GroupID, user -> AccountID, invitee -> Email)
	Name         string `json:"name"`           // (group -> GroupName, user -> DisplayName, invitee -> Email)
	Type         string `json:"type"`           // (group -> GroupType.Tag, user -> "", invitee -> "")
	Email        string `json:"email"`          // (group -> "", user -> Email, invitee -> Email)
	SameTeam     bool   `json:"same_team"`      // (group -> false, user -> SameTeam, invitee -> false)
	TeamMemberID string `json:"team_member_id"` // (group -> "", user -> TeamMemberID, invitee -> "")
	IsMember     bool   `json:"is_member"`      // (group -> IsMember, user -> false, invitee -> false)
	IsOwner      bool   `json:"is_owner"`       // (group -> IsOwner, user -> false, invitee -> false)
	MemberCount  int    `json:"member_count"`   // (group -> MemberCount, user -> 0, invitee -> 0)
	IsInherited  bool   `json:"is_inherited"`   // (group -> IsInherited, user -> IsInherited, invitee -> IsInherited)
}

// It checks if a file has been shared to a new person
func onUserShareAFileToNewPerson(req static.AreaRequest) shared.AreaResponse {

	// Check file members
	// Check if the file has been shared to a new person
	body, err := json.Marshal(map[string]interface{}{
		"file":              (*req.Store)["req:file:id"],
		"include_inherited": true,
		"limit":             300,
	})
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	if (*req.Store)["ctx:entity:type"] == nil {
		(*req.Store)["ctx:entity:type"] = "user"
		if (*req.Store)["req:entity:type"] != nil {
			(*req.Store)["ctx:entity:type"] = (*req.Store)["req:entity:type"]
		}
	}

	entities, _, errr := req.Service.Endpoints["ListFileMembersEndpoint"].CallEncode([]interface{}{
		req.Authorization,
		string(body),
	})

	if errr != nil {
		return shared.AreaResponse{Error: errr}
	}

	var listFileMembersResponse ListFileMembersResponse
	err = json.Unmarshal(entities, &listFileMembersResponse)
	if err != nil {
		return shared.AreaResponse{Error: err}
	}

	var nbEntities int
	var newEntity CommonEntity
	switch (*req.Store)["ctx:entity:type"] {
	case "user":
		nbEntities = len(listFileMembersResponse.Users)
		newEntity = CommonEntity{
			ID:           listFileMembersResponse.Users[0].User.AccountID,
			Name:         listFileMembersResponse.Users[0].User.DisplayName,
			Type:         "",
			Email:        listFileMembersResponse.Users[0].User.Email,
			SameTeam:     listFileMembersResponse.Users[0].User.SameTeam,
			TeamMemberID: listFileMembersResponse.Users[0].User.TeamMemberID,
			IsMember:     false,
			IsOwner:      false,
			MemberCount:  0,
			IsInherited:  listFileMembersResponse.Users[0].IsInherited,
		}
	case "group":
		nbEntities = len(listFileMembersResponse.Groups)
		newEntity = CommonEntity{
			ID:           listFileMembersResponse.Groups[0].Group.GroupID,
			Name:         listFileMembersResponse.Groups[0].Group.GroupName,
			Type:         listFileMembersResponse.Groups[0].Group.GroupType.Tag,
			Email:        "",
			SameTeam:     false,
			TeamMemberID: "",
			IsMember:     listFileMembersResponse.Groups[0].Group.IsMember,
			IsOwner:      listFileMembersResponse.Groups[0].Group.IsOwner,
			MemberCount:  listFileMembersResponse.Groups[0].Group.MemberCount,
			IsInherited:  listFileMembersResponse.Groups[0].IsInherited,
		}
	case "invitee":
		nbEntities = len(listFileMembersResponse.Invitees)
		newEntity = CommonEntity{
			ID:           listFileMembersResponse.Invitees[0].Invitee.Email,
			Name:         listFileMembersResponse.Invitees[0].Invitee.Email,
			Type:         "",
			Email:        listFileMembersResponse.Invitees[0].Invitee.Email,
			SameTeam:     false,
			TeamMemberID: "",
			IsMember:     false,
			IsOwner:      false,
			MemberCount:  0,
			IsInherited:  listFileMembersResponse.Invitees[0].IsInherited,
		}
	}

	ok, errL := utils.IsLatestBasic(req.Store, nbEntities)
	if errL != nil {
		return shared.AreaResponse{Error: errL}
	}
	if !ok {
		return shared.AreaResponse{Success: false}
	}

	return shared.AreaResponse{
		Success: true,
		Data: map[string]interface{}{
			"dropbox:entity:id":             newEntity.ID,
			"dropbox:entity:name":           newEntity.Name,
			"dropbox:entity:type":           newEntity.Type,
			"dropbox:entity:email":          newEntity.Email,
			"dropbox:entity:same_team":      newEntity.SameTeam,
			"dropbox:entity:team_member_id": newEntity.TeamMemberID,
			"dropbox:entity:is_member":      newEntity.IsMember,
			"dropbox:entity:is_owner":       newEntity.IsOwner,
			"dropbox:entity:member_count":   newEntity.MemberCount,
			"dropbox:entity:is_inherited":   newEntity.IsInherited,
		},
	}
}

// It returns a service area that describes the Dropbox action "User share a file to new entity"
func DescriptorForDropboxActionUserShareAFileToNewEntity() static.ServiceArea {
	return static.ServiceArea{
		Name:        "user_share_a_file_to_new_entity",
		Description: "User share a file",
		RequestStore: map[string]static.StoreElement{
			"req:file:id": {
				Description: "The file to check if it has been shared to a new person",
				Type:        "select_uri",
				Required:    true,
				Values:      []string{"/files?path=${req:directory:path}"},
			},
			"req:directory:path": {
				Priority:    1,
				Description: "The path to the directory",
				Type:        "select_uri",
				Required:    false,
				Values:      []string{"/directories"},
			},
			"req:entity:type": {
				Priority:    1,
				Description: "The type of the entity (group, user or invitee) (default: user)",
				Type:        "select",
				Required:    false,
				Values: []string{
					"group",
					"user",
					"invitee",
				},
			},
		},
		Method: onUserShareAFileToNewPerson,
		Components: []string{
			"dropbox:entity:id",
			"dropbox:entity:name",
			"dropbox:entity:type",
			"dropbox:entity:email",
			"dropbox:entity:same_team",
			"dropbox:entity:team_member_id",
			"dropbox:entity:is_member",
			"dropbox:entity:is_owner",
			"dropbox:entity:member_count",
			"dropbox:entity:is_inherited",
		},
	}
}
