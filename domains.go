package fishfish

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Domain struct {
	Domain   string          `json:"name"`
	Category Category        `json:"category"`
	Apex     bool            `json:"apex"`
	Meta     *DomainMetadata `json:"meta,omitempty"`
}

type DomainMetadata struct {
	Path    string    `json:"path,omitempty"`
	Urlscan string    `json:"urlscan,omitempty"`
	Active  time.Time `json:"active,omitempty"`
	Target  string    `json:"target,omitempty"`
}

func (c *Client) GetDomains(category Category, recent bool) (*[]string, error) {
	query := makeQuery(map[string]string{
		"category": string(category),
		"recent":   strconv.FormatBool(recent),
	})
	res, err := c.makeRequest("GET", "/domains", query, nil, c.defaultAuthType)

	if err != nil {
		return nil, err
	}

	return readBody[[]string](res)
}

func (c *Client) GetDomain(domain string, detailed bool) (*Domain, error) {
	query := makeQuery(map[string]string{"detailed": strconv.FormatBool(detailed)})
	path := fmt.Sprintf("/domains/%s", domain)

	res, err := c.makeRequest("GET", path, query, nil, c.defaultAuthType)

	if err != nil {
		return nil, err
	}

	return readBody[Domain](res)
}

func (c *Client) GetDomainsFull() (*[]Domain, error) {
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

func (c *Client) AddDomain(domain string, category Category, apex bool) (*Domain, error) {
	if !c.HasPermission(APIPermissionDomains) {
		return nil, errors.New("missing permission: domains")
	}

	body, err := json.Marshal(Domain{
		Domain:   domain,
		Category: category,
		Apex:     apex,
	})

	if err != nil {
		return nil, fmt.Errorf("error creating body for AddDomain: %s", err)
	}

	res, err := c.makeRequest("POST", "/domains", nil, bytes.NewBuffer(body), authTypeSession)

	if err != nil {
		return nil, err
	}

	return readBody[Domain](res)
}

func (c *Client) UpdateDomain(domain string, category Category) (*Domain, error) {
	if !c.HasPermission(APIPermissionDomains) {
		return nil, errors.New("missing permission: domains")
	}

	body, err := json.Marshal(Domain{
		Category: category,
	})

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

func (c *Client) UpdateDomainMetadata(domain string, metadata DomainMetadata) (*DomainMetadata, error) {
	if !c.HasPermission(APIPermissionDomains) {
		return nil, errors.New("missing permission: domains")
	}
	body, err := json.Marshal(metadata)

	if err != nil {
		return nil, fmt.Errorf("error creating body for UpdateDomainMetadata: %s", err)
	}

	path := fmt.Sprintf("/domains/%s/metadata", domain)
	res, err := c.makeRequest("PATCH", path, nil, bytes.NewBuffer(body), authTypeSession)

	if err != nil {
		return nil, err
	}

	return readBody[DomainMetadata](res)
}

func (c *Client) DeleteDomain(domain string) error {
	if !c.HasPermission(APIPermissionDomains) {
		return errors.New("missing permission: domains")
	}

	path := fmt.Sprintf("/domains/%s", domain)
	_, err := c.makeRequest("DELETE", path, nil, nil, authTypeSession)

	// No need to check if err is nil, only returning err
	return err
}
