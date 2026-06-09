package utils

import (
	"TravelSphere/models"
	"fmt"
	"os"
)

// WeatherClientInterface is used to mock the weather client in tests.
type WeatherClientInterface interface {
	FetchCurrentWeather(city string) (*models.WeatherDTO, error)
}

// WeatherClient is a client for interacting with the Weather API.
type WeatherClient struct {
	BaseURL    string
	APIKey     string
	HTTPClient HTTPClient
}

// NewWeatherClient creates a new WeatherAPI client.
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

// IsAvailable checks whether the weather API key is available.
func (c *WeatherClient) IsAvailable() bool {
	return c.APIKey != ""
}

// FetchCurrentWeather fetches the current weather using the city name.
func (c *WeatherClient) FetchCurrentWeather(city string) (*models.WeatherDTO, error) {

	if !c.IsAvailable() {
		return &models.WeatherDTO{Available: false}, nil
	}

	url := fmt.Sprintf("%s/current.json?key=%s&q=%s&aqi=no",
		c.BaseURL, c.APIKey, city)

	var response models.WeatherResponse
	if err := FetchJSON(c.HTTPClient, url, &response); err != nil {

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
