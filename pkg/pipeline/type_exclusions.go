package pipeline

import "strings"

type Exclusions struct {
	Domains map[string]struct{}
	IPs     map[string]struct{}
	URLs    map[string]struct{}
}

func normalize_domain(domain string) string {
	return strings.ToLower(domain)
}

func normalize_ip(ip string) string {
	return ip // For now, no normalization needed for IPs
}

func normalize_url(url string) string {
	return strings.ToLower(url)
}

// MakeExclusion initializes a new Exclusions struct
func MakeExclusion() Exclusions {
	return Exclusions{
		Domains: make(map[string]struct{}),
		IPs:     make(map[string]struct{}),
		URLs:    make(map[string]struct{}),
	}
}

// Insert adds all elements from the given Surface to the Exclusions
func (e *Exclusions) Insert(s *Surface) {
	for _, domain := range s.Domains {
		e.Domains[normalize_domain(domain)] = struct{}{}
	}

	for _, ip := range s.IPs {
		e.IPs[normalize_ip(ip)] = struct{}{}
	}

	for _, url := range s.URLs {
		e.URLs[normalize_url(url)] = struct{}{}
	}
}

// Contains_domain checks if a domain is in the exclusions
func (e *Exclusions) Contains_domain(domain string) bool {
	_, exists := e.Domains[normalize_domain(domain)]

	//TODO: if domain is not a TLD, check if its parent
    //      is in the exclusion list.
	return exists
}

// Contains_ip checks if an IP is in the exclusions
func (e *Exclusions) Contains_ip(ip string) bool {
	_, exists := e.IPs[normalize_ip(ip)]
	return exists
}

// Contains_url checks if a URL is in the exclusions
func (e *Exclusions) Contains_url(url string) bool {
	_, exists := e.URLs[normalize_url(url)]
	return exists
}

// Contains checks if a string is in any of the exclusion lists
func (e *Exclusions) Contains(s string) bool {
	if e.Contains_domain(s) {
		return true
	}

	if e.Contains_ip(s) {
		return true
	}

	if e.Contains_url(s) {
		return true
	}

	return false
}
