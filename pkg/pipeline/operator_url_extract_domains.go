package pipeline

import (
	"net/url"
	"strings"
)

// ExtractDomains takes a slice of URLs and returns a slice of unique domains
func URLExtractDomains(urls []string) []string {
	// Use a map to track unique domains
	uniqueDomains := make(map[string]struct{})

	for _, rawURL := range urls {
		// If URL doesn't have a scheme, add one to make parsing work
		parsableURL := rawURL
		if !strings.HasPrefix(strings.ToLower(parsableURL), "http://") &&
			!strings.HasPrefix(strings.ToLower(parsableURL), "https://") {
			parsableURL = "https://" + parsableURL
		}

		parsedURL, err := url.Parse(parsableURL)
		if err != nil {
			// Skip invalid URLs
			continue
		}

		host := parsedURL.Hostname() // This removes any port number
		if host == "" {
			// skip empty hosts
			continue
		}

		uniqueDomains[host] = struct{}{}
	}

	// Convert the map keys to a slice
	domains := make([]string, 0, len(uniqueDomains))
	for domain := range uniqueDomains {
		domains = append(domains, domain)
	}

	return domains
}
