package utils

import (
	"strings"
	"unicode"
)

// IsValidSlug checks whether a slug is valid.
// A valid slug contains only lowercase letters, numbers, and hyphens.
func IsValidSlug(slug string) bool {
	if slug == "" {
		return false
	}
	for _, ch := range slug {
		if ch == '-' || ch == '_' {
			continue
		}
		if ch == '\'' {
			continue
		}
		if unicode.IsLower(ch) || unicode.IsDigit(ch) {
			continue
		}
		return false
	}
	return true
}

// IsValidSearch search query validate
func IsValidSearch(query string) bool {
	query = strings.TrimSpace(query)
	return len(query) <= 100 // max 100 characters
}

// IsValidRegion region filter validate
func IsValidRegion(region string) bool {
	if region == "" || region == "all" {
		return true
	}
	validRegions := map[string]bool{
		"Africa":    true,
		"Americas":  true,
		"Asia":      true,
		"Europe":    true,
		"Oceania":   true,
		"Antarctic": true,
	}
	return validRegions[region]
}

// SanitizeString removes dangerous characters from a string.
func SanitizeString(s string) string {
	s = strings.TrimSpace(s)
	// Basic XSS prevention
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}

// TruncateString truncates a string to a maximum length.
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
