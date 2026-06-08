package utils_test

import (
	"TravelSphere/models"
	"TravelSphere/utils"
	"testing"
)

func TestFormatPopulation(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{170000000, "170.0M"},
		{2400000, "2.4M"},
		{47400, "47.4K"},
		{500, "500"},
		{1000000000, "1.0B"},
	}
	for _, tt := range tests {
		result := utils.FormatPopulation(tt.input)
		if result != tt.expected {
			t.Errorf("FormatPopulation(%d) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestFormatLanguages(t *testing.T) {
	langs := map[string]string{"ben": "Bengali"}
	result := utils.FormatLanguages(langs)
	if result != "Bengali" {
		t.Errorf("expected 'Bengali', got %q", result)
	}

	empty := map[string]string{}
	if utils.FormatLanguages(empty) != "N/A" {
		t.Error("empty map should return N/A")
	}
}

func TestFormatCurrenciesWithCode(t *testing.T) {
	curr := map[string]models.Currency{
		"BDT": {Name: "Bangladeshi taka", Symbol: "৳"},
	}
	result := utils.FormatCurrenciesWithCode(curr)
	if result != "BDT (Bangladeshi taka)" {
		t.Errorf("unexpected result: %q", result)
	}
}

func TestNameToSlug(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"United States", "united-states"},
		{"Bangladesh", "bangladesh"},
		{"Côte d'Ivoire", "côte-d'ivoire"},
	}
	for _, tt := range tests {
		result := utils.NameToSlug(tt.input)
		if result != tt.expected {
			t.Errorf("NameToSlug(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestIsValidSlug(t *testing.T) {
	if !utils.IsValidSlug("bangladesh") {
		t.Error("bangladesh should be valid")
	}
	if !utils.IsValidSlug("united-states") {
		t.Error("united-states should be valid")
	}
	if utils.IsValidSlug("") {
		t.Error("empty should be invalid")
	}
	if utils.IsValidSlug("United States") {
		t.Error("uppercase/space should be invalid")
	}
}

func TestFormatKinds(t *testing.T) {
	result := utils.FormatKinds("museums,historic_architecture")
	if len(result) != 2 {
		t.Errorf("expected 2 kinds, got %d", len(result))
	}
	if result[1] != "historic architecture" {
		t.Errorf("expected 'historic architecture', got %q", result[1])
	}
}

func TestIsValidRegion(t *testing.T) {
	if !utils.IsValidRegion("Asia") {
		t.Error("Asia should be valid")
	}
	if !utils.IsValidRegion("") {
		t.Error("empty should be valid (means all)")
	}
	if utils.IsValidRegion("Mars") {
		t.Error("Mars should be invalid")
	}
}