package search

import "strings"

// matchPath returns true if specified path exists in haystack
func matchPath(path string, haystack []string) bool {
	for _, v := range haystack {
		if strings.HasPrefix(v, path) {
			return true
		}
	}
	return false
}
