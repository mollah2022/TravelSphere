package utils

import "regexp"

// slugRegex defines a valid URL slug pattern:
// lowercase letters (a-z)
// numbers (0-9)
// words separated by single hyphens (-)
// Example: "travel-dhaka-2026"
var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// IsValidSlug checks whether a given string is a valid slug or not
func IsValidSlug(slug string) bool {
	if len(slug) == 0 || len(slug) > 100 {
		return false
	}
	return slugRegex.MatchString(slug)
}
