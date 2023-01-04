package fishfish

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

type APIPermission string
type authType int

const (
	APIPermissionDomains = "domains"
	APIPermissionURLs    = "urls"
	APIPermissionAdmin   = "admin"
)

const (
	authTypePrimary authType = iota
	authTypeSession
	authTypeNone
)

type SessionToken struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}

type CreateSessionTokenRequest struct {
	Permissions []APIPermission `json:"permissions"`
}

func (c *RawClient) CreateSessionToken() (*SessionToken, error) {
	if c.primaryToken == "" {
		return nil, errors.New("invalid authentication token")
	}

	// Apply permissions to token
	body, err := json.Marshal(CreateSessionTokenRequest{Permissions: c.permissions})

	if err != nil {
		return nil, fmt.Errorf("unable to marshal permissions: %s", err)
	}

	res, err := c.makeRequest("POST", "/users/@me/tokens", nil, bytes.NewBuffer(body), authTypePrimary)

	if err != nil {
		// Special 403 case for CreateSessionToken
		if res.StatusCode == 403 {
			return nil, fmt.Errorf("unauthorized for specified permission(s)")
		}

		return nil, err
	}

	return readBody[SessionToken](res)
}

type PartialMainToken struct {
	ID          int64           `json:"id"`
	Permissions []APIPermission `json:"permissions"`
}

func (c *RawClient) GetMainToken(userID, tokenID int64) (*PartialMainToken, error) {
	if !c.HasPermission(APIPermissionAdmin) {
		return nil, errors.New("missing permission: admin")
	}

	path := fmt.Sprintf("/users/%d/tokens/%d", userID, tokenID)
	res, err := c.makeRequest("GET", path, nil, nil, authTypeSession)

	if err != nil {
		return nil, err
	}

	return readBody[PartialMainToken](res)
}

type CreateMainTokenRequest struct {
	Permissions []APIPermission `json:"permissions"`
}

type CreateMainTokenResponse struct {
	Expires int64  `json:"expires"`
	ID      int64  `json:"id"`
	Token   string `json:"token"`
}

func (c *RawClient) CreateMainToken(userID int64, options CreateMainTokenRequest) (*CreateMainTokenResponse, error) {
	if !c.HasPermission(APIPermissionAdmin) {
		return nil, errors.New("missing permission: admin")
	}

	body, err := json.Marshal(options)

	if err != nil {
		return nil, fmt.Errorf("error creating body for CreateMainToken: %s", err)
	}

	path := fmt.Sprintf("/users/%d/tokens", userID)
	res, err := c.makeRequest("GET", path, nil, bytes.NewBuffer(body), authTypeSession)

	if err != nil {
		return nil, err
	}

	return readBody[CreateMainTokenResponse](res)
}

func (c *RawClient) DeleteMainToken(userID, tokenID int64) error {
	if !c.HasPermission(APIPermissionAdmin) {
		return errors.New("missing permission: admin")
	}

	path := fmt.Sprintf("/users/%d/tokens/%d", userID, tokenID)
	_, err := c.makeRequest("DELETE", path, nil, nil, authTypeSession)

	return err
}

// Allow external refresh of the session token

func (c *RawClient) SetSessionToken(token SessionToken) {
	c.sessionToken = token
}

// Check if the client has the specified permission
func (c *RawClient) HasPermission(permission APIPermission) bool {
	for _, v := range c.permissions {
		if v == permission {
			return true
		}
	}

	return false
}
