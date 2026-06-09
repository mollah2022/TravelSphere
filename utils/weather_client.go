package utils

import (
	"TravelSphere/models"
	"fmt"
	"os"
)

// WeatherClientInterface mock করার জন্য interface
type WeatherClientInterface interface {
	FetchCurrentWeather(city string) (*models.WeatherDTO, error)
}

// WeatherClient WeatherAPI client
type WeatherClient struct {
	BaseURL    string
	APIKey     string
	HTTPClient HTTPClient
}

// NewWeatherClient নতুন WeatherAPI client তৈরি করে
func NewWeatherClient() *WeatherClient {
	baseURL := os.Getenv("WEATHER_API_BASE_URL")
	if baseURL == "" {
		baseURL = "http://api.weatherapi.com/v1"
	}
	apiKey := os.Getenv("WEATHER_API_KEY")

	return &WeatherClient{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		HTTPClient: NewHTTPClient(8),
	}
}

// IsAvailable check করে weather API key আছে কিনা
func (c *WeatherClient) IsAvailable() bool {
	return c.APIKey != ""
}

// FetchCurrentWeather city name দিয়ে current weather fetch করে
func (c *WeatherClient) FetchCurrentWeather(city string) (*models.WeatherDTO, error) {
	// API key না থাকলে unavailable return করো
	if !c.IsAvailable() {
		return &models.WeatherDTO{Available: false}, nil
	}

	url := fmt.Sprintf("%s/current.json?key=%s&q=%s&aqi=no",
		c.BaseURL, c.APIKey, city)

	var response models.WeatherResponse
	if err := FetchJSON(c.HTTPClient, url, &response); err != nil {
		// Weather fail হলে app crash করবে না — graceful fallback
		return &models.WeatherDTO{Available: false}, nil
	}

	return &models.WeatherDTO{
		Location:   response.Location.Name,
		Country:    response.Location.Country,
		TempC:      response.Current.TempC,
		Condition:  response.Current.Condition.Text,
		Icon:       "https:" + response.Current.Condition.Icon,
		Humidity:   response.Current.Humidity,
		WindKph:    response.Current.WindKph,
		FeelsLikeC: response.Current.FeelsLikeC,
		Available:  true,
	}, nil
}
