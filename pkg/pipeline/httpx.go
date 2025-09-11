package pipeline

import (
	"fmt"
	"sync"

	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/httpx/runner"
)

// Result represents the result of checking a URL
type Result struct {
	URL        string
	StatusCode int
	Error      error
}

// Httpx takes a Surface struct and returns a list of results
func Httpx(surface Surface, threads int) ([]Result, error) {
	// Combine all targets
	var targets []string
	targets = append(targets, surface.URLs...)
	targets = append(targets, surface.Domains...)
	targets = append(targets, surface.IPs...)
	
	// Create a slice to store results
	var results []Result
	var mu sync.Mutex
	
	// Set up httpx options
	options := runner.Options{
		Methods:         "GET",
		InputTargetHost: goflags.StringSlice(targets),
		Threads:         threads,
		OnResult: func(r runner.Result) {
			result := Result{
				StatusCode: r.StatusCode,
				Error:      r.Err,
			}
			
			// If no error, add URL information
			if r.Err == nil {
				result.URL = r.URL
			}
			
			// Thread-safe append to results
			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		},
	}
	
	// Validate options
	if err := options.ValidateOptions(); err != nil {
		return nil, err
	}
	
	// Create and run httpx
	httpxRunner, err := runner.New(&options)
	if err != nil {
		return nil, err
	}
	defer httpxRunner.Close()
	
	// Run the enumeration
	httpxRunner.RunEnumeration()
	
	return results, nil
}

