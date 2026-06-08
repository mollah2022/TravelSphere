package utils_test

import (
	"testing"
	"TravelSphere/models"
	"TravelSphere/utils"
)

// ── FormatPopulation Tests ──

func TestFormatPopulation_Billions(t *testing.T) {
	result := utils.FormatPopulation(1_400_000_000)
	if result != "1.4B" {
		t.Errorf("expected '1.4B', got %q", result)
	}
}

func TestFormatPopulation_Millions(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{170_000_000, "170.0M"},
		{2_400_000, "2.4M"},
		{1_000_000, "1.0M"},
	}
	for _, tt := range tests {
		got := utils.FormatPopulation(tt.input)
		if got != tt.expected {
			t.Errorf("FormatPopulation(%d) = %q, want %q",
				tt.input, got, tt.expected)
		}
	}
}

func TestFormatPopulation_Thousands(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{47_400, "47.4K"},
		{1_000, "1.0K"},
		{500_000, "500.0K"},
	}
	for _, tt := range tests {
		got := utils.FormatPopulation(tt.input)
		if got != tt.expected {
			t.Errorf("FormatPopulation(%d) = %q, want %q",
				tt.input, got, tt.expected)
		}
	}
}

func TestFormatPopulation_Small(t *testing.T) {
	result := utils.FormatPopulation(800)
	if result != "800" {
		t.Errorf("expected '800', got %q", result)
	}
}

func TestFormatPopulation_Zero(t *testing.T) {
	result := utils.FormatPopulation(0)
	if result != "0" {
		t.Errorf("expected '0', got %q", result)
	}
}

// ── FormatLanguages Tests ──

func TestFormatLanguages_Single(t *testing.T) {
	langs := map[string]string{"ben": "Bengali"}
	result := utils.FormatLanguages(langs)
	if result != "Bengali" {
		t.Errorf("expected 'Bengali', got %q", result)
	}
}

func TestFormatLanguages_Empty(t *testing.T) {
	result := utils.FormatLanguages(map[string]string{})
	if result != "N/A" {
		t.Errorf("expected 'N/A', got %q", result)
	}
}

func TestFormatLanguages_Nil(t *testing.T) {
	result := utils.FormatLanguages(nil)
	if result != "N/A" {
		t.Errorf("expected 'N/A' for nil, got %q", result)
	}
}

// ── FormatCurrenciesWithCode Tests ──

func TestFormatCurrenciesWithCode_Single(t *testing.T) {
	curr := map[string]models.Currency{
		"BDT": {Name: "Bangladeshi taka", Symbol: "৳"},
	}
	result := utils.FormatCurrenciesWithCode(curr)
	if result != "BDT (Bangladeshi taka)" {
		t.Errorf("unexpected result: %q", result)
	}
}

func TestFormatCurrenciesWithCode_Empty(t *testing.T) {
	result := utils.FormatCurrenciesWithCode(map[string]models.Currency{})
	if result != "N/A" {
		t.Errorf("expected 'N/A', got %q", result)
	}
}

// ── FormatCapital Tests ──

func TestFormatCapital_HasCapital(t *testing.T) {
	result := utils.FormatCapital([]string{"Dhaka"})
	if result != "Dhaka" {
		t.Errorf("expected 'Dhaka', got %q", result)
	}
}

func TestFormatCapital_Empty(t *testing.T) {
	result := utils.FormatCapital([]string{})
	if result != "N/A" {
		t.Errorf("expected 'N/A', got %q", result)
	}
}

func TestFormatCapital_MultipleReturnsFirst(t *testing.T) {
	result := utils.FormatCapital([]string{"Seat of Government", "Capital"})
	if result != "Seat of Government" {
		t.Errorf("expected first capital, got %q", result)
	}
}

// ── NameToSlug Tests ──

func TestNameToSlug_Simple(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Bangladesh", "bangladesh"},
		{"United States", "united-states"},
		{"New Zealand", "new-zealand"},
		{"Sri Lanka", "sri-lanka"},
	}
	for _, tt := range tests {
		got := utils.NameToSlug(tt.input)
		if got != tt.expected {
			t.Errorf("NameToSlug(%q) = %q, want %q",
				tt.input, got, tt.expected)
		}
	}
}

func TestNameToSlug_SpecialChars(t *testing.T) {
	// apostrophe remove হবে
	result := utils.NameToSlug("Côte d'Ivoire")
	if result == "" {
		t.Error("slug should not be empty")
	}
}

// ── SlugToName Tests ──

func TestSlugToName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"bangladesh", "Bangladesh"},
		{"united-states", "United States"},
		{"new-zealand", "New Zealand"},
	}
	for _, tt := range tests {
		got := utils.SlugToName(tt.input)
		if got != tt.expected {
			t.Errorf("SlugToName(%q) = %q, want %q",
				tt.input, got, tt.expected)
		}
	}
}

// ── FormatKinds Tests ──

func TestFormatKinds_Normal(t *testing.T) {
	result := utils.FormatKinds("museums,historic_architecture")
	if len(result) != 2 {
		t.Fatalf("expected 2 items, got %d", len(result))
	}
	if result[0] != "museums" {
		t.Errorf("expected 'museums', got %q", result[0])
	}
	if result[1] != "historic architecture" {
		t.Errorf("expected 'historic architecture', got %q", result[1])
	}
}

func TestFormatKinds_Empty(t *testing.T) {
	result := utils.FormatKinds("")
	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestFormatKinds_Single(t *testing.T) {
	result := utils.FormatKinds("natural")
	if len(result) != 1 || result[0] != "natural" {
		t.Errorf("unexpected result: %v", result)
	}
}

// ── FormatRegion Tests ──

func TestFormatRegion_WithSubregion(t *testing.T) {
	result := utils.FormatRegion("Asia", "Southern Asia")
	if result != "Asia - Southern Asia" {
		t.Errorf("expected 'Asia - Southern Asia', got %q", result)
	}
}

func TestFormatRegion_WithoutSubregion(t *testing.T) {
	result := utils.FormatRegion("Asia", "")
	if result != "Asia" {
		t.Errorf("expected 'Asia', got %q", result)
	}
}