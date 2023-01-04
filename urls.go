package fishfish

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type URL struct {
	URL      string       `json:"url"`
	Category Category     `json:"category"`
	Meta     *URLMetadata `json:"meta,omitempty"`
}

type URLMetadata struct {
	Urlscan string    `json:"urlscan,omitempty"`
	Active  time.Time `json:"active,omitempty"`
	Target  string    `json:"target,omitempty"`
}

func (c *RawClient) GetURLs(category Category, recent bool) (*[]string, error) {
	query := makeQuery(map[string]string{
		"category": string(category),
		"recent":   strconv.FormatBool(recent),
	})
	res, err := c.makeRequest("GET", "/urls", query, nil, c.defaultAuthType)

	if err != nil {
		return nil, err
	}

	return readBody[[]string](res)
}

func (c *RawClient) GetURL(url string, detailed bool) (*URL, error) {
	query := makeQuery(map[string]string{"detailed": strconv.FormatBool(detailed)})
	path := fmt.Sprintf("/urls/%s", url)

	res, err := c.makeRequest("GET", path, query, nil, c.defaultAuthType)

	if err != nil {
		return nil, err
	}

	return readBody[URL](res)
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

func (c *RawClient) AddURL(url string, category Category) (*URL, error) {
	if !c.HasPermission(APIPermissionURLs) {
		return nil, errors.New("missing permission: urls")
	}

	body, err := json.Marshal(URL{
		URL:      url,
		Category: category,
	})

	if err != nil {
		return nil, fmt.Errorf("error creating body for AddURL: %s", err)
	}

	res, err := c.makeRequest("POST", "/urls", nil, bytes.NewBuffer(body), authTypeSession)

	if err != nil {
		return nil, err
	}

	return readBody[URL](res)
}

func (c *RawClient) UpdateURL(url string, category Category) (*URL, error) {
	if !c.HasPermission(APIPermissionURLs) {
		return nil, errors.New("missing permission: urls")
	}

	body, err := json.Marshal(URL{
		Category: category,
	})

	if err != nil {
		return nil, fmt.Errorf("error creating body for UpdateURLs: %s", err)
	}

	path := fmt.Sprintf("/urls/%s", url)
	res, err := c.makeRequest("PATCH", path, nil, bytes.NewBuffer(body), authTypeSession)

	if err != nil {
		return nil, err
	}

	return readBody[URL](res)
}

func (c *RawClient) UpdateURLMetadata(url string, metadata URLMetadata) (*URLMetadata, error) {
	if !c.HasPermission(APIPermissionURLs) {
		return nil, errors.New("missing permission: urls")
	}
	body, err := json.Marshal(metadata)

	if err != nil {
		return nil, fmt.Errorf("error creating body for UpdateURLMetadata: %s", err)
	}

	path := fmt.Sprintf("/urls/%s/metadata", url)
	res, err := c.makeRequest("PATCH", path, nil, bytes.NewBuffer(body), authTypeSession)

	if err != nil {
		return nil, err
	}

	return readBody[URLMetadata](res)
}

func (c *RawClient) DeleteURL(url string) error {
	if !c.HasPermission(APIPermissionURLs) {
		return errors.New("missing permission: urls")
	}

	path := fmt.Sprintf("/urls/%s", url)
	_, err := c.makeRequest("DELETE", path, nil, nil, authTypeSession)

	// No need to check if err is nil, only returning err
	return err
}
