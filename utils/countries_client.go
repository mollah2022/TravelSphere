package utils

import (
	"TravelSphere/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// CountriesClient handles communication with the REST Countries API
type CountriesClient struct {
	baseURL string
	client  *http.Client
}

// NewCountriesClient creates and returns a new CountriesClient instance
func NewCountriesClient() *CountriesClient {
	base := os.Getenv("REST_COUNTRIES_BASE_URL")
	if base == "" {
		base = "https://restcountries.com/v3.1"
	}
	return &CountriesClient{
		baseURL: base,
		client:  NewHTTPClient(10),
	}
}

// GetAll fetches all countries data from the API
func (c *CountriesClient) GetAll() ([]models.CountryResponse, error) {
	fields := "name,cca2,cca3,capital,region,subregion,population,flags,currencies,languages,latlng,timezones,tld"
	url := fmt.Sprintf("%s/all?fields=%s", c.baseURL, fields)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("countries API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("countries API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var countries []models.CountryResponse
	if err := json.Unmarshal(body, &countries); err != nil {
		return nil, fmt.Errorf("failed to parse countries response: %w", err)
	}

	return countries, nil
}

// GetByName fetches a single country by its full name from the API
func (c *CountriesClient) GetByName(name string) (*models.CountryResponse, error) {
	url := fmt.Sprintf("%s/name/%s?fullText=true", c.baseURL, name)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("countries API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, models.ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("countries API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var countries []models.CountryResponse
	if err := json.Unmarshal(body, &countries); err != nil {
		return nil, fmt.Errorf("failed to parse country response: %w", err)
	}

	if len(countries) == 0 {
		return nil, models.ErrNotFound
	}

	return &countries[0], nil
}

// SlugToName converts URL slug into a readable country name
func SlugToName(slug string) string {
	words := strings.Split(slug, "-")
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}
