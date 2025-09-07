package pipeline

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestDnsxFilterWildcards(t *testing.T) {
	tests := []struct {
		name         string
		inputDomains []string
		wildcards    []string
		expected     []string
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
				"test.com",      // *.test.com is a wildcard
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
			wildcards: []string{},
			expected:  []string{},
		},
		// {
		// 	name: "All domains are wildcards",
		// 	inputDomains: []string{
		// 		"example.com",
		// 		"test.org",
		// 		"dev.com",
		// 	},
		// 	wildcards: []string{
		// 		"example.com",
		// 		"test.org",
		// 		"dev.com",
		// 	},
		// 	expected: []string{
		// 		"example.com",
		// 		"test.org",
		// 		"dev.com",
		// 	},
		// },
		// {
		// 	name: "Nested wildcards - child should be detected, parent ignored",
		// 	inputDomains: []string{
		// 		"example.com",
		// 		"a.example.com",
		// 		"b.a.example.com",
		// 		"c.a.example.com",
		// 		"b.example.com",
		// 	},
		// 	wildcards: []string{
		// 		"example.com",   // *.example.com is a wildcard
		// 		"a.example.com", // *.a.example.com is also a wildcard
		// 	},
		// 	expected: []string{
		// 		"example.com",
		// 		"a.example.com",
		// 	},
		// },
		// {
		// 	name: "Multiple TLDs",
		// 	inputDomains: []string{
		// 		"example.com",
		// 		"sub.example.com",
		// 		"example.org",
		// 		"sub.example.org",
		// 		"example.net",
		// 		"sub.example.net",
		// 		"example.co.uk",
		// 		"sub.example.co.uk",
		// 	},
		// 	wildcards: []string{
		// 		"example.com",
		// 		"example.net",
		// 	},
		// 	expected: []string{
		// 		"example.com",
		// 		"example.net",
		// 	},
		// },
		// {
		// 	name: "Complex effective TLDs",
		// 	inputDomains: []string{
		// 		"example.co.uk",
		// 		"sub.example.co.uk",
		// 		"example.com.au",
		// 		"sub.example.com.au",
		// 		"example.github.io",
		// 		"sub.example.github.io",
		// 	},
		// 	wildcards: []string{
		// 		"example.co.uk",
		// 		"example.github.io",
		// 	},
		// 	expected: []string{
		// 		"example.co.uk",
		// 		"example.github.io",
		// 	},
		// },
		// {
		// 	name: "Wildcards at different levels",
		// 	inputDomains: []string{
		// 		"example.com",
		// 		"a.example.com",
		// 		"b.example.com",
		// 		"x.a.example.com",
		// 		"y.a.example.com",
		// 		"z.b.example.com",
		// 	},
		// 	wildcards: []string{
		// 		"b.example.com", // *.b.example.com is a wildcard
		// 	},
		// 	expected: []string{
		// 		"b.example.com",
		// 	},
		// },
		// {
		// 	name: "Partial wildcards in group",
		// 	inputDomains: []string{
		// 		"a.example.com",
		// 		"b.example.com",
		// 		"c.example.com",
		// 		"sub.a.example.com",
		// 		"sub.b.example.com",
		// 		"sub.c.example.com",
		// 	},
		// 	wildcards: []string{
		// 		"a.example.com",
		// 		"c.example.com",
		// 	},
		// 	expected: []string{
		// 		"a.example.com",
		// 		"c.example.com",
		// 	},
		// },
		// {
		// 	name: "Sibling wildcards",
		// 	inputDomains: []string{
		// 		"a.example.com",
		// 		"b.example.com",
		// 		"x.a.example.com",
		// 		"y.a.example.com",
		// 		"x.b.example.com",
		// 		"y.b.example.com",
		// 	},
		// 	wildcards: []string{
		// 		"a.example.com",
		// 		"b.example.com",
		// 	},
		// 	expected: []string{
		// 		"a.example.com",
		// 		"b.example.com",
		// 	},
		// },
		// {
		// 	name: "Original example case",
		// 	inputDomains: []string{
		// 		"a.example.com",
		// 		"b.example.com",
		// 		"b.a.example.com",
		// 		"c.a.example.com",
		// 		"a.test.com",
		// 		"b.test.com",
		// 	},
		// 	wildcards: []string{
		// 		"a.example.com",
		// 		"test.com",
		// 	},
		// 	expected: []string{
		// 		"a.example.com",
		// 		"a.test.com",
		// 		"b.test.com",
		// 	},
		// },
		// {
		// 	name: "Wildcard with no subdomains in input",
		// 	inputDomains: []string{
		// 		"example.com",
		// 		"test.com",
		// 		"domain.org",
		// 	},
		// 	wildcards: []string{
		// 		"example.com",
		// 	},
		// 	expected: []string{
		// 		"example.com",
		// 	},
		// },
		// {
		// 	name: "Very deep domain hierarchy",
		// 	inputDomains: []string{
		// 		"a.b.c.d.e.f.example.com",
		// 		"x.y.z.example.com",
		// 		"example.com",
		// 	},
		// 	wildcards: []string{
		// 		"c.d.e.f.example.com",
		// 		"z.example.com",
		// 	},
		// 	expected: []string{
		// 		"c.d.e.f.example.com",
		// 	},
		// },
		// {
		// 	name: "Non-existent parent domains",
		// 	inputDomains: []string{
		// 		"a.example.com",
		// 		"b.c.example.com",
		// 	},
		// 	wildcards: []string{
		// 		"example.com", // Parent domain with wildcard, but not in input list
		// 	},
		// 	expected: []string{
		// 		// Empty because example.com isn't in the input domains
		// 	},
		// },
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock lookup function
			mockLookup := func(domain string) ([]string, error) {
				// For any test domain (containing a UUID), check if it's under a wildcard
				for _, wildcard := range tt.wildcards {
					if strings.HasSuffix(domain, "."+wildcard) {
						// This is under a wildcard, so it should resolve
						return []string{"192.0.2.1"}, nil // Return a dummy IP
					}
				}
				
				// Not under a wildcard
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
