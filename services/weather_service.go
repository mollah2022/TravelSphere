package services

import (
	"TravelSphere/models"
	"TravelSphere/utils"
)

// WeatherServiceInterface is used to mock the weather service in tests.
type WeatherServiceInterface interface {
	GetWeather(city string) *models.WeatherDTO
}

// WeatherService weather related business logic
type WeatherService struct {
	client *utils.WeatherClient
}

// NewWeatherService creates a new instance of WeatherService.
func NewWeatherService(client *utils.WeatherClient) *WeatherService {
	return &WeatherService{client: client}
}

// GetWeather returns weather information for a city.
// If the API key is missing or the request fails, it returns an unavailable DTO.
func (s *WeatherService) GetWeather(city string) *models.WeatherDTO {
	if city == "" {
		return &models.WeatherDTO{Available: false}
	}
	weather, _ := s.client.FetchCurrentWeather(city)
	if weather == nil {
		return &models.WeatherDTO{Available: false}
	}
	return weather
}
