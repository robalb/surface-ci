package configfiles

import (
	"fmt"
	"github.com/robalb/tinyasm/pkg/pipeline"
	"github.com/robalb/tinyasm/pkg/validation"
	"gopkg.in/yaml.v3"
	"os"
)

type scopeFileData struct {
	Scope      pipeline.Surface `yaml:"scope"`
	Exclusions pipeline.Surface `yaml:"exclusions,omitempty"`
}

func parseScope(filePath string) (*scopeFileData, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read scope file at %s: %w", filePath, err)
	}

	config := scopeFileData{
		Scope: pipeline.Surface{
			Domains: []string{},
			IPs:     []string{},
			URLs:    []string{},
		},
		Exclusions: pipeline.Surface{
			Domains: []string{},
			IPs:     []string{},
			URLs:    []string{},
		},
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("Failed to parse scope file at %s: Invalid Syntax: %w", filePath, err)
	}

	// Check for empty scope
	s := config.Scope
	if len(s.Domains) == 0 && len(s.IPs) == 0 && len(s.URLs) == 0 {
		return nil, fmt.Errorf("Failed to parse scope file at %s: the scope cannot be emtpy", filePath)
	}

	err = validateSurface(&config.Scope)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse scope file at %s: In section 'scope': %w", filePath, err)
	}

	err = validateSurface(&config.Exclusions)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse scope file at %s: In section 'exclusions': %w", filePath, err)
	}

	return &config, nil
}

func validateSurface(s *pipeline.Surface) error {

	// Validate domains
	for i, domain := range s.Domains {
		if err := validation.ValidateDomain(domain); err != nil {
			return fmt.Errorf("Invalid domain at index %d: %w", i, err)
		}
	}

	// Validate IPs
	for i, ip := range s.IPs {
		if err := validation.ValidateIP(ip); err != nil {
			return fmt.Errorf("Invalid IP at index %d: %w", i, err)
		}
	}

	// Validate URLs
	for i, url := range s.URLs {
		if err := validation.ValidateURL(url); err != nil {
			return fmt.Errorf("Invalid url at index %d: %w", i, err)
		}
	}

	return nil
}
