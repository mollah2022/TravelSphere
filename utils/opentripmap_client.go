package utils

import (
	"TravelSphere/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// OpenTripMapClient handles communication with the OpenTripMap API
type OpenTripMapClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// NewOpenTripMapClient creates and returns a new API client instance
func NewOpenTripMapClient() *OpenTripMapClient {
	return &OpenTripMapClient{
		baseURL: os.Getenv("OPENTRIPMAP_BASE_URL"),
		apiKey:  os.Getenv("OPENTRIPMAP_API_KEY"),
		client:  NewHTTPClient(10),
	}
}

// GetAttractionsByRadius fetches nearby attractions based on latitude, longitude, and radius
func (c *OpenTripMapClient) GetAttractionsByRadius(lat, lon float64, radius int) ([]models.AttractionDTO, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("OPENTRIPMAP_API_KEY not set")
	}

	url := fmt.Sprintf(
		"%s/places/radius?radius=%d&lon=%f&lat=%f&kinds=interesting_places&format=geojson&limit=10&apikey=%s",
		c.baseURL, radius, lon, lat, c.apiKey,
	)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("opentripmap request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("opentripmap returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result models.AttractionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse attractions: %w", err)
	}

	var dtos []models.AttractionDTO
	for _, f := range result.Features {
		if f.Properties.Name == "" {
			continue
		}
		// Map API data to internal DTO
		dto := models.AttractionDTO{
			XID:       f.Properties.XID,
			Name:      f.Properties.Name,
			Kinds:     f.Properties.Kinds,
			Distance:  f.Properties.Dist,
			Latitude:  f.Properties.Point.Lat,
			Longitude: f.Properties.Point.Lon,
		}
		dtos = append(dtos, dto)
	}

	return dtos, nil
}
