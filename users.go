package fishfish

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

type User struct {
	ExternalServiceID string          `json:"external_service_id,omitempty"`
	ID                int64           `json:"id"`
	Permissions       []APIPermission `json:"permissions"`
	Username          string          `json:"username"`
}

func (c *RawClient) GetUser(id int64) (*User, error) {
	if !c.HasPermission(APIPermissionAdmin) {
		return nil, errors.New("missing permission: admin")
	}

	path := fmt.Sprintf("/users/%d", id)
	res, err := c.makeRequest("GET", path, nil, nil, authTypeSession)

	if err != nil {
		return nil, err
	}

	return readBody[User](res)
}

type CreateUserRequest struct {
	ExternalServiceID string `json:"external_service_id,omitempty"`
	Username          string `json:"username"`
}

func (c *RawClient) CreateUser(options CreateDomainRequest) (*User, error) {
	if !c.HasPermission(APIPermissionAdmin) {
		return nil, errors.New("missing permission: admin")
	}

	body, err := json.Marshal(options)

	if err != nil {
		return nil, fmt.Errorf("error creating body for CreateUser: %s", err)
	}

	res, err := c.makeRequest("POST", "/users", nil, bytes.NewBuffer(body), authTypeSession)

	if err != nil {
		return nil, err
	}

	return readBody[User](res)
}

type UpdateUserRequest struct {
	Permissions []APIPermission `json:"permissions"`
	Username    string          `json:"username"`
}

func (c *RawClient) UpdateUser(id int64, options UpdateUserRequest) error {
	if !c.HasPermission(APIPermissionAdmin) {
		return errors.New("missing permission: admin")
	}

	body, err := json.Marshal(options)

	if err != nil {
		return fmt.Errorf("error creating body for UpdateUser: %s", err)
	}

	path := fmt.Sprintf("/users/%d", id)
	_, err = c.makeRequest("PATCH", path, nil, bytes.NewBuffer(body), authTypeSession)

	return err
}

func (c *RawClient) DeleteUser(id int64) error {
	if !c.HasPermission(APIPermissionAdmin) {
		return errors.New("missing permission: admin")
	}

	path := fmt.Sprintf("/users/%d", id)
	_, err := c.makeRequest("DELETE", path, nil, nil, authTypeSession)

	return err
}
