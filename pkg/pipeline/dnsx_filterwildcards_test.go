package pipeline

import (
	"slices"
	"reflect"
	"sort"
	"testing"
)

func TestDnsxFilterWildcards(t *testing.T) {
	tests := []struct {
		name         string
		inputDomains []string
		expected     []string

		wildcards    []string
	    unregistered []string
	}{
		{
			name: "Basic wildcard detection",
			inputDomains: []string{
				"a.example.com",
				"b.example.com",
				"b.a.example.com",
				"c.a.example.com",
				"a.test.com",
				"b.test.com",
			},
			wildcards: []string{
				"a.example.com", // *.a.example.com is a wildcard
				"test.com",    // *.test.com is a wildcard
			},
			unregistered: []string{
				"example.com",
			},
			expected: []string{
				"a.example.com",
				"a.test.com",
				"b.test.com",
			},
		},
		{
			name: "No wildcards",
			inputDomains: []string{
				"example.com",
				"sub.example.com",
				"another.sub.example.com",
				"test.org",
				"dev.test.org",
			},
			unregistered: []string{},
			wildcards: []string{},
			expected:  []string{},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock dns lookup function
			mockLookup := func(domain string) ([]string, error) {
				// if its under a wildcard, resolve a dummy ip.
				for _, wildcard := range tt.wildcards {
					if isSubdomain(domain, wildcard) {
						return []string{"192.0.1.1"}, nil
					}
				}

				// if it's in the unregistered list, return nothing
				if slices.Contains(tt.unregistered, domain) {
					return []string{}, nil 
				}

				// if it's in the input domains, return a dummy ip.
				if slices.Contains(tt.inputDomains, domain) {
					return []string{"192.0.3.1"}, nil
				}
				
				// anything else falls here
				return []string{}, nil
			}
			
			cache := NewDNSCache()
			
			// Run the function with our mock lookup
			got := dnsxFilterWildcards(tt.inputDomains, cache, mockLookup)
			
			// Sort both slices to ensure order doesn't matter
			sort.Strings(got)
			sort.Strings(tt.expected)
			
			// Compare results
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("DnsxFindWildcards() = %v, want %v", got, tt.expected)
			}
		})
	}
}





