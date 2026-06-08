package utils

import (
	"TravelSphere/models"
	"fmt"
	"os"
)

// OpenTripMapClientInterface mock করার জন্য interface
type OpenTripMapClientInterface interface {
	FetchAttractionsByCoords(lat, lon float64, radius int) ([]models.AttractionDTO, error)
}

// OpenTripMapClient OpenTripMap API client
type OpenTripMapClient struct {
	BaseURL    string
	APIKey     string
	HTTPClient HTTPClient
}

// NewOpenTripMapClient নতুন OpenTripMap client তৈরি করে
func NewOpenTripMapClient() *OpenTripMapClient {
	baseURL := os.Getenv("OPENTRIPMAP_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.opentripmap.com/0.1/en"
	}
	apiKey := os.Getenv("OPENTRIPMAP_API_KEY")

	return &OpenTripMapClient{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		HTTPClient: NewHTTPClient(10),
	}
}

// FetchAttractionsByCoords coordinates দিয়ে nearby attractions fetch করে
func (c *OpenTripMapClient) FetchAttractionsByCoords(lat, lon float64, radius int) ([]models.AttractionDTO, error) {
	if c.APIKey == "" {
		return []models.AttractionDTO{}, nil
	}

	url := fmt.Sprintf(
		"%s/places/radius?radius=%d&lon=%f&lat=%f&kinds=interesting_places,museums,historic,tourist_facilities&limit=10&apikey=%s",
		c.BaseURL, radius, lon, lat, c.APIKey,
	)

	var response models.AttractionResponse
	if err := FetchJSON(c.HTTPClient, url, &response); err != nil {
		return nil, fmt.Errorf("FetchAttractionsByCoords failed: %w", err)
	}

	return TransformAttractions(response), nil
}

// TransformAttractions AttractionResponse কে []AttractionDTO তে convert করে
func TransformAttractions(resp models.AttractionResponse) []models.AttractionDTO {
	result := make([]models.AttractionDTO, 0, len(resp.Features))

	for _, f := range resp.Features {
		// নাম ছাড়া attraction skip করো
		name := f.Properties.Name
		if name == "" {
			name = f.Name
		}
		if name == "" {
			continue
		}

		dto := models.AttractionDTO{
			XID:       f.Properties.XID,
			Name:      name,
			Kinds:     f.Properties.Kinds,
			KindsList: FormatKinds(f.Properties.Kinds),
			Distance:  f.Properties.Dist,
			Latitude:  f.Properties.Point.Lat,
			Longitude: f.Properties.Point.Lon,
		}
		result = append(result, dto)
	}
	return result
}

// GetPopularAttractions home page এর জন্য static popular attractions
// Free API tier এ global attractions পাওয়া যায় না তাই static
func GetPopularAttractions() []models.PopularAttraction {
	return []models.PopularAttraction{
		{Name: "Eiffel Tower", Kinds: "architecture,historic", Country: "France"},
		{Name: "Grand Canyon", Kinds: "natural", Country: "USA"},
		{Name: "Sydney Opera House", Kinds: "architecture,theatre", Country: "Australia"},
		{Name: "Colosseum", Kinds: "historic,architecture", Country: "Italy"},
		{Name: "Taj Mahal", Kinds: "historic,architecture", Country: "India"},
		{Name: "Machu Picchu", Kinds: "historic,natural", Country: "Peru"},
	}
}