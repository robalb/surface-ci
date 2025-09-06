package datafiles

import (
	"fmt"
	"os"

	"github.com/robalb/tinyasm/pkg/surface"
	"github.com/robalb/tinyasm/pkg/validation"
	"gopkg.in/yaml.v3"
)

type knownSurfaceFileData struct {
    KnownSurface surface.Surface `yaml:"known_surface"`
}

func parseKnownSurface(filePath string) (*knownSurfaceFileData, error){
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("Failed to read scope file at %s: %w", filePath, err)
    }

    surface := knownSurfaceFileData{
        KnownSurface: surface.Surface{
            Domains: []string{},
            IPs:     []string{},
            URLs:    []string{},
        },
    }
    if err := yaml.Unmarshal(data, &surface); err != nil {
		return nil, fmt.Errorf("Failed to parse known-surface file at %s: Invalid Syntax: %w", filePath, err)
    }

	s := surface.KnownSurface

	// Validate domains
	for i, domain := range s.Domains {
		if err := validation.ValidateDomain(domain); err != nil {
			return nil, fmt.Errorf("Invalid domain at index %d: %w", i, err)
		}
	}

	// Validate IPs
	for i, ip := range s.IPs {
		if err := validation.ValidateIP(ip); err != nil {
			return nil, fmt.Errorf("Invalid IP at index %d: %w", i, err)
		}
	}

	// Validate URLs
	for i, url := range s.URLs {
		if err := validation.ValidateURL(url); err != nil {
			return nil, fmt.Errorf("Invalid url at index %d: %w", i, err)
		}
	}


	return &surface, nil

}



