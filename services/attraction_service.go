package services

import (
	"TravelSphere/models"
	"TravelSphere/utils"
	"strings"
)

// AttractionService handles all attraction-related business logic
// It uses OpenTripMap API client to fetch data
type AttractionService struct {
	client *utils.OpenTripMapClient
}

// NewAttractionService creates a new AttractionService instance
func NewAttractionService(client *utils.OpenTripMapClient) *AttractionService {
	return &AttractionService{client: client}
}

// GetAttractionsByCountry returns attractions based on latitude and longitude
func (s *AttractionService) GetAttractionsByCountry(lat, lon float64) ([]models.AttractionDTO, error) {
	attractions, err := s.client.GetAttractionsByRadius(lat, lon, 50000)
	if err != nil {
		return nil, err
	}

	for i := range attractions {
		if attractions[i].Kinds != "" {
			attractions[i].KindsList = strings.Split(attractions[i].Kinds, ",")
		}
	}

	return attractions, nil
}

// GetPopularAttractions returns a static list of famous attractions
// This is used when we want to show default/popular places
func (s *AttractionService) GetPopularAttractions() []models.PopularAttraction {
	return []models.PopularAttraction{
		{Name: "Eiffel Tower", Kinds: "architecture", Country: "France"},
		{Name: "Colosseum", Kinds: "historic", Country: "Italy"},
		{Name: "Machu Picchu", Kinds: "archaeological", Country: "Peru"},
		{Name: "Great Wall of China", Kinds: "historic", Country: "China"},
		{Name: "Taj Mahal", Kinds: "architecture", Country: "India"},
		{Name: "Sagrada Família", Kinds: "architecture", Country: "Spain"},
		{Name: "Angkor Wat", Kinds: "archaeological", Country: "Cambodia"},
		{Name: "Petra", Kinds: "archaeological", Country: "Jordan"},
	}
}
