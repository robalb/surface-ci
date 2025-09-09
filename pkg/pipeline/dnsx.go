package pipeline

import (
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/projectdiscovery/dnsx/libs/dnsx"
	"golang.org/x/net/publicsuffix"
)

// DNSLookupFunc defines the signature for a DNS lookup functions
type DNSLookupFunc func(domain string) ([]string, error)

// FilterActiveDomains takes a list of domains and returns those with valid DNS records
// If a DNSCache is provided, it will check the cache before making DNS queries
// and update the cache with successful resolutions
func DnsxFilterActive(domains []string, cache *DNSCache) []string {
	// First, check all domains against the cache
	var validDomains []string
	var domainsToResolve []string

	for _, domain := range domains {
		ips, found := cache.Get(domain)
		if found {
			if len(ips) > 0 {
				validDomains = append(validDomains, domain)
			}
		} else {
			domainsToResolve = append(domainsToResolve, domain)
		}
	}

	// If there are no domains to resolve, return the valid ones
	if len(domainsToResolve) == 0 {
		return validDomains
	}

	// Create DNS Resolver with default options
	dnsClient, err := dnsx.New(dnsx.DefaultOptions)
	if err != nil {
		return validDomains
	}

	// Resolve all the domains not found in cache
	for _, domain := range domainsToResolve {
		// Use Lookup to get IP addresses
		ips, err := dnsClient.Lookup(domain)
		if err != nil || len(ips) == 0 {
			// Store empty result to prevent future lookups
			cache.Set(domain, []string{})
			continue
		}

		// Store the results in cache
		cache.Set(domain, ips)
		validDomains = append(validDomains, domain)
	}

	return validDomains
}

// DnsxFilterWildcards takes a list of domains and returns those that
// are the root of a wildcard domain.
func DnsxFilterWildcards(domains []string, cache *DNSCache) []string {
	dnsClient, err := dnsx.New(dnsx.DefaultOptions)
	if err != nil {
		return nil
	}
	var defaultDNSLookup DNSLookupFunc = func(domain string) ([]string, error) {
		return dnsClient.Lookup(domain)
	}
	return dnsxFilterWildcards(domains, cache, defaultDNSLookup)
}

// TODO: set a depth limit
// wildcard discovery algorythm.
// in a situation were the following wildcards exist:
// *.a.example.com
// *.test.com
//
// and we are given the following domains list:
// a.example.com
// b.example.com
// b.a.example.com
// c.a.example.com
// a.test.com
// b.test.com
//
// we expect to receive this:
// a.example.com
// a.test.com
// b.test.com
//
// note how in the case of test.com the real wildcard is *.test.com, but we could not
// detect that becasuse test.com is not in scope.
// If we only have a subdomain in scope, the parent domain
// is not intended to be in scope, we are not allowed to go there.
func dnsxFilterWildcards(domains []string, cache *DNSCache, dnsLookup DNSLookupFunc) []string {

	// Group domains by base domain (effective TLD+1)
	domainGroups := make(map[string][]string)

	wildcardDomains := []string{}

	for _, domain := range domains {
		baseDomain, err := publicsuffix.EffectiveTLDPlusOne(domain)
		if err != nil {
			continue
		}
		domainGroups[baseDomain] = append(domainGroups[baseDomain], domain)
	}

	// Process each group
	for _, group := range domainGroups {
		// Sort by depth (number of dots), ascending (process parents first)
		sort.Slice(group, func(i, j int) bool {
			return countDots(group[i]) < countDots(group[j])
		})

		// Check each domain
		for _, domain := range group {
			// Skip if this domain is a child of a known wildcard
			shouldSkip := false
			for _, parent := range wildcardDomains {
				if isSubdomain(domain, parent) {
					shouldSkip = true
					break
				}
			}

			if shouldSkip {
				continue
			}

			// Generate a random subdomain to test
			randomPart := strings.ReplaceAll(uuid.New().String(), "-", "")
			testDomain := fmt.Sprintf("%s.%s", randomPart, domain)

			// Check if the random subdomain resolves
			ips, found := cache.Get(testDomain)
			if !found {
				var err error
				ips, err = dnsLookup(testDomain)
				if err != nil {
					ips = []string{}
				}
				cache.Set(testDomain, ips) // Cache the result
			}

			// If the random subdomain resolves, we've found a wildcard
			if len(ips) > 0 {
				wildcardDomains = append(wildcardDomains, domain)
				fmt.Printf("  - resolved! %v \n", wildcardDomains)
			}
		}
	}

	return wildcardDomains
}

// countDots counts the number of dots in a domain name
func countDots(domain string) int {
	return strings.Count(domain, ".")
}

// isSubdomain checks if child is a subdomain of parent
func isSubdomain(child, parent string) bool {
	// Check if child ends with .parent
	suffix := "." + parent
	return strings.HasSuffix(child, suffix)
}
