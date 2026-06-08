package utils

import (
	"TravelSphere/models"
	"fmt"
	"os"
	"strings"
)

// CountriesClientInterface mock করার জন্য interface
type CountriesClientInterface interface {
	FetchAll() ([]models.CountryResponse, error)
	FetchByName(name string) ([]models.CountryResponse, error)
}

// CountriesClient REST Countries API client
type CountriesClient struct {
	BaseURL    string
	HTTPClient HTTPClient
}

// NewCountriesClient নতুন REST Countries client তৈরি করে
func NewCountriesClient() *CountriesClient {
	baseURL := os.Getenv("REST_COUNTRIES_BASE_URL")
	if baseURL == "" {
		baseURL = "https://restcountries.com/v3.1"
	}
	return &CountriesClient{
		BaseURL:    baseURL,
		HTTPClient: NewHTTPClient(10),
	}
}

// FetchAll সব দেশের তথ্য fetch করে
// tld আর timezones বাদ দেওয়া হয়েছে — API 400 error দিচ্ছিল
func (c *CountriesClient) FetchAll() ([]models.CountryResponse, error) {
	url := fmt.Sprintf("%s/all?fields=name,cca2,cca3,capital,region,population,flags,currencies,languages,latlng", c.BaseURL)

	var countries []models.CountryResponse
	if err := FetchJSON(c.HTTPClient, url, &countries); err != nil {
		return nil, fmt.Errorf("FetchAll failed: %w", err)
	}
	return countries, nil
}

// FetchByName নাম দিয়ে দেশ খোঁজে
func (c *CountriesClient) FetchByName(name string) ([]models.CountryResponse, error) {
	url := fmt.Sprintf("%s/name/%s?fields=name,cca2,cca3,capital,region,population,flags,currencies,languages,latlng", c.BaseURL, name)

	var countries []models.CountryResponse
	if err := FetchJSON(c.HTTPClient, url, &countries); err != nil {
		return nil, fmt.Errorf("FetchByName failed: %w", err)
	}
	return countries, nil
}
// TransformCountryToDTO CountryResponse কে CountryDTO তে convert করে
func TransformCountryToDTO(c models.CountryResponse) models.CountryDTO {
	lat, lon := 0.0, 0.0
	if len(c.LatLng) >= 2 {
		lat = c.LatLng[0]
		lon = c.LatLng[1]
	}

	return models.CountryDTO{
		Slug:         NameToSlug(c.Name.Common),
		Name:         c.Name.Common,
		OfficialName: c.Name.Official,
		Capital:      FormatCapital(c.Capital),
		Region:       c.Region,
		Subregion:    c.Subregion,
		Population:   c.Population,
		FlagURL:      c.Flags.PNG,
		FlagAlt:      c.Flags.Alt,
		Currencies:   FormatCurrenciesWithCode(c.Currencies),
		Languages:    FormatLanguages(c.Languages),
		CCA2:         c.CCA2,
		CCA3:         c.CCA3,
		Latitude:     lat,
		Longitude:    lon,
	}
}

// FilterCountries search ও region দিয়ে filter করে
func FilterCountries(countries []models.CountryDTO, search, region string) []models.CountryDTO {
	if search == "" && (region == "" || region == "all") {
		return countries
	}

	search = strings.ToLower(strings.TrimSpace(search))
	result := make([]models.CountryDTO, 0)

	for _, c := range countries {
		// Region filter
		if region != "" && region != "all" {
			if !strings.EqualFold(c.Region, region) {
				continue
			}
		}
		// Search filter — name বা capital দিয়ে খোঁজে
		if search != "" {
			nameMatch := strings.Contains(strings.ToLower(c.Name), search)
			capitalMatch := strings.Contains(strings.ToLower(c.Capital), search)
			if !nameMatch && !capitalMatch {
				continue
			}
		}
		result = append(result, c)
	}
	return result
}