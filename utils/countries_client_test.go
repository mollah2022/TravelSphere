package utils_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"TravelSphere/models"
	"TravelSphere/utils"
)

type fakeClient struct {
	resp *http.Response
	err  error
}

func (f *fakeClient) Get(url string) (*http.Response, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.resp, nil
}

func (f *fakeClient) Do(req *http.Request) (*http.Response, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.resp, nil
}

func TestFetchJSON_Success(t *testing.T) {
	jsonBody := `[{"name": {"common":"X","official":"X"}}]`
	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(jsonBody)),
	}
	client := &fakeClient{resp: resp}

	var out []models.CountryResponse
	err := utils.FetchJSON(client, "http://example.test", &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 || out[0].Name.Common != "X" {
		t.Fatalf("unexpected decode result: %#v", out)
	}
}

func TestFetchJSON_Non200(t *testing.T) {
	resp := &http.Response{
		StatusCode: 500,
		Body:       io.NopCloser(strings.NewReader("error")),
	}
	client := &fakeClient{resp: resp}

	var out []models.CountryResponse
	err := utils.FetchJSON(client, "http://example.test", &out)
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestTransformCountryToDTO_FilterAndTransform(t *testing.T) {
	cr := models.CountryResponse{
		Name:       models.CountryName{Common: "United States", Official: "United States of America"},
		Capital:    []string{"Washington"},
		Region:     "Americas",
		Subregion:  "North America",
		Population: 330000000,
		Flags:      models.CountryFlag{PNG: "flag.png", Alt: "flag alt"},
		Currencies: map[string]models.Currency{"USD": {Name: "US Dollar", Symbol: "$"}},
		Languages:  map[string]string{"eng": "English"},
		CCA2:       "US",
		CCA3:       "USA",
		LatLng:     []float64{38.0, -97.0},
	}
	dto := utils.TransformCountryToDTO(cr)
	if dto.Name != "United States" || dto.Capital != "Washington" || dto.CCA2 != "US" {
		t.Fatalf("unexpected dto: %#v", dto)
	}
}

func TestFilterCountries_SearchAndRegion(t *testing.T) {
	countries := []models.CountryDTO{
		{Name: "France", Capital: "Paris", Region: "Europe"},
		{Name: "Japan", Capital: "Tokyo", Region: "Asia"},
		{Name: "United States", Capital: "Washington", Region: "Americas"},
	}

	// search by name
	res := utils.FilterCountries(countries, "japan", "")
	if len(res) != 1 || res[0].Name != "Japan" {
		t.Fatalf("search filter failed: %#v", res)
	}

	// region filter
	res = utils.FilterCountries(countries, "", "Europe")
	if len(res) != 1 || res[0].Region != "Europe" {
		t.Fatalf("region filter failed: %#v", res)
	}

	// combined search+region (no match)
	res = utils.FilterCountries(countries, "paris", "Asia")
	if len(res) != 0 {
		t.Fatalf("expected no matches, got: %#v", res)
	}

	// empty search+all region => return all
	res = utils.FilterCountries(countries, "", "all")
	if len(res) != 3 {
		t.Fatalf("expected 3, got %d", len(res))
	}
}

func TestNewCountriesClient_UsesEnvironment(t *testing.T) {
	orig := os.Getenv("REST_COUNTRIES_BASE_URL")
	defer os.Setenv("REST_COUNTRIES_BASE_URL", orig)

	os.Setenv("REST_COUNTRIES_BASE_URL", "https://mockapi.test")
	client := utils.NewCountriesClient()
	if client.BaseURL != "https://mockapi.test" {
		t.Fatalf("expected custom base url, got %q", client.BaseURL)
	}
}

func TestCountriesClient_FetchAll(t *testing.T) {
	jsonBody := `[{"name":{"common":"Chad"},"cca2":"TD","cca3":"TCD","capital":["N'Djamena"],"region":"Africa","population":17000000,"flags":{"png":"flag.png"},"currencies":{"XAF":{"name":"Central African CFA franc"}},"languages":{"fra":"French"},"latlng":[12.1,15.0]}]`
	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(jsonBody)),
	}
	client := &utils.CountriesClient{
		BaseURL:    "http://example.test",
		HTTPClient: &fakeClient{resp: resp},
	}
	countries, err := client.FetchAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(countries) != 1 || countries[0].Name.Common != "Chad" {
		t.Fatalf("unexpected result: %#v", countries)
	}
}

func TestCountriesClient_FetchByName(t *testing.T) {
	jsonBody := `[{"name":{"common":"India"},"cca2":"IN","cca3":"IND","capital":["New Delhi"],"region":"Asia","population":1380000000,"flags":{"png":"flag.png"},"currencies":{"INR":{"name":"Indian rupee"}},"languages":{"hin":"Hindi"},"latlng":[20.0,77.0]}]`
	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(jsonBody)),
	}
	client := &utils.CountriesClient{
		BaseURL:    "http://example.test",
		HTTPClient: &fakeClient{resp: resp},
	}
	countries, err := client.FetchByName("India")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(countries) != 1 || countries[0].Name.Common != "India" {
		t.Fatalf("unexpected result: %#v", countries)
	}
}

func TestCountriesClient_FetchAll_HTTPFailure(t *testing.T) {
	client := &utils.CountriesClient{
		BaseURL:    "http://example.test",
		HTTPClient: &fakeClient{err: errors.New("boom")},
	}
	countries, err := client.FetchAll()
	if err != nil {
		t.Fatalf("expected no error (fallback to mock), got %v", err)
	}
	if len(countries) == 0 {
		t.Fatal("expected mock data to be returned on HTTP failure")
	}
}

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
		Name:       models.CountryName{Common: "Bangladesh", Official: "People's Republic of Bangladesh"},
		CCA2:       "BD",
		CCA3:       "BGD",
		Capital:    []string{"Dhaka"},
		Region:     "Asia",
		Subregion:  "Southern Asia",
		Population: 170000000,
		Flags:      models.CountryFlag{PNG: "https://flag.png", Alt: "Flag of Bangladesh"},
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
