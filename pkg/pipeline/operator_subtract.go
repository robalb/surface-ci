package pipeline

func Subtract(original []string, toRemove []string) []string {
    // Create a map for efficient lookup of elements to remove
    removeMap := make(map[string]bool)
    for _, item := range toRemove {
        removeMap[item] = true
    }
    
    // Create result slice with initial capacity of original
    result := make([]string, 0, len(original))
    
    // Add elements from original that aren't in toRemove
    for _, item := range original {
        if !removeMap[item] {
            result = append(result, item)
        }
    }
    
    return result
}
