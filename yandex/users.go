package yandex

import (
	"fmt"

	"github.com/levigross/grequests"
)

type User struct {
	ID         string `json:"_id,omitempty"`
	Name       string `json:"fullname"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Nickname   string `json:"nickname"`
	CostCenter string `json:"cost_center"`
	IsActive   bool   `json:"is_active"`
	Role       struct {
		ID      string   `json:"role_id,omitempty"`
		Name    string   `json:"name,omitempty"`
		Limit   float32  `json:"limit,omitempty"`
		Classes []string `json:"classes,omitempty"`
	} `json:"role"`
}

type UserResponse struct {
	Items []User `json:"items"`
}

func (a API) buildURL(action string) string {
	return fmt.Sprintf("%s/client/%s/%s", apiURL, a.clientID, action)
}

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

	err = resp.JSON(usersResponse)
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

func (a *API) DisableUser(user User) error {
	userID := user.ID
	user.ID = ""
	user.IsActive = false
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

func (a *API) EnableUser(user User) error {
	userID := user.ID
	user.ID = ""
	user.IsActive = true
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
