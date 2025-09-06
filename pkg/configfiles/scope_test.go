package configfiles

import (
	"reflect"
	"strings"
	"testing"

	"github.com/robalb/tinyasm/pkg/surface"
)

func TestBadConfigurations(t *testing.T) {
	tests := []struct {
		name        string
		filePath    string
		errContains string
	}{
		{"empty_domain", "testdata/scope/empty_domain.yaml", "Domain cannot be empty"},
		{"empty_file", "testdata/scope/empty_file.yaml", "scope cannot be emtpy"},
		{"missing_file", "testdata/scope/DO_NOT_CREATE_ME", "no such file or directory"},
		{"empty_section", "testdata/scope/empty_section.yaml", "scope cannot be emtpy"},
		{"empty_section2", "testdata/scope/empty_section.yaml", "scope cannot be emtpy"},
		{"malformed_yaml", "testdata/scope/malformed_yaml.yaml", "Invalid Syntax"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseScope(tt.filePath)

			// t.Fatalf("error for %s, %v", tt.name, err)
			if err == nil {
				t.Fatalf("Expected error for %s, got nil", tt.name)
			}
			if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
				t.Fatalf("Error message doesn't contain %q: %v", tt.errContains, err)
			}
		})
	}
}

func TestValidConfigurations(t *testing.T) {
	tests := []struct {
		name         string
		filePath     string
		expectedData *scopeFileData
	}{
		{
			name:     "basic_valid_config",
			filePath: "testdata/scope/valid_basic.yaml",
			expectedData: &scopeFileData{
				Scope: surface.Surface{
					Domains: []string{"example.com", "test-domain.org"},
					IPs:     []string{"192.168.1.1", "10.0.0.0/24"},
					URLs:    []string{"https://example.com/api", "example.org/endpoint"},
				},
			},
		},
		{
			name:     "valid_with_comments",
			filePath: "testdata/scope/valid_with_comments.yaml",
			expectedData: &scopeFileData{
				Scope: surface.Surface{
					Domains: []string{"example.com"},
					IPs:     []string{},
					URLs:    []string{"https://example.com/api", "http://test.com"},
				},
			},
		},
		{
			name:     "domains_only",
			filePath: "testdata/scope/valid_domains_only.yaml",
			expectedData: &scopeFileData{
				Scope: surface.Surface{
					Domains: []string{"example.com", "sub.example.org", "test-site.co.uk"},
					IPs:     []string{},
					URLs:    []string{},
				},
			},
		},
		{
			name:     "with_exclusions",
			filePath: "testdata/scope/valid_with_exclusions.yaml",
			expectedData: &scopeFileData{
				Scope: surface.Surface{
					Domains: []string{"example.com"},
					IPs:     []string{"192.168.0.0/16"},
					URLs:    []string{"https://example.com/api"},
				},
				Exclusions: surface.Surface{
					Domains: []string{"admin.example.com", "internal.example.com"},
					IPs:     []string{"192.168.1.100", "192.168.2.0/24"},
					URLs:    []string{"https://example.com/admin", "example.com/internal"},
				},
			},
		},
		{
			name:     "ipv6_addresses",
			filePath: "testdata/scope/valid_ipv6.yaml",
			expectedData: &scopeFileData{
				Scope: surface.Surface{
					Domains: []string{"example.com"},
					IPs:     []string{"2001:db8::1", "2001:db8::/64", "::1"},
					URLs:    []string{"https://example.com/api"},
				},
			},
		},
		{
			name:     "minimal_valid",
			filePath: "testdata/scope/valid_minimal.yaml",
			expectedData: &scopeFileData{
				Scope: surface.Surface{
					Domains: []string{"example.com"},
					IPs:     []string{},
					URLs:    []string{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := parseScope(tt.filePath)
			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}

			// Compare expected and actual configurations
			// Check Scope section
			if !reflect.DeepEqual(config.Scope.Domains, tt.expectedData.Scope.Domains) {
				t.Errorf("Domains mismatch.\nExpected: %v\nGot: %v",
					tt.expectedData.Scope.Domains, config.Scope.Domains)
			}
			if !reflect.DeepEqual(config.Scope.IPs, tt.expectedData.Scope.IPs) {
				t.Errorf("IPs mismatch.\nExpected: %v\nGot: %v",
					tt.expectedData.Scope.IPs, config.Scope.IPs)
			}
			if !reflect.DeepEqual(config.Scope.URLs, tt.expectedData.Scope.URLs) {
				t.Errorf("URLs mismatch.\nExpected: %v\nGot: %v",
					tt.expectedData.Scope.URLs, config.Scope.URLs)
			}

			// Check Exclusions section if it exists
			if tt.expectedData.Exclusions.Domains != nil {
				if !reflect.DeepEqual(config.Exclusions.Domains, tt.expectedData.Exclusions.Domains) {
					t.Errorf("Exclusion Domains mismatch.\nExpected: %v\nGot: %v",
						tt.expectedData.Exclusions.Domains, config.Exclusions.Domains)
				}
				if !reflect.DeepEqual(config.Exclusions.IPs, tt.expectedData.Exclusions.IPs) {
					t.Errorf("Exclusion IPs mismatch.\nExpected: %v\nGot: %v",
						tt.expectedData.Exclusions.IPs, config.Exclusions.IPs)
				}
				if !reflect.DeepEqual(config.Exclusions.URLs, tt.expectedData.Exclusions.URLs) {
					t.Errorf("Exclusion URLs mismatch.\nExpected: %v\nGot: %v",
						tt.expectedData.Exclusions.URLs, config.Exclusions.URLs)
				}
			}
		})
	}
}
