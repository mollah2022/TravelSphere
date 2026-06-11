package services_test

import (
	"errors"
	"testing"

	"TravelSphere/models"
	"TravelSphere/services"
)

type failingAttractionClient struct{}

func (m *failingAttractionClient) FetchAttractionsByCoords(lat, lon float64, radius int) ([]models.AttractionDTO, error) {
	return nil, errors.New("failed")
}

type workingAttractionClient struct{}

func (m *workingAttractionClient) FetchAttractionsByCoords(lat, lon float64, radius int) ([]models.AttractionDTO, error) {
	return []models.AttractionDTO{{Name: "Park", Kinds: "natural"}}, nil
}

func TestGetAttractionsByCountry_ErrorReturnsEmpty(t *testing.T) {
	svc := services.NewAttractionService(&failingAttractionClient{})
	items, err := svc.GetAttractionsByCountry(23.7, 90.4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected empty slice on error, got %d", len(items))
	}
}

func TestGetAttractionsByCountry_Success(t *testing.T) {
	svc := services.NewAttractionService(&workingAttractionClient{})
	items, err := svc.GetAttractionsByCountry(23.7, 90.4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) == 0 {
		t.Fatalf("expected attractions, got 0")
	}
	if items[0].Name != "Park" {
		t.Errorf("expected Park, got %q", items[0].Name)
	}
}

func TestGetPopularAttractions_ReturnsStaticList(t *testing.T) {
	svc := services.NewAttractionService(&workingAttractionClient{})
	items := svc.GetPopularAttractions()
	if len(items) == 0 {
		t.Fatal("expected popular attractions list")
	}
}
