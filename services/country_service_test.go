package services_test

import (
	"TravelSphere/models"
	"TravelSphere/services"
	"errors"
	"testing"
)

// ── Mock Countries Client ──

type mockCountriesClient struct {
	data      []models.CountryResponse
	shouldErr bool
}

func (m *mockCountriesClient) FetchAll() ([]models.CountryResponse, error) {
	if m.shouldErr {
		return nil, errors.New("mock fetch error")
	}
	return m.data, nil
}

func (m *mockCountriesClient) FetchByName(name string) ([]models.CountryResponse, error) {
	return m.FetchAll()
}

// ── Test Data ──

func mockCountries() []models.CountryResponse {
	return []models.CountryResponse{
		{
			Name:      models.CountryName{Common: "Bangladesh", Official: "People's Republic of Bangladesh"},
			CCA2:      "BD",
			CCA3:      "BGD",
			Capital:   []string{"Dhaka"},
			Region:    "Asia",
			Subregion: "Southern Asia",
			Population: 170_000_000,
			Flags:      models.CountryFlag{PNG: "https://bd-flag.png"},
			Currencies: map[string]models.Currency{
				"BDT": {Name: "Bangladeshi taka"},
			},
			Languages: map[string]string{"ben": "Bengali"},
			LatLng:    []float64{24.0, 90.0},
		},
		{
			Name:      models.CountryName{Common: "France", Official: "French Republic"},
			CCA2:      "FR",
			CCA3:      "FRA",
			Capital:   []string{"Paris"},
			Region:    "Europe",
			Subregion: "Western Europe",
			Population: 67_000_000,
			Flags:      models.CountryFlag{PNG: "https://fr-flag.png"},
			Currencies: map[string]models.Currency{
				"EUR": {Name: "Euro"},
			},
			Languages: map[string]string{"fra": "French"},
			LatLng:    []float64{46.0, 2.0},
		},
		{
			Name:      models.CountryName{Common: "Japan", Official: "Japan"},
			CCA2:      "JP",
			CCA3:      "JPN",
			Capital:   []string{"Tokyo"},
			Region:    "Asia",
			Subregion: "Eastern Asia",
			Population: 125_000_000,
			Flags:      models.CountryFlag{PNG: "https://jp-flag.png"},
			Currencies: map[string]models.Currency{
				"JPY": {Name: "Japanese yen"},
			},
			Languages: map[string]string{"jpn": "Japanese"},
			LatLng:    []float64{36.0, 138.0},
		},
	}
}

func newCountryService(shouldErr bool) *services.CountryService {
	client := &mockCountriesClient{
		data:      mockCountries(),
		shouldErr: shouldErr,
	}
	return services.NewCountryService(client)
}

// ── Tests ──

func TestGetAllCountries_Success(t *testing.T) {
	svc := newCountryService(false)

	countries, err := svc.GetAllCountries()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(countries) != 3 {
		t.Errorf("expected 3 countries, got %d", len(countries))
	}
}

func TestGetAllCountries_APIError(t *testing.T) {
	svc := newCountryService(true)

	_, err := svc.GetAllCountries()
	if err == nil {
		t.Error("expected error when API fails")
	}
}

func TestGetAllCountries_Cached(t *testing.T) {
	svc := newCountryService(false)

	// প্রথম call
	first, _ := svc.GetAllCountries()
	// দ্বিতীয় call — cache থেকে আসবে
	second, _ := svc.GetAllCountries()

	if len(first) != len(second) {
		t.Error("cached result should be same as first result")
	}
}

func TestSearchCountries_ByName(t *testing.T) {
	svc := newCountryService(false)

	results, err := svc.SearchCountries("bang", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result for 'bang', got %d", len(results))
	}
	if results[0].Name != "Bangladesh" {
		t.Errorf("expected Bangladesh, got %q", results[0].Name)
	}
}

func TestSearchCountries_ByRegion(t *testing.T) {
	svc := newCountryService(false)

	results, err := svc.SearchCountries("", "Asia")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 Asian countries, got %d", len(results))
	}
}

func TestSearchCountries_NoMatch(t *testing.T) {
	svc := newCountryService(false)

	results, err := svc.SearchCountries("zzz", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestGetCountryBySlug_Found(t *testing.T) {
	svc := newCountryService(false)

	country, err := svc.GetCountryBySlug("bangladesh")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if country.Name != "Bangladesh" {
		t.Errorf("expected Bangladesh, got %q", country.Name)
	}
}

func TestGetCountryBySlug_ByCCA2(t *testing.T) {
	svc := newCountryService(false)

	country, err := svc.GetCountryBySlug("fr")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if country.Name != "France" {
		t.Errorf("expected France, got %q", country.Name)
	}
}

func TestGetCountryBySlug_ByCCA3(t *testing.T) {
	svc := newCountryService(false)

	country, err := svc.GetCountryBySlug("jpn")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if country.Name != "Japan" {
		t.Errorf("expected Japan, got %q", country.Name)
	}
}

func TestGetCountryBySlug_NotFound(t *testing.T) {
	svc := newCountryService(false)

	_, err := svc.GetCountryBySlug("nonexistent-country")
	if err != models.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestGetSearchSuggestions_LimitedTo8(t *testing.T) {
	// 10টা country দিয়ে test করো
	data := make([]models.CountryResponse, 10)
	for i := range data {
		data[i] = models.CountryResponse{
			Name:    models.CountryName{Common: "Country"},
			CCA2:    "XX",
			Capital: []string{"Capital"},
			Region:  "Asia",
			Flags:   models.CountryFlag{},
		}
	}
	client := &mockCountriesClient{data: data}
	svc := services.NewCountryService(client)

	results, err := svc.GetSearchSuggestions("c")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) > 8 {
		t.Errorf("suggestions should be limited to 8, got %d", len(results))
	}
}

func TestGetSearchSuggestions_EmptyQuery(t *testing.T) {
	svc := newCountryService(false)

	results, err := svc.GetSearchSuggestions("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("empty query should return 0 suggestions, got %d",
			len(results))
	}
}

func TestGetFeaturedCountries_Success(t *testing.T) {
	svc := newCountryService(false)

	results, err := svc.GetFeaturedCountries()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Errorf("expected featured countries, got %d", len(results))
	}
	// Should include Bangladesh which is in the featured list
	found := false
	for _, c := range results {
		if c.Name == "Bangladesh" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected Bangladesh in featured countries")
	}
}

func TestGetFeaturedCountries_APIError(t *testing.T) {
	svc := newCountryService(true)

	_, err := svc.GetFeaturedCountries()
	if err == nil {
		t.Error("expected error when API fails")
	}
}

func TestSearchCountriesWithPagination_Success(t *testing.T) {
	svc := newCountryService(false)

	results, err := svc.SearchCountriesWithPagination("", "", 10, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Errorf("expected 3 countries, got %d", len(results))
	}
}

func TestSearchCountriesWithPagination_WithOffset(t *testing.T) {
	svc := newCountryService(false)

	results, err := svc.SearchCountriesWithPagination("", "", 10, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 country with offset 2, got %d", len(results))
	}
}

func TestSearchCountriesWithPagination_WithLimit(t *testing.T) {
	svc := newCountryService(false)

	results, err := svc.SearchCountriesWithPagination("", "", 2, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 countries with limit 2, got %d", len(results))
	}
}

func TestSearchCountriesWithPagination_APIError(t *testing.T) {
	svc := newCountryService(true)

	_, err := svc.SearchCountriesWithPagination("", "", 10, 0)
	if err == nil {
		t.Error("expected error when API fails")
	}
}