package pipeline

import "testing"

func TestSelectSubdomains(t *testing.T) {
	tests := []struct {
		name            string
		domains         []string
		selectorDomains []string
		expected        []string
	}{
		{
			name:            "empty inputs",
			domains:         []string{},
			selectorDomains: []string{},
			expected:        []string{},
		},
		{
			name:            "no matches",
			domains:         []string{"example.com", "test.org"},
			selectorDomains: []string{"google.com"},
			expected:        []string{},
		},
		{
			name:            "exact matches",
			domains:         []string{"example.com", "google.com", "test.org"},
			selectorDomains: []string{"google.com"},
			expected:        []string{"google.com"},
		},
		{
			name:            "subdomain matches",
			domains:         []string{"mail.google.com", "docs.google.com", "example.com"},
			selectorDomains: []string{"google.com"},
			expected:        []string{"mail.google.com", "docs.google.com"},
		},
		{
			name:            "multiple selectors",
			domains:         []string{"mail.google.com", "example.com", "blog.example.com"},
			selectorDomains: []string{"google.com", "example.com"},
			expected:        []string{"mail.google.com", "example.com", "blog.example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SelectSubdomains(tt.domains, tt.selectorDomains)
			
			// Simple equality check - doesn't care about order
			if len(result) != len(tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, result)
			}
			
			// Check that all expected elements are in the result
			for _, e := range tt.expected {
				found := false
				for _, r := range result {
					if e == r {
						found = true
						break
					}
				}
				if !found {
					t.Fatalf("expected %s to be in result %v", e, result)
				}
			}
		})
	}
}
