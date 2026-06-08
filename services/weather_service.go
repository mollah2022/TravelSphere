package services

import (
	"TravelSphere/models"
	"TravelSphere/utils"
	"log"
)

// WeatherService handles weather-related business logic
// It communicates with Weather API client to get live weather data
type WeatherService struct {
	client *utils.WeatherClient
}

// NewWeatherService creates a new WeatherService instance
func NewWeatherService(client *utils.WeatherClient) *WeatherService {
	return &WeatherService{client: client}
}

// GetWeather returns current weather for a given city
func (s *WeatherService) GetWeather(city string) *models.WeatherDTO {
	if city == "" {
		return &models.WeatherDTO{Available: false}
	}

	dto, err := s.client.GetCurrent(city)
	if err != nil {
		log.Printf("[WARN] WeatherService: %v", err)
		return &models.WeatherDTO{Available: false}
	}

	return dto
}
