package pipeline

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/net/publicsuffix"
)

// TrimSubdomains takes a slice of domains and returns a new slice
// where subdomains of lower hierarchies have been filtered out
// for example: [a.sub.example.com, sub.example.com] -> [sub.example.com]
func TrimSubdomains(domains []string) ([]string, error) {
	if len(domains) == 0 {
		return []string{}, nil
	}

	// Group domains by their effective TLD+1 (registered domain)
	domainGroups := make(map[string][]string)

	for _, domain := range domains {
		// Get the registered domain (eTLD+1) using lowercase for consistency
		etldPlusOne, err := publicsuffix.EffectiveTLDPlusOne(strings.ToLower(domain))
		if err != nil {
			return nil, fmt.Errorf("error processing domain %s: %v", domain, err)
		}

		domainGroups[etldPlusOne] = append(domainGroups[etldPlusOne], strings.ToLower(domain))
	}

	var result []string

	// Process each group separately
	for _, group := range domainGroups {
		// Sort by domain parts (longer domains first)
		sort.Slice(group, func(i, j int) bool {
			return len(strings.Split(group[i], ".")) > len(strings.Split(group[j], "."))
		})

		// Track which domains should be kept
		shouldKeep := make([]bool, len(group))
		for i := range shouldKeep {
			shouldKeep[i] = true
		}

		// Mark subdomains for removal - with case insensitive comparison
		for i := range group {
			if !shouldKeep[i] {
				continue
			}

			for j := i + 1; j < len(group); j++ {
				// Case insensitive suffix check
				if strings.HasSuffix(group[i], "."+group[j]) {
					shouldKeep[i] = false
					break
				}
			}
		}

		// Add domains that should be kept to the result
		for i, domain := range group {
			if shouldKeep[i] {
				result = append(result, domain)
			}
		}
	}

	return result, nil
}
