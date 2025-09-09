package pipeline

import (
	"reflect"
	"sort"
	"testing"
)

func TestTrimSubdomains(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
		wantErr  bool
	}{
		{
			name:     "Basic example",
			input:    []string{"aa.sub.example.com", "bb.sub.example.com", "sub.example.com", "test.com"},
			expected: []string{"sub.example.com", "test.com"},
			wantErr:  false,
		},
		{
			name:     "Empty input",
			input:    []string{},
			expected: []string{},
			wantErr:  false,
		},
		{
			name:     "Single domain",
			input:    []string{"example.com"},
			expected: []string{"example.com"},
			wantErr:  false,
		},
		{
			name:     "Different TLDs",
			input:    []string{"example.com", "example.org", "example.net"},
			expected: []string{"example.com", "example.org", "example.net"},
			wantErr:  false,
		},
		{
			name:     "Multiple hierarchy levels",
			input:    []string{"deep.deeper.sub.example.com", "deeper.sub.example.com", "sub.example.com", "example.com"},
			expected: []string{"example.com"},
			wantErr:  false,
		},
		{
			name:     "Effective TLDs",
			input:    []string{"sub.example.co.uk", "example.co.uk", "test.co.uk"},
			expected: []string{"example.co.uk", "test.co.uk"},
			wantErr:  false,
		},
		{
			name:     "Effective TLDs mixed cases",
			input:    []string{"sub.example.co.uk", "example.CO.UK", "test.co.uk"},
			expected: []string{"example.co.uk", "test.co.uk"},
			wantErr:  false,
		},
		// {
		// 	name:     "Mixed case",
		// 	input:    []string{"SUB.Example.COM", "sub.example.com"},
		// 	expected: []string{"sub.example.com"}, // Assuming case-insensitive comparison
		// 	wantErr:  false,
		// },
		{
			name:     "With invalid domain",
			input:    []string{"example.com", "invalid..domain"},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := TrimSubdomains(tt.input)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterSubdomains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we expect an error, don't check the result
			if tt.wantErr {
				return
			}

			// Sort results for consistent comparison
			sort.Strings(result)
			sort.Strings(tt.expected)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FilterSubdomains() = %v, want %v", result, tt.expected)
			}
		})
	}
}
