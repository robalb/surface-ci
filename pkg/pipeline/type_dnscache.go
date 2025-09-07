package pipeline

import (
	"sync"
)

type DNSCache struct {
	cache map[string][]string
	mutex sync.RWMutex
}

func NewDNSCache() *DNSCache {
	return &DNSCache{
		cache: make(map[string][]string),
	}
}

// Get returns cached IPs for a domain if they exist
func (c *DNSCache) Get(domain string) ([]string, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	ips, exists := c.cache[domain]
	return ips, exists
}

// Set caches the IPs for a domain
func (c *DNSCache) Set(domain string, ips []string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache[domain] = ips
}
