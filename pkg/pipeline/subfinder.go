package pipeline

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
)

// ExpandDomains takes a list of domains and exclusions, uses subfinder to discover subdomains,
// and returns the expanded list of domains, filtering out any excluded domains.
func Subfinder(
	ctx context.Context,
	domains []string,
) ([]string, error) {
	if len(domains) == 0 {
		return []string{}, nil
	}

	subfinderOpts := &runner.Options{
		Threads:            5,
		Timeout:            10,
		MaxEnumerationTime: 30,
		Silent:             false,
	}

	subfinderRunner, err := runner.NewRunner(subfinderOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create subfinder runner: %v", err)
	}

	// Convert domains array to reader expected by subfinder
	domainsStr := strings.Join(domains, "\n")
	domainsReader := strings.NewReader(domainsStr)

	// Buffer to capture subfinder output
	output := &bytes.Buffer{}

	// Run subfinder on domains
	if err := subfinderRunner.EnumerateMultipleDomainsWithCtx(ctx, domainsReader, []io.Writer{output}); err != nil {
		return nil, fmt.Errorf("subfinder enumeration failed: %v", err)
	}

	// Process results
	results := output.String()
	if results == "" {
		return []string{}, nil
	}

	discoveredDomains := strings.Split(strings.TrimSpace(results), "\n")
	return discoveredDomains, nil

	// Filter out excluded domains
	// var expandedDomains []string
	// for _, domain := range discoveredDomains {
	// 	if domain == "" {
	// 		continue
	// 	}
	//
	// 	excluded := false
	// 	for _, excludedDomain := range excludedDomains {
	// 		// Check if domain exactly matches an excluded domain
	// 		if domain == excludedDomain {
	// 			excluded = true
	// 			break
	// 		}
	//
	// 		// Check if domain is a subdomain of an excluded domain
	// 		if strings.HasSuffix(domain, "."+excludedDomain) {
	// 			excluded = true
	// 			break
	// 		}
	// 	}
	//
	// 	if !excluded {
	// 		expandedDomains = append(expandedDomains, domain)
	// 	}
	// }
	// return expandedDomains, nil
}
