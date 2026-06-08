package services

import (
	"TravelSphere/models"
	"TravelSphere/utils"
)

// AttractionServiceInterface mock করার জন্য interface
type AttractionServiceInterface interface {
	GetAttractionsByCountry(lat, lon float64) ([]models.AttractionDTO, error)
	GetPopularAttractions() []models.PopularAttraction
}

// AttractionService attraction related business logic
type AttractionService struct {
	client utils.OpenTripMapClientInterface
}

// NewAttractionService নতুন AttractionService তৈরি করে
func NewAttractionService(client utils.OpenTripMapClientInterface) *AttractionService {
	return &AttractionService{client: client}
}

// GetAttractionsByCountry দেশের coordinates দিয়ে attractions আনে
func (s *AttractionService) GetAttractionsByCountry(lat, lon float64) ([]models.AttractionDTO, error) {
	// 50km radius এ attractions খোঁজো
	attractions, err := s.client.FetchAttractionsByCoords(lat, lon, 50000)
	if err != nil {
		// Attraction fail হলে empty list দাও, app crash করবে না
		return []models.AttractionDTO{}, nil
	}
	return attractions, nil
}

// GetPopularAttractions home page এর জন্য static popular attractions
func (s *AttractionService) GetPopularAttractions() []models.PopularAttraction {
	return utils.GetPopularAttractions()
}