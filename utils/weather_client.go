package utils

import (
	"TravelSphere/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// WeatherClient handles communication with the Weather API
type WeatherClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// NewWeatherClient creates and returns a new WeatherClient instance
func NewWeatherClient() *WeatherClient {
	return &WeatherClient{
		baseURL: os.Getenv("WEATHER_API_BASE_URL"),
		apiKey:  os.Getenv("WEATHER_API_KEY"),
		client:  NewHTTPClient(8),
	}
}

// GetCurrent fetches current weather data for a given city
func (c *WeatherClient) GetCurrent(city string) (*models.WeatherDTO, error) {
	if c.apiKey == "" || c.baseURL == "" {
		return nil, fmt.Errorf("weather API not configured")
	}

	reqURL := fmt.Sprintf("%s/current.json?key=%s&q=%s&aqi=no",
		c.baseURL, c.apiKey, url.QueryEscape(city),
	)

	resp, err := c.client.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("weather request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read weather response: %w", err)
	}

	var raw models.WeatherResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse weather response: %w", err)
	}

	// Convert raw API response into clean DTO format
	return &models.WeatherDTO{
		Location:   raw.Location.Name,
		Country:    raw.Location.Country,
		TempC:      raw.Current.TempC,
		Condition:  raw.Current.Condition.Text,
		Icon:       raw.Current.Condition.Icon,
		Humidity:   raw.Current.Humidity,
		WindKph:    raw.Current.WindKph,
		FeelsLikeC: raw.Current.FeelsLikeC,
		Available:  true,
	}, nil
}
