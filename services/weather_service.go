package services

import (
	"TravelSphere/models"
	"TravelSphere/utils"
)

// WeatherServiceInterface mock করার জন্য interface
type WeatherServiceInterface interface {
	GetWeather(city string) *models.WeatherDTO
}

// WeatherService weather related business logic
type WeatherService struct {
	client *utils.WeatherClient
}

// NewWeatherService নতুন WeatherService তৈরি করে
func NewWeatherService(client *utils.WeatherClient) *WeatherService {
	return &WeatherService{client: client}
}

// GetWeather city এর weather আনে
// API key না থাকলে বা fail হলে unavailable DTO দেয়
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