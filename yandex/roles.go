package yandex

import (
	"encoding/json"
	"strings"

	"github.com/levigross/grequests"
)

// Role is a group of Yandex.Taxi users with shared configuration
type Role struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

// RoleResponse is Yandex.Taxi API DTO for roles
type RoleResponse struct {
	Items []Role `json:"items"`
}

// GetRoles fetches all existing roles from Yandex.Taxi API
func (a API) GetRoles() (roles map[string]Role, err error) {
	resp, err := a.session.Get(apiURL+"/client/"+a.clientID+"/role/", &grequests.RequestOptions{
		Params: map[string]string{
			"limit": "10000",
		},
	})
	if err != nil {
		return
	}

	rolesResponse := &RoleResponse{}
	stringResponse := strings.Replace(resp.String(), "Infinity", "10000000", -1)
	err = json.Unmarshal([]byte(stringResponse), rolesResponse)
	if err != nil {
		return
	}

	roles = make(map[string]Role, len(rolesResponse.Items))
	for i := range rolesResponse.Items {
		roles[rolesResponse.Items[i].ID] = rolesResponse.Items[i]
	}

	return
}
