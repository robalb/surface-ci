package validation

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
)

func ValidateDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("Domain cannot be empty")
	}

	// Remove leading/trailing whitespace
	domain = strings.TrimSpace(domain)

	// Check length
	if len(domain) > 253 {
		return fmt.Errorf("domain '%s' exceeds maximum length of 253 characters", domain)
	}

	if strings.Contains(domain, "://") {
		return fmt.Errorf("domain '%s' contains an invalid prefix. This is probably an URL, not a domain.", domain)
	}

	// Validate domain format
	var domainRegex = regexp.MustCompile(`^[a-zA-Z0-9\-]+\.[a-zA-Z0-9\-\.]+$`)
	if !domainRegex.MatchString(domain) {
		return fmt.Errorf("domain '%s' has invalid format", domain)
	}

	return nil
}

func ValidateIP(ip string) error {
	if ip == "" {
		return fmt.Errorf("IP cannot be empty")
	}

	// Remove leading/trailing whitespace
	ip = strings.TrimSpace(ip)

	// Check if it's a CIDR notation
	if strings.Contains(ip, "/") {
		_, _, err := net.ParseCIDR(ip)
		if err != nil {
			return fmt.Errorf("invalid CIDR notation '%s'", ip)
		}
		return nil
	}

	// Check if it's a valid IP address
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("invalid IP address '%s'", ip)
	}

	return nil
}

func ValidateURL(endpoint string) error {
	if endpoint == "" {
		return fmt.Errorf("endpoint cannot be empty")
	}

	// Remove leading/trailing whitespace
	endpoint = strings.TrimSpace(endpoint)

	// If it doesn't start with a scheme, add one for validation
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}

	// Parse as URL
	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("invalid endpoint URL '%s': %w", endpoint, err)
	}

	// Check that it has a host
	if u.Host == "" {
		return fmt.Errorf("endpoint '%s' must have a host", endpoint)
	}

	return nil
}
