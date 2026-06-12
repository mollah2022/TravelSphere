package utils

import (
	"TravelSphere/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// CountriesClientInterface mock করার জন্য interface
type CountriesClientInterface interface {
	FetchAll() ([]models.CountryResponse, error)
	FetchByName(name string) ([]models.CountryResponse, error)
}

// CountriesClient REST Countries API client
type CountriesClient struct {
	BaseURL    string
	APIKey     string
	HTTPClient HTTPClient
	rawClient  *http.Client
}

// APIResponse wrapper for new v5 API response format
type APIResponse struct {
	Success bool                     `json:"success"`
	Data    []models.CountryResponse `json:"data"`
	Errors  []map[string]string      `json:"errors"`
}

// V5APIResponse wrapper for new v5 API format
type V5APIResponse struct {
	Data struct {
		Objects []V5CountryData `json:"objects"`
	} `json:"data"`
}

// V5CountryData represents a country from the new v5 API
type V5CountryData struct {
	Names struct {
		Common   string `json:"common"`
		Official string `json:"official"`
	} `json:"names"`
	Codes struct {
		Alpha2 string `json:"alpha_2"`
		Alpha3 string `json:"alpha_3"`
	} `json:"codes"`
	Capitals []struct {
		Name        string `json:"name"`
		Coordinates struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"coordinates"`
	} `json:"capitals"`
	Region     string           `json:"region"`
	Subregion  string           `json:"subregion"`
	Population int64            `json:"population"`
	Flag       V5FlagData       `json:"flag"`
	Currencies []V5CurrencyData `json:"currencies"`
	Languages  []V5LanguageData `json:"languages"`
}

// V5FlagData represents flag data from v5 API
type V5FlagData struct {
	Description string `json:"description"`
	Emoji       string `json:"emoji"`
	HtmlEntity  string `json:"html_entity"`
	Unicode     string `json:"unicode"`
	UrlPng      string `json:"url_png"`
	UrlSvg      string `json:"url_svg"`
}

// V5CurrencyData represents currency from v5 API (note: it's a list, not a map)
type V5CurrencyData struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

// V5LanguageData represents language from v5 API
type V5LanguageData struct {
	Bcp47      string `json:"bcp47"`
	Iso6391    string `json:"iso639_1"`
	Iso6392b   string `json:"iso639_2b"`
	Iso6392t   string `json:"iso639_2t"`
	Iso6393    string `json:"iso639_3"`
	Name       string `json:"name"`
	NativeName string `json:"native_name"`
}

// NewCountriesClient নতুন REST Countries client তৈরি করে
func NewCountriesClient() *CountriesClient {
	baseURL := os.Getenv("REST_COUNTRIES_BASE_URL")
	// if baseURL == "" {
	// 	baseURL = "https://api.restcountries.com/countries/v5"
	// }
	apiKey := os.Getenv("REST_COUNTRIES_API_KEY")
	// if apiKey == "" {
	// 	// Fallback API key if not provided
	// 	apiKey = "rc_live_e3842f257f19449086b10cfcfba5e3f2"
	// }
	return &CountriesClient{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		HTTPClient: NewHTTPClient(10),
		rawClient:  &http.Client{Timeout: 10 * time.Second},
	}
}

// FetchAll সব দেশের তথ্য fetch করে (250+ countries)
// Uses pagination to fetch all countries from the v5 API
// Makes multiple requests with limit=100 and increasing offsets
func (c *CountriesClient) FetchAll() ([]models.CountryResponse, error) {
	allCountries := make([]models.CountryResponse, 0)
	const maxLimit = 100
	offset := 0
	consecutiveEmptyResponses := 0
	maxAttempts := 5 // Prevent infinite loops; REST Countries has ~254 countries

	// Keep fetching until we get an empty response (pagination exhausted)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		params := url.Values{}
		params.Add("limit", strconv.Itoa(maxLimit))
		params.Add("offset", strconv.Itoa(offset))

		urlStr := fmt.Sprintf("%s?%s", c.BaseURL, params.Encode())

		req, err := http.NewRequest("GET", urlStr, nil)
		if err != nil {
			// If we already have some data, return it; otherwise use mock data
			if len(allCountries) > 0 {
				return allCountries, nil
			}
			return GetMockCountriesData(), nil
		}

		// Add Authorization header
		if c.APIKey != "" {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
		}

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			// If we already have some data, return it; otherwise use mock data
			if len(allCountries) > 0 {
				return allCountries, nil
			}
			return GetMockCountriesData(), nil
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			// If we already have some data, return it; otherwise use mock data
			if len(allCountries) > 0 {
				return allCountries, nil
			}
			return GetMockCountriesData(), nil
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			// If we already have some data, return it; otherwise use mock data
			if len(allCountries) > 0 {
				return allCountries, nil
			}
			return GetMockCountriesData(), nil
		}

		// Try to unmarshal as v5 API format
		var v5Response V5APIResponse
		if err := json.Unmarshal(body, &v5Response); err == nil && len(v5Response.Data.Objects) > 0 {
			convertedCountries := convertV5ToCountryResponse(v5Response.Data.Objects)
			allCountries = append(allCountries, convertedCountries...)
			consecutiveEmptyResponses = 0
			offset += maxLimit
			continue
		}

		// Try to unmarshal as wrapper format
		var wrapper APIResponse
		if err := json.Unmarshal(body, &wrapper); err == nil && wrapper.Success && len(wrapper.Data) > 0 {
			allCountries = append(allCountries, wrapper.Data...)
			consecutiveEmptyResponses = 0
			offset += maxLimit
			continue
		}

		// Fallback: try to unmarshal as direct array format
		var countries []models.CountryResponse
		if err := json.Unmarshal(body, &countries); err == nil && len(countries) > 0 {
			allCountries = append(allCountries, countries...)
			consecutiveEmptyResponses = 0
			offset += maxLimit
			continue
		}

		// If this request returned no data, increment empty counter
		consecutiveEmptyResponses++
		if consecutiveEmptyResponses >= 2 {
			// Two consecutive empty responses means we've fetched all data
			break
		}

		offset += maxLimit
	}

	// If we got some data, return it; otherwise use mock data
	if len(allCountries) > 0 {
		return allCountries, nil
	}

	return GetMockCountriesData(), nil
}

// FetchByName নাম দিয়ে দেশ খোঁজে
// Uses the new v5 API format with query parameter and Authorization header
func (c *CountriesClient) FetchByName(name string) ([]models.CountryResponse, error) {
	url := fmt.Sprintf("%s?q=%s", c.BaseURL, name)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// Failed to create request, search in mock data
		return searchMockCountriesByName(name), nil
	}

	// Add Authorization header
	if c.APIKey != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		// API call failed, search in mock data
		return searchMockCountriesByName(name), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Unexpected status code, search in mock data
		return searchMockCountriesByName(name), nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// Failed to read response, search in mock data
		return searchMockCountriesByName(name), nil
	}

	// Try to unmarshal as v5 API format
	var v5Response V5APIResponse
	if err := json.Unmarshal(body, &v5Response); err == nil && len(v5Response.Data.Objects) > 0 {
		return convertV5ToCountryResponse(v5Response.Data.Objects), nil
	}

	// Try to unmarshal as wrapper format
	var wrapper APIResponse
	if err := json.Unmarshal(body, &wrapper); err == nil && wrapper.Success && len(wrapper.Data) > 0 {
		return wrapper.Data, nil
	}

	// Fallback: try to unmarshal as direct array format
	var countries []models.CountryResponse
	if err := json.Unmarshal(body, &countries); err == nil && len(countries) > 0 {
		return countries, nil
	}

	// If all else fails, search in mock data
	return searchMockCountriesByName(name), nil
}

// FetchWithPagination fetches countries with pagination support
// limit: number of countries per page (1-100, default 25)
// offset: number of countries to skip (default 0)
func (c *CountriesClient) FetchWithPagination(limit int, offset int) ([]models.CountryResponse, error) {
	if limit < 1 || limit > 100 {
		limit = 25
	}
	if offset < 0 {
		offset = 0
	}

	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("offset", strconv.Itoa(offset))

	return c.fetchWithParams(params)
}

// FetchByRegion fetches countries filtered by region
// region: region name (e.g., "Europe", "Americas", "Africa", "Asia", "Oceania")
func (c *CountriesClient) FetchByRegion(region string) ([]models.CountryResponse, error) {
	if region == "" {
		return c.FetchAll()
	}

	params := url.Values{}
	params.Add("region", region)

	return c.fetchWithParams(params)
}

// FetchBySubregion fetches countries filtered by subregion
func (c *CountriesClient) FetchBySubregion(subregion string) ([]models.CountryResponse, error) {
	if subregion == "" {
		return c.FetchAll()
	}

	params := url.Values{}
	params.Add("subregion", subregion)

	return c.fetchWithParams(params)
}

// FetchWithSearch fetches countries matching search query
func (c *CountriesClient) FetchWithSearch(query string) ([]models.CountryResponse, error) {
	if query == "" {
		return c.FetchAll()
	}

	params := url.Values{}
	params.Add("q", query)

	return c.fetchWithParams(params)
}

// FetchWithFilters fetches countries with multiple filter options
type CountryFilters struct {
	Limit      int    // 1-100, default 25
	Offset     int    // default 0
	Region     string // filter by region
	Subregion  string // filter by subregion
	Search     string // search query
	FieldsOnly string // comma-separated fields to include
}

// FetchWithFilters fetches countries with advanced filtering
func (c *CountriesClient) FetchWithFilters(filters CountryFilters) ([]models.CountryResponse, error) {
	if filters.Limit < 1 || filters.Limit > 100 {
		filters.Limit = 25
	}
	if filters.Offset < 0 {
		filters.Offset = 0
	}

	params := url.Values{}
	params.Add("limit", strconv.Itoa(filters.Limit))
	params.Add("offset", strconv.Itoa(filters.Offset))

	if filters.Region != "" {
		params.Add("region", filters.Region)
	}
	if filters.Subregion != "" {
		params.Add("subregion", filters.Subregion)
	}
	if filters.Search != "" {
		params.Add("q", filters.Search)
	}
	if filters.FieldsOnly != "" {
		params.Add("response_fields", filters.FieldsOnly)
	}

	return c.fetchWithParams(params)
}

// fetchWithParams is the internal method that handles all API calls with query parameters
func (c *CountriesClient) fetchWithParams(params url.Values) ([]models.CountryResponse, error) {
	urlStr := c.BaseURL
	if len(params) > 0 {
		urlStr = fmt.Sprintf("%s?%s", c.BaseURL, params.Encode())
	}

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		// Failed to create request, use mock data
		return GetMockCountriesData(), nil
	}

	// Add Authorization header
	if c.APIKey != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	}

	resp, err := c.rawClient.Do(req)
	if err != nil {
		// API call failed, use mock data
		return GetMockCountriesData(), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Unexpected status code, use mock data
		return GetMockCountriesData(), nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// Failed to read response, use mock data
		return GetMockCountriesData(), nil
	}

	// Try to unmarshal as v5 API format
	var v5Response V5APIResponse
	if err := json.Unmarshal(body, &v5Response); err == nil && len(v5Response.Data.Objects) > 0 {
		return convertV5ToCountryResponse(v5Response.Data.Objects), nil
	}

	// Try to unmarshal as wrapper format
	var wrapper APIResponse
	if err := json.Unmarshal(body, &wrapper); err == nil && wrapper.Success && len(wrapper.Data) > 0 {
		return wrapper.Data, nil
	}

	// Fallback: try to unmarshal as direct array format
	var countries []models.CountryResponse
	if err := json.Unmarshal(body, &countries); err == nil && len(countries) > 0 {
		return countries, nil
	}

	// If all else fails, use mock data
	return GetMockCountriesData(), nil
}

// convertV5ToCountryResponse converts v5 API response to internal CountryResponse format
func convertV5ToCountryResponse(v5Countries []V5CountryData) []models.CountryResponse {
	result := make([]models.CountryResponse, 0, len(v5Countries))

	for _, v5 := range v5Countries {
		// Convert currencies list to map format
		currencies := make(map[string]models.Currency)
		for _, curr := range v5.Currencies {
			currencies[curr.Code] = models.Currency{
				Name:   curr.Name,
				Symbol: curr.Symbol,
			}
		}

		// Convert languages list to map format
		languages := make(map[string]string)
		for _, lang := range v5.Languages {
			// Use ISO 639-1 code as key, name as value
			if lang.Iso6391 != "" {
				languages[lang.Iso6391] = lang.Name
			}
		}

		// Extract capital name
		capitals := make([]string, 0)
		for _, cap := range v5.Capitals {
			if cap.Name != "" {
				capitals = append(capitals, cap.Name)
			}
		}

		// Extract latitude and longitude from first capital
		latLng := []float64{0.0, 0.0}
		if len(v5.Capitals) > 0 {
			latLng = []float64{v5.Capitals[0].Coordinates.Lat, v5.Capitals[0].Coordinates.Lng}
		}

		// Get flag PNG and SVG URLs
		flagPNG := v5.Flag.UrlPng
		flagSVG := v5.Flag.UrlSvg
		flagAlt := v5.Flag.Description

		country := models.CountryResponse{
			Name: models.CountryName{
				Common:   v5.Names.Common,
				Official: v5.Names.Official,
			},
			CCA2:       v5.Codes.Alpha2,
			CCA3:       v5.Codes.Alpha3,
			Capital:    capitals,
			Region:     v5.Region,
			Subregion:  v5.Subregion,
			Population: v5.Population,
			Flags: models.CountryFlag{
				PNG: flagPNG,
				SVG: flagSVG,
				Alt: flagAlt,
			},
			Currencies: currencies,
			Languages:  languages,
			LatLng:     latLng,
		}
		result = append(result, country)
	}

	return result
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
