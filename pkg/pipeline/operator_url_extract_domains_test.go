package pipeline

import (
	"reflect"
	"sort"
	"testing"
)

func TestExtractDomains(t *testing.T) {
	tests := []struct {
		name     string
		urls     []string
		expected []string
	}{
		{
			name:     "Basic URLs",
			urls:     []string{"https://example.com", "http://test.org"},
			expected: []string{"example.com", "test.org"},
		},
		{
			name:     "URLs with subdomains",
			urls:     []string{"https://www.example.com", "http://blog.example.com", "api.test.org"},
			expected: []string{"www.example.com", "blog.example.com", "api.test.org"},
		},
		{
			name:     "URLs with paths and query parameters",
			urls:     []string{"https://example.com/path", "http://test.org/page?q=value#fragment"},
			expected: []string{"example.com", "test.org"},
		},
		{
			name:     "URLs with port numbers",
			urls:     []string{"https://example.com:443", "http://test.org:8080/path"},
			expected: []string{"example.com", "test.org"},
		},
		{
			name:     "URLs without schemes",
			urls:     []string{"example.com", "test.org/path"},
			expected: []string{"example.com", "test.org"},
		},
		{
			name:     "Mixed valid and invalid URLs",
			urls:     []string{"https://example.com", "@", "test.org"},
			expected: []string{"example.com", "test.org"},
		},
		{
			name:     "Duplicate domains",
			urls:     []string{"https://example.com/page1", "http://example.com/page2", "example.com/page3"},
			expected: []string{"example.com"},
		},
		{
			name:     "Empty input",
			urls:     []string{},
			expected: []string{},
		},
		{
			name:     "URLs with international domains",
			urls:     []string{"https://例子.测试", "http://例子.测试/path"},
			expected: []string{"例子.测试"},
		},
		{
			name:     "URLs with IPv6 addresses",
			urls:     []string{"http://[2001:db8:85a3:8d3:1319:8a2e:370:7348]", "https://[::1]:8080"},
			expected: []string{"2001:db8:85a3:8d3:1319:8a2e:370:7348", "::1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := URLExtractDomains(tt.urls)
			
			// Sort both slices for consistent comparison
			sort.Strings(result)
			sort.Strings(tt.expected)
			
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExtractDomains() = %v, want %v", result, tt.expected)
			}
		})
	}
}
