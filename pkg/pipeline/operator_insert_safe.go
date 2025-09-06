package pipeline

// insert_safe_string inserts all elements from source into target,
// avoiding duplicates and excluded values.
// Note: the target will be modified in place
func insert_safe_string(source []string, checkExclusion func(string) bool, target *[]string) {
	// Create a map to track existing values in target for O(1) lookup
	existing := make(map[string]struct{})
	for _, val := range *target {
		existing[val] = struct{}{}
	}

	// Add elements from source that aren't in existing or exclusions
	for _, val := range source {
		// Skip if value is in exclusions
		if checkExclusion(val) {
			continue
		}

		// Skip if value already exists in target
		if _, exists := existing[val]; exists {
			continue
		}

		// Add the value to target and mark it as existing
		*target = append(*target, val)
		existing[val] = struct{}{}
	}
}

// insert_safe inserts all elements of a Surface from source into target,
// avoiding duplicates and excluded values.
// Note: the target will be modified in place
func insert_safe(source Surface, exclusions Exclusions, target *Surface) {
	// Handle domains
	insert_safe_string(source.Domains, exclusions.Contains_domain, &target.Domains)

	// Handle IPs
	insert_safe_string(source.IPs, exclusions.Contains_ip, &target.IPs)

	// Handle URLs
	insert_safe_string(source.URLs, exclusions.Contains_url, &target.URLs)
}
