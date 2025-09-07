package pipeline

import (
	"bytes"
	"io"
	"strings"

	"github.com/projectdiscovery/alterx"
)

// Alterx takes a list of domains and returns plausible alternative domains
// generated using the alterx tool, a wordlist and a patterns list
func Alterx(domains []string) ([]string, error) {

	// Configure alterx options
	alterxOpts := &alterx.Options{
		Domains: domains,
		MaxSize: 1000,
		Enrich:  true,
		//TODO: configure words and patters using an LLM
	}

	// Create a new alterx instance
	m, err := alterx.New(alterxOpts)
	if err != nil {
		return nil, err
	}

	// Capture the output in a buffer instead of writing to stdout
	var buffer bytes.Buffer
	err = m.ExecuteWithWriter(&buffer)
	if err != nil {
		return nil, err
	}

	// Convert buffer to string slice
	var results []string
	for {
		line, err := buffer.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return results, err
		}
		// Trim newline and add to results
		domain := strings.TrimSpace(line)
		if domain != "" {
			results = append(results, domain)
		}
	}

	return results, nil
}
