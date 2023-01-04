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
	mx sync.RWMutex
	// Store as map to pointer to avoid copying every domain
	index   map[string]*Domain
	domains []Domain
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

	c.cache.domains = *domains
	c.cache.index = map[string]*Domain{}

	// Index domains in map for faster lookup
	for i, domain := range c.cache.domains {
		c.cache.index[domain.Domain] = &c.cache.domains[i]
	}

	return nil
}

func (c *AutoSyncClient) StartAutoSync(interval time.Duration) {
	context, cancel := context.WithCancel(context.Background())
	c.context.ctx = context
	c.context.cancel = cancel

	c.cacheTicker = time.NewTicker(interval)
	// Session tokens expire after an hour, refresh 15 minutes early just in case
	c.sessionTicker = time.NewTicker(time.Second * 10)

	// Generate Session Token
	// The client should already be able to successfully create a token from initialization
	token, _ := c.raw.CreateSessionToken()
	c.raw.UpdateSessionToken(*token)

	// Initial Sync
	c.ForceSync()

	// Start automatically syncing domains
	go func(client *AutoSyncClient) {
		for {
			select {
			case <-c.cacheTicker.C:
				c.ForceSync()
				fmt.Println("synced!")
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
				fmt.Println("old token: ", token)
				token, _ := c.raw.CreateSessionToken()
				c.raw.UpdateSessionToken(*token)
				fmt.Println("new token: ", token)
			case <-c.context.ctx.Done():
				return
			}
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

	return c.cache.domains
}

func (c *AutoSyncClient) GetDomain(domain string) (*Domain, error) {
	c.cache.mx.RLock()
	defer c.cache.mx.RUnlock()

	d, ok := c.cache.index[domain]

	if !ok {
		return nil, fmt.Errorf("domain %s not found", domain)
	}

	return d, nil
}
