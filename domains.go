package fishfish

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type Domain struct {
	Domain      string   `json:"name"`
	Description string   `json:"description"`
	Category    Category `json:"category"`
	Added       int64    `json:"added"`
	Checked     int64    `json:"checked"`
	Target      string   `json:"target,omitempty"`
}

func (c *RawClient) GetDomain(domain string) (*Domain, error) {
	path := fmt.Sprintf("/domains/%s", domain)
	res, err := c.makeRequest("GET", path, nil, nil, c.defaultAuthType)

	if err != nil {
		return nil, err
	}

	return readBody[Domain](res)
}

func (c *RawClient) GetDomains(category Category) (*[]string, error) {
	query := makeQuery(map[string]string{
		"category": string(category),
	})
	res, err := c.makeRequest("GET", "/domains", query, nil, c.defaultAuthType)

	if err != nil {
		return nil, err
	}

	return readBody[[]string](res)
}

func (c *RawClient) GetDomainsFull() (*[]Domain, error) {
	// Requires auth
	if c.defaultAuthType == authTypeNone {
		return nil, errors.New("GetDomainsFull requires authentication")
	}

	query := makeQuery(map[string]string{"full": strconv.FormatBool(true)})
	res, err := c.makeRequest("GET", "/domains", query, nil, authTypeSession)

	if err != nil {
		return nil, err
	}

	return readBody[[]Domain](res)
}

type CreateDomainRequest struct {
	Category    Category `json:"category"`
	Description string   `json:"description"`
	Target      string   `json:"target,omitempty"`
}

func (c *RawClient) AddDomain(domain string, options CreateDomainRequest) (*Domain, error) {
	if !c.HasPermission(APIPermissionDomains) {
		return nil, errors.New("missing permission: domains")
	}

	body, err := json.Marshal(options)

	if err != nil {
		return nil, fmt.Errorf("error creating body for AddDomain: %s", err)
	}

	path := fmt.Sprintf("/domains/%s", domain)
	res, err := c.makeRequest("POST", path, nil, bytes.NewBuffer(body), authTypeSession)

	if err != nil {
		return nil, err
	}

	return readBody[Domain](res)
}

type UpdateDomainRequest struct {
	Category    Category `json:"category,omitempty"`
	Description string   `json:"description,omitempty"`
	Target      string   `json:"target,omitempty"`
}

func (c *RawClient) UpdateDomain(domain string, options UpdateDomainRequest) (*Domain, error) {
	if !c.HasPermission(APIPermissionDomains) {
		return nil, errors.New("missing permission: domains")
	}

	body, err := json.Marshal(options)

	if err != nil {
		return nil, fmt.Errorf("error creating body for UpdateDomain: %s", err)
	}

	path := fmt.Sprintf("/domains/%s", domain)
	res, err := c.makeRequest("PATCH", path, nil, bytes.NewBuffer(body), authTypeSession)

	if err != nil {
		return nil, err
	}

	return readBody[Domain](res)
}

func (c *RawClient) DeleteDomain(domain string) error {
	if !c.HasPermission(APIPermissionDomains) {
		return errors.New("missing permission: domains")
	}

	path := fmt.Sprintf("/domains/%s", domain)
	_, err := c.makeRequest("DELETE", path, nil, nil, authTypeSession)

	// No need to check if err is nil, only returning err
	return err
}
