package fishfish

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type AutoSyncClient struct {
	raw           RawClient
	cache         domainCache
	cacheTicker   *time.Ticker
	sessionTicker *time.Ticker
	context       syncContext
}

type domainCache struct {
	mx          sync.RWMutex
	domainIndex map[string]Domain
	urlIndex    map[string]URL
}

type syncContext struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewAutoSync(primaryToken string, permissions []APIPermission) (*AutoSyncClient, error) {
	rawClient, err := NewRaw(primaryToken, permissions)

	if err != nil {
		return nil, err
	}

	client := AutoSyncClient{
		raw:   *rawClient,
		cache: domainCache{},
	}

	return &client, nil
}

func (c *AutoSyncClient) ForceSync() error {
	c.cache.mx.Lock()
	defer c.cache.mx.Unlock()

	domains, err := c.raw.GetDomainsFull()

	if err != nil {
		return fmt.Errorf("failed to sync: %s", err)
	}

	domainArr := *domains
	c.cache.domainIndex = map[string]Domain{}

	// Index domains in map for faster lookup
	for i, domain := range domainArr {
		c.cache.domainIndex[domain.Domain] = domainArr[i]
	}

	return nil
}

func (c *AutoSyncClient) StartAutoSync() {
	context, cancel := context.WithCancel(context.Background())
	c.context.ctx = context
	c.context.cancel = cancel

	// Force update the cache every hour
	c.cacheTicker = time.NewTicker(time.Hour)
	// Session tokens expire after an hour, refresh 15 minutes early just in case
	c.sessionTicker = time.NewTicker(time.Second * 10)

	// Generate Session Token
	// The client should already be able to successfully create a token from initialization
	token, _ := c.raw.CreateSessionToken()
	c.raw.SetSessionToken(*token)

	// Initial Sync
	c.ForceSync()

	// Start automatically syncing domains/urls
	go func(client *AutoSyncClient) {
		for {
			select {
			case <-c.cacheTicker.C:
				c.ForceSync()
			case <-c.context.ctx.Done():
				return
			}
		}
	}(c)

	// Start automatically refreshing the session token
	go func(client *AutoSyncClient) {
		for {
			select {
			case <-c.sessionTicker.C:
				token, _ := c.raw.CreateSessionToken()
				c.raw.SetSessionToken(*token)
			case <-c.context.ctx.Done():
				return
			}
		}
	}(c)

	// Start the websocket to add new domains
	go func(client *AutoSyncClient) {
		ch := make(chan WSEvent)
		go client.raw.ConnectWS(c.context.ctx, ch)

		for {
			data := <-ch

			if data.Data == nil {
				// WebSocket was closed
				return
			}

			c.cache.mx.Lock()

			dataAsMap := data.Data.(map[string]interface{})

			switch data.Type {
			case WSEventTypeDomainCreate:
				createData, err := JSONStructToMap[WSCreateDomainData](dataAsMap)

				if err != nil {
					continue
				}

				now := time.Now().Unix()
				domain := Domain{
					Domain:      createData.Domain,
					Description: createData.Description,
					Category:    createData.Category,
					Target:      createData.Target,
					Added:       now,
					Checked:     now,
				}
				c.cache.domainIndex[domain.Domain] = domain
			case WSEventTypeDomainUpdate:
				updateData, err := JSONStructToMap[WSUpdateDomainData](dataAsMap)

				if err != nil {
					continue
				}

				currentDomain := c.cache.domainIndex[updateData.Domain]

				if updateData.Category != "" {
					currentDomain.Category = updateData.Category
				}
				if updateData.Description != "" {
					currentDomain.Description = updateData.Description
				}
				if updateData.Target != "" {
					currentDomain.Target = updateData.Target
				}
				currentDomain.Checked = updateData.Checked
			case WSEventTypeDomainDelete:
				deleteData, err := JSONStructToMap[WSDeleteDomainData](dataAsMap)

				if err != nil {
					continue
				}

				delete(c.cache.domainIndex, deleteData.Domain)
			case WSEventTypeURLCreate:
				createData, err := JSONStructToMap[WSCreateURLData](dataAsMap)

				if err != nil {
					continue
				}

				now := time.Now().Unix()
				url := URL{
					URL:         createData.URL,
					Description: createData.Description,
					Category:    createData.Category,
					Target:      createData.Target,
					Added:       now,
					Checked:     now,
				}

				c.cache.urlIndex[url.URL] = url
			case WSEventTypeURLUpdate:
				updateData, err := JSONStructToMap[WSUpdateURLData](dataAsMap)

				if err != nil {
					continue
				}

				currentURL := c.cache.urlIndex[updateData.URL]

				if updateData.Category != "" {
					currentURL.Category = updateData.Category
				}
				if updateData.Description != "" {
					currentURL.Description = updateData.Description
				}
				if updateData.Target != "" {
					currentURL.Target = updateData.Target
				}
				currentURL.Checked = updateData.Checked
				c.cache.urlIndex[currentURL.URL] = currentURL

			case WSEventTypeURLDelete:
				deleteData, err := JSONStructToMap[WSDeleteURLData](dataAsMap)

				if err != nil {
					continue
				}

				delete(c.cache.urlIndex, deleteData.URL)
			}

			c.cache.mx.Unlock()
		}
	}(c)
}

func (c *AutoSyncClient) StopAutoSync() {
	c.cacheTicker.Stop()
	c.sessionTicker.Stop()
	c.context.cancel()
}

func (c *AutoSyncClient) GetDomains() []Domain {
	c.cache.mx.RLock()
	defer c.cache.mx.RUnlock()

	values := make([]Domain, 0, len(c.cache.domainIndex))
	for _, d := range c.cache.domainIndex {
		values = append(values, d)
	}

	return values
}

func (c *AutoSyncClient) GetURLs() []URL {
	c.cache.mx.RLock()
	defer c.cache.mx.RUnlock()

	values := make([]URL, 0, len(c.cache.urlIndex))
	for _, u := range c.cache.urlIndex {
		values = append(values, u)
	}

	return values
}

func (c *AutoSyncClient) GetDomain(domain string) (*Domain, error) {
	c.cache.mx.RLock()
	defer c.cache.mx.RUnlock()

	d, ok := c.cache.domainIndex[domain]

	if !ok {
		return nil, fmt.Errorf("domain %s not found", domain)
	}

	return &d, nil
}

func (c *AutoSyncClient) GetURL(url string) (*URL, error) {
	c.cache.mx.RLock()
	defer c.cache.mx.RUnlock()

	u, ok := c.cache.urlIndex[url]

	if !ok {
		return nil, fmt.Errorf("url %s not found", url)
	}

	return &u, nil
}
