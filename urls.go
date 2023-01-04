package fishfish

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type URL struct {
	URL         string   `json:"url"`
	Description string   `json:"description"`
	Category    Category `json:"category"`
	Added       int64    `json:"added"`
	Checked     int64    `json:"checked"`
	Target      string   `json:"target,omitempty"`
}

func (c *RawClient) GetURL(url string) (*URL, error) {
	path := fmt.Sprintf("/urls/%s", url)
	res, err := c.makeRequest("GET", path, nil, nil, c.defaultAuthType)

	if err != nil {
		return nil, err
	}

	return readBody[URL](res)
}

func (c *RawClient) GetURLs(category Category) (*[]string, error) {
	query := makeQuery(map[string]string{
		"category": string(category),
	})
	res, err := c.makeRequest("GET", "/urls", query, nil, c.defaultAuthType)

	if err != nil {
		return nil, err
	}

	return readBody[[]string](res)
}

func (c *RawClient) GetURLsFull() (*[]URL, error) {
	// Requires auth
	if c.defaultAuthType == authTypeNone {
		return nil, errors.New("GetURLsFull requires authentication")
	}

	query := makeQuery(map[string]string{"full": strconv.FormatBool(true)})
	res, err := c.makeRequest("GET", "/urls", query, nil, authTypeSession)

	if err != nil {
		return nil, err
	}

	return readBody[[]URL](res)
}

type CreateURLRequest struct {
	Category    Category `json:"category"`
	Description string   `json:"description"`
	Target      string   `json:"target,omitempty"`
}

func (c *RawClient) AddURL(url string, options CreateURLRequest) (*URL, error) {
	if !c.HasPermission(APIPermissionURLs) {
		return nil, errors.New("missing permission: urls")
	}

	body, err := json.Marshal(options)

	if err != nil {
		return nil, fmt.Errorf("error creating body for AddURL: %s", err)
	}

	path := fmt.Sprintf("/urls/%s", url)
	res, err := c.makeRequest("POST", path, nil, bytes.NewBuffer(body), authTypeSession)

	if err != nil {
		return nil, err
	}

	return readBody[URL](res)
}

type UpdateURLRequest struct {
	Category    Category `json:"category,omitempty"`
	Description string   `json:"description,omitempty"`
	Target      string   `json:"target,omitempty"`
}

func (c *RawClient) UpdateURL(url string, options UpdateURLRequest) error {
	if !c.HasPermission(APIPermissionURLs) {
		return errors.New("missing permission: urls")
	}

	body, err := json.Marshal(options)

	if err != nil {
		return fmt.Errorf("error creating body for UpdateURLs: %s", err)
	}

	path := fmt.Sprintf("/urls/%s", url)
	_, err = c.makeRequest("PATCH", path, nil, bytes.NewBuffer(body), authTypeSession)

	return err
}

func (c *RawClient) DeleteURL(url string) error {
	if !c.HasPermission(APIPermissionURLs) {
		return errors.New("missing permission: urls")
	}

	path := fmt.Sprintf("/urls/%s", url)
	_, err := c.makeRequest("DELETE", path, nil, nil, authTypeSession)

	return err
}
