package utils_test

import (
	"TravelSphere/models"
	"TravelSphere/utils"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

// MockHTTPClient test এ real HTTP call না করে fake response দেয়
type MockHTTPClient struct {
	ResponseBody string
	StatusCode   int
	ShouldError  bool
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	if m.ShouldError {
		return nil, fmt.Errorf("mock network error")
	}
	return &http.Response{
		StatusCode: m.StatusCode,
		Body:       io.NopCloser(strings.NewReader(m.ResponseBody)),
	}, nil
}

func TestTransformCountryToDTO(t *testing.T) {
	country := models.CountryResponse{
		Name:      models.CountryName{Common: "Bangladesh", Official: "People's Republic of Bangladesh"},
		CCA2:      "BD",
		CCA3:      "BGD",
		Capital:   []string{"Dhaka"},
		Region:    "Asia",
		Subregion: "Southern Asia",
		Population: 170000000,
		Flags:     models.CountryFlag{PNG: "https://flag.png", Alt: "Flag of Bangladesh"},
		Currencies: map[string]models.Currency{
			"BDT": {Name: "Bangladeshi taka", Symbol: "৳"},
		},
		Languages: map[string]string{"ben": "Bengali"},
		LatLng:    []float64{24.0, 90.0},
	}

	dto := utils.TransformCountryToDTO(country)

	if dto.Slug != "bangladesh" {
		t.Errorf("expected slug 'bangladesh', got %q", dto.Slug)
	}
	if dto.Capital != "Dhaka" {
		t.Errorf("expected capital 'Dhaka', got %q", dto.Capital)
	}
	if dto.Latitude != 24.0 {
		t.Errorf("expected lat 24.0, got %f", dto.Latitude)
	}
	if dto.Languages != "Bengali" {
		t.Errorf("expected 'Bengali', got %q", dto.Languages)
	}
}

func TestFilterCountries(t *testing.T) {
	countries := []models.CountryDTO{
		{Name: "Bangladesh", Capital: "Dhaka", Region: "Asia"},
		{Name: "France", Capital: "Paris", Region: "Europe"},
		{Name: "Japan", Capital: "Tokyo", Region: "Asia"},
	}

	// Search test
	result := utils.FilterCountries(countries, "bang", "")
	if len(result) != 1 || result[0].Name != "Bangladesh" {
		t.Errorf("search filter failed: %+v", result)
	}

	// Region test
	result2 := utils.FilterCountries(countries, "", "Asia")
	if len(result2) != 2 {
		t.Errorf("region filter failed, expected 2 got %d", len(result2))
	}

	// No filter
	result3 := utils.FilterCountries(countries, "", "all")
	if len(result3) != 3 {
		t.Errorf("no filter should return all 3, got %d", len(result3))
	}
}