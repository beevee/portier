package yandex

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/levigross/grequests"
)

// User is a single Yandex.Taxi user
type User struct {
	ID         string `json:"_id,omitempty"`
	Name       string `json:"fullname"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Nickname   string `json:"nickname"`
	CostCenter string `json:"cost_center"`
	IsActive   bool   `json:"is_active"`
	Role       struct {
		ID              string   `json:"role_id,omitempty"`
		Name            string   `json:"name,omitempty"`
		Limit           float32  `json:"limit,omitempty"`
		Classes         []string `json:"classes,omitempty"`
		NoSpecificLimit bool     `json:"no_specific_limit,omitempty"`
		Restrictions    []struct {
			Days      []string `json:"days,omitempty"`
			StartTime string   `json:"start_time,omitempty"`
			EndTime   string   `json:"end_time,omitempty"`
			Type      string   `json:"type,omitempty"`
		} `json:"restrictions,omitempty"`
	} `json:"role"`
}

// UserResponse is Yandex.Taxi API DTO for users
type UserResponse struct {
	Items []User `json:"items"`
}

func (a API) buildURL(action string) string {
	return fmt.Sprintf("%s/client/%s/%s", apiURL, a.clientID, action)
}

// GetUsersByRole fetches all existing users and then filters them by role from Yandex.Taxi API
func (a API) GetUsersByRole(role string) (users map[string]User, err error) {
	usersResponse := &UserResponse{}
	resp, err := a.session.Get(a.buildURL("user"), &grequests.RequestOptions{
		Params: map[string]string{
			"limit": "10000",
		},
	})
	if err != nil {
		return
	}

	stringResponse := strings.Replace(resp.String(), "Infinity", "10000000", -1)
	err = json.Unmarshal([]byte(stringResponse), usersResponse)
	if err != nil {
		return
	}
	if len(usersResponse.Items) == 0 {
		err = fmt.Errorf("API returned zero users")
		return
	}

	roles, err := a.GetRoles()
	if err != nil {
		return
	}

	users = make(map[string]User, len(usersResponse.Items))
	for i := range usersResponse.Items {
		if usersResponse.Items[i].Role.Name == role || roles[usersResponse.Items[i].Role.ID].Name == role {
			users[usersResponse.Items[i].ID] = usersResponse.Items[i]
		}
	}

	return
}

// DisableUser makes it impossible for user to call a taxi
func (a *API) DisableUser(user User) error {
	userID := user.ID
	user.ID = ""
	user.IsActive = false
	user.Role.Name = "" // weirdly, leaving name non-empty leads to 400 Bad Request in this API request
	resp, err := a.session.Put(a.buildURL("user")+"/"+userID+"/", &grequests.RequestOptions{
		JSON: user,
	})
	if err != nil {
		return err
	}
	if !resp.Ok {
		return fmt.Errorf("disable request failed with code %d", resp.StatusCode)
	}
	return nil
}

// EnableUser makes it possible for user to call a taxi
func (a *API) EnableUser(user User) error {
	userID := user.ID
	user.ID = ""
	user.IsActive = true
	user.Role.Name = "" // weirdly, leaving name non-empty leads to 400 Bad Request in this API request
	resp, err := a.session.Put(a.buildURL("user")+"/"+userID+"/", &grequests.RequestOptions{
		JSON: user,
	})
	if err != nil {
		return err
	}
	if !resp.Ok {
		return fmt.Errorf("enable request failed with code %d", resp.StatusCode)
	}
	return nil
}
