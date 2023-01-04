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
)

const (
	authTypePrimary authType = iota
	authTypeSession
	authTypeNone
)

type sessionToken struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}

type createSessionToken struct {
	Permissions []APIPermission `json:"permissions"`
}

func (c *RawClient) CreateSessionToken() (*sessionToken, error) {
	if c.primaryToken == "" {
		return nil, errors.New("invalid authentication token")
	}

	// Apply permissions to token
	body, err := json.Marshal(createSessionToken{Permissions: c.permissions})

	if err != nil {
		return nil, fmt.Errorf("unable to marshal permissions: %s", err)
	}

	res, err := c.makeRequest("POST", "/users/@me/tokens", nil, bytes.NewBuffer(body), authTypePrimary)

	if err != nil {
		return nil, err
	}

	return readBody[sessionToken](res)
}

// Allow external refresh of the session token

func (c *RawClient) UpdateSessionToken(token sessionToken) {
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
