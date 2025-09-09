package pipeline

import "strings"

// SelectSubdomains takes a slice of source domains,
// a list of selector domains,
// and returns all source domains that are also
// subdomains of at least one selector domain
func SelectSubdomains(domains []string, selectorDomains []string) []string {
	var result []string
	
	for _, domain := range domains {
		for _, selector := range selectorDomains {
			// A domain is a subdomain if it equals the selector domain
			// or if it ends with "."+selector
			if domain == selector || strings.HasSuffix(domain, "."+selector) {
				result = append(result, domain)
				break // Once we find a match, no need to check other selectors
			}
		}
	}
	
	return result
}

