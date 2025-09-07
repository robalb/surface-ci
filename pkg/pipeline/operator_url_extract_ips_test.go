package pipeline

import (
	"reflect"
	"sort"
	"testing"
)

func TestExtractIPs(t *testing.T) {
	tests := []struct {
		name     string
		urls     []string
		expected []string
	}{
		{
			name:     "Basic IPv4 URLs",
			urls:     []string{"https://192.168.1.1", "http://10.0.0.1"},
			expected: []string{"192.168.1.1", "10.0.0.1"},
		},
		{
			name:     "IPv4 URLs with paths and query parameters",
			urls:     []string{"https://192.168.1.1/path", "http://10.0.0.1/page?q=value#fragment"},
			expected: []string{"192.168.1.1", "10.0.0.1"},
		},
		{
			name:     "IPv4 URLs with port numbers",
			urls:     []string{"https://192.168.1.1:443", "http://10.0.0.1:8080/path"},
			expected: []string{"192.168.1.1", "10.0.0.1"},
		},
		{
			name:     "IPv4 without schemes",
			urls:     []string{"192.168.1.1", "10.0.0.1/path"},
			expected: []string{"192.168.1.1", "10.0.0.1"},
		},
		{
			name:     "IPv6 URLs",
			urls:     []string{"https://[2001:db8::1]", "http://[::1]"},
			expected: []string{"2001:db8::1", "::1"},
		},
		{
			name:     "IPv6 URLs with ports and paths",
			urls:     []string{"https://[2001:db8::1]:443/path", "http://[::1]:8080"},
			expected: []string{"2001:db8::1", "::1"},
		},
		{
			name:     "Mixed IPv4 and IPv6",
			urls:     []string{"https://192.168.1.1", "http://[2001:db8::1]"},
			expected: []string{"192.168.1.1", "2001:db8::1"},
		},
		{
			name:     "Mixed IPs and domains",
			urls:     []string{"https://192.168.1.1", "https://example.com", "http://[2001:db8::1]"},
			expected: []string{"192.168.1.1", "2001:db8::1"},
		},
		{
			name:     "Duplicate IPs",
			urls:     []string{"https://192.168.1.1/page1", "http://192.168.1.1/page2", "192.168.1.1/page3"},
			expected: []string{"192.168.1.1"},
		},
		{
			name:     "Empty input",
			urls:     []string{},
			expected: []string{},
		},
		{
			name:     "No IPs in URLs",
			urls:     []string{"https://example.com", "http://test.org"},
			expected: []string{},
		},
		{
			name:     "Invalid URLs",
			urls:     []string{"invalid:/192.168.1.1", "http://[invalid-ipv6]"},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractIPs(tt.urls)
			
			// Sort both slices for consistent comparison
			sort.Strings(result)
			sort.Strings(tt.expected)
			
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExtractIPs() = %v, want %v", result, tt.expected)
			}
		})
	}
}
