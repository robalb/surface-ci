// Package datafiles provides an abstraction layer between internal data structures and external data-storage files
package datafiles

import (
	"fmt"
	"github.com/robalb/tinyasm/pkg/surface"
	"os"
	"path"
)

var (
	knownSurfaceFileName = "discovered-surface.yaml"
	knownIssuesFileName  = "discovered-issues.yaml"
	datafileHeader       = "## This is a program-generated data file. Do not edit. ##"
)

type DataFiles struct {
	KnownSurface surface.Surface
	// knownIssues Issues TODO
}

func New(dataFolder string) (d *DataFiles, fileMissing bool, err error) {
	knownSurfaceFilePath := path.Join(dataFolder, knownSurfaceFileName)

	// check if the directory exists
	if _, err := os.Stat(dataFolder); os.IsNotExist(err) {
		return d, true, fmt.Errorf("data directory does not exist: %s", dataFolder)
	}

	fileMissing, err = initFile(knownSurfaceFilePath)
	if err != nil {
		return
	}

	knownSurfaceData, err := parseKnownSurface(knownSurfaceFilePath)
	if err != nil {
		return
	}

	d = &DataFiles{
		knownSurfaceData.KnownSurface,
	}
	return
}

func (d *DataFiles) Summary() string {
	return fmt.Sprintf(
		"Known surface elements discovered in the past: {Domains[%d], IPs[%d], Endpoints[%d]}",
		len(d.KnownSurface.Domains),
		len(d.KnownSurface.IPs),
		len(d.KnownSurface.URLs),
	)
}

func initFile(filePath string) (fileMissing bool, err error) {
	// Check if file exists
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return true, fmt.Errorf("failed to create file at %s: %w", filePath, err)
		}
		defer file.Close()

		_, err = file.WriteString(datafileHeader)
		if err != nil {
			return true, fmt.Errorf("failed to write content to file %s: %w", filePath, err)
		}

		return true, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to check if file exists at %s: %w", filePath, err)
	}

	// File already exists
	return false, nil
}
