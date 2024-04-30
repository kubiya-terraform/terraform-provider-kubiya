package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// api/v1/users/self
func (c *Client) self() (*User, error) {
	m := "GET"
	t := "%s/api/v1/users/self"
	uri := c.queryParams(fmt.Sprintf(t, c.host))

	req, err := http.NewRequest(m, uri, nil)
	if err != nil || req == nil {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to create *http.Request")
	}

	body, err := c.doBytesHttpRequest(req)
	if err != nil {
		return nil, err
	}

	var result *User
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// api/v1/users
func (c *Client) users() ([]*User, error) {
	m := "GET"
	t := "%s/api/v1/users"
	uri := c.queryParams(fmt.Sprintf(t, c.host))

	req, err := http.NewRequest(m, uri, nil)
	if err != nil || req == nil {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to create *http.Request")
	}

	body, err := c.doBytesHttpRequest(req)
	if err != nil {
		return nil, err
	}

	var result []*User
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// api/v1/users/groups
func (c *Client) groups() ([]*Group, error) {
	m := "GET"
	t := "%s/api/v1/users/groups"
	uri := c.queryParams(fmt.Sprintf(t, c.host))

	req, err := http.NewRequest(m, uri, nil)
	if err != nil || req == nil {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to create *http.Request")
	}

	body, err := c.doBytesHttpRequest(req)
	if err != nil {
		return nil, err
	}

	var result []*Group
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) GetUserByEmail(email string) (*User, error) {
	items, err := c.users()
	if err != nil {
		return nil, err
	}

	var user User

	for _, item := range items {
		if item.Email == email {
			user = User{
				UUID:       item.UUID,
				Name:       item.Name,
				Email:      item.Email,
				Image:      item.Image,
				Status:     item.Status,
				CreateAt:   item.CreateAt,
				UserStatus: item.UserStatus,
				InviteLink: item.InviteLink,
				Roles:      append(make([]any, 0), item.Roles...),
				Groups:     append(make([]string, 0), item.Groups...),
			}
			break
		}
	}

	return &user, nil
}

func (c *Client) GetGroupByName(name string) (*Group, error) {
	items, err := c.groups()
	if err != nil {
		return nil, err
	}

	var group Group

	for _, item := range items {
		if item.Name == name {
			group = Group{
				UUID:        item.UUID,
				Name:        item.Name,
				System:      item.System,
				CreateAt:    item.CreateAt,
				Description: item.Description,
				Roles:       append(make([]string, 0), item.Roles...),
			}
			break
		}
	}

	return &group, nil
}
