package pipeline

import (
	"net"
	"net/url"
	"strings"
)

// URLExtractIPs takes a slice of URLs and returns a slice of unique IP addresses
func URLExtractIPs(urls []string) []string {
	// Use a map to track unique IPs
	uniqueIPs := make(map[string]struct{})

	for _, rawURL := range urls {
		// If URL doesn't have a scheme, add one to make parsing work
		parsableURL := rawURL
		if !strings.HasPrefix(strings.ToLower(parsableURL), "http://") &&
			!strings.HasPrefix(strings.ToLower(parsableURL), "https://") {
			parsableURL = "https://" + parsableURL
		}

		// Parse the URL
		parsedURL, err := url.Parse(parsableURL)
		if err != nil {
			// Skip invalid URLs
			continue
		}

		// Extract the host part
		host := parsedURL.Hostname() // This removes any port number

		// Skip empty hosts
		if host == "" {
			continue
		}

		// Check if the host is an IPv6 address
		// net.ParseIP can handle IPv6 addresses without brackets
		ip := net.ParseIP(host)
		if ip != nil {
			// It's a valid IP (either IPv4 or IPv6)
			uniqueIPs[host] = struct{}{}
			continue
		}

		// Check if it's an IPv6 in brackets (which might not be properly parsed by Hostname())
		if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
			ipv6 := host[1 : len(host)-1]
			if ip := net.ParseIP(ipv6); ip != nil {
				uniqueIPs[ipv6] = struct{}{}
			}
		}
	}

	// Convert the map keys to a slice
	ips := make([]string, 0, len(uniqueIPs))
	for ip := range uniqueIPs {
		ips = append(ips, ip)
	}

	return ips
}
