package utils_test

import (
	"testing"
	"TravelSphere/utils"
)

// ── IsValidSlug Tests ──

func TestIsValidSlug_Valid(t *testing.T) {
	validSlugs := []string{
		"bangladesh",
		"united-states",
		"new-zealand",
		"usa",
		"bd",
		"123",
		"country-123",
	}
	for _, slug := range validSlugs {
		if !utils.IsValidSlug(slug) {
			t.Errorf("expected %q to be valid", slug)
		}
	}
}

func TestIsValidSlug_Invalid(t *testing.T) {
	invalidSlugs := []string{
		"",
		"United States",
		"BANGLADESH",
		"country name",
		"country!",
		"country@123",
		"country/sub",
	}
	for _, slug := range invalidSlugs {
		if utils.IsValidSlug(slug) {
			t.Errorf("expected %q to be invalid", slug)
		}
	}
}

// ── IsValidSearch Tests ──

func TestIsValidSearch_Valid(t *testing.T) {
	tests := []string{"", "bangladesh", "united states", "a"}
	for _, q := range tests {
		if !utils.IsValidSearch(q) {
			t.Errorf("expected %q to be valid search", q)
		}
	}
}

func TestIsValidSearch_TooLong(t *testing.T) {
	// 101 characters
	long := ""
	for i := 0; i < 101; i++ {
		long += "a"
	}
	if utils.IsValidSearch(long) {
		t.Error("query over 100 chars should be invalid")
	}
}

func TestIsValidSearch_Exactly100(t *testing.T) {
	// exactly 100 characters — valid
	exact := ""
	for i := 0; i < 100; i++ {
		exact += "a"
	}
	if !utils.IsValidSearch(exact) {
		t.Error("exactly 100 chars should be valid")
	}
}

// ── IsValidRegion Tests ──

func TestIsValidRegion_ValidRegions(t *testing.T) {
	valid := []string{
		"", "all", "Africa", "Americas",
		"Asia", "Europe", "Oceania", "Antarctic",
	}
	for _, r := range valid {
		if !utils.IsValidRegion(r) {
			t.Errorf("expected %q to be valid region", r)
		}
	}
}

func TestIsValidRegion_Invalid(t *testing.T) {
	invalid := []string{"Mars", "Unknown", "asia", "EUROPE"}
	for _, r := range invalid {
		if utils.IsValidRegion(r) {
			t.Errorf("expected %q to be invalid region", r)
		}
	}
}

// ── SanitizeString Tests ──

func TestSanitizeString_XSS(t *testing.T) {
	input := "<script>alert('xss')</script>"
	result := utils.SanitizeString(input)
	// < এবং > escape হবে
	if result == input {
		t.Error("XSS string should be sanitized")
	}
	if len(result) == 0 {
		t.Error("sanitized string should not be empty")
	}
}

func TestSanitizeString_Normal(t *testing.T) {
	input := "Visit Dhaka 2026"
	result := utils.SanitizeString(input)
	if result != input {
		t.Errorf("normal string should not change: %q", result)
	}
}

func TestSanitizeString_TrimSpace(t *testing.T) {
	result := utils.SanitizeString("  hello  ")
	if result != "hello" {
		t.Errorf("expected trimmed string, got %q", result)
	}
}

// ── TruncateString Tests ──

func TestTruncateString_Short(t *testing.T) {
	result := utils.TruncateString("hello", 10)
	if result != "hello" {
		t.Errorf("short string should not be truncated, got %q", result)
	}
}

func TestTruncateString_Long(t *testing.T) {
	result := utils.TruncateString("hello world", 5)
	if result != "hello..." {
		t.Errorf("expected 'hello...', got %q", result)
	}
}

func TestTruncateString_Exact(t *testing.T) {
	result := utils.TruncateString("hello", 5)
	if result != "hello" {
		t.Errorf("exact length should not be truncated, got %q", result)
	}
}