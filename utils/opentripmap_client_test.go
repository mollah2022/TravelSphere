package utils_test

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"TravelSphere/models"
	"TravelSphere/utils"
)

func TestTransformAttractions_SkipsUnnamedFeatures(t *testing.T) {
	resp := models.AttractionResponse{
		Features: []models.AttractionFeature{
			{Properties: models.AttractionProperties{XID: "1", Name: "Site A", Kinds: "historic", Dist: 3.5, Point: models.AttractionPoint{Lat: 1.2, Lon: 2.3}}},
			{Properties: models.AttractionProperties{XID: "2", Name: "", Kinds: "museums", Dist: 5.1, Point: models.AttractionPoint{Lat: 3.4, Lon: 4.5}}},
			{Name: "Fallback Name", Properties: models.AttractionProperties{XID: "3", Kinds: "culture", Dist: 2.2, Point: models.AttractionPoint{Lat: 5.6, Lon: 6.7}}},
		},
	}

	items := utils.TransformAttractions(resp)
	if len(items) != 2 {
		t.Fatalf("expected 2 attractions, got %d", len(items))
	}
	if items[0].Name != "Site A" || items[1].Name != "Fallback Name" {
		t.Fatalf("unexpected names: %v", items)
	}
}

func TestGetPopularAttractions_ReturnsDefaultList(t *testing.T) {
	items := utils.GetPopularAttractions()
	if len(items) != 6 {
		t.Fatalf("expected 6 popular attractions, got %d", len(items))
	}
}

func TestOpenTripMapClient_NoAPIKeyReturnsEmpty(t *testing.T) {
	orig := os.Getenv("OPENTRIPMAP_API_KEY")
	defer os.Setenv("OPENTRIPMAP_API_KEY", orig)
	os.Setenv("OPENTRIPMAP_API_KEY", "")

	client := utils.NewOpenTripMapClient()
	items, err := client.FetchAttractionsByCoords(23.7, 90.4, 1000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected empty list without API key, got %d", len(items))
	}
}

func TestNewOpenTripMapClient_UsesEnvironment(t *testing.T) {
	origKey := os.Getenv("OPENTRIPMAP_API_KEY")
	origBase := os.Getenv("OPENTRIPMAP_BASE_URL")
	defer os.Setenv("OPENTRIPMAP_API_KEY", origKey)
	defer os.Setenv("OPENTRIPMAP_BASE_URL", origBase)

	os.Setenv("OPENTRIPMAP_API_KEY", "test-key")
	os.Setenv("OPENTRIPMAP_BASE_URL", "https://mockapi.test")

	client := utils.NewOpenTripMapClient()
	if client.APIKey != "test-key" {
		t.Fatalf("expected APIKey=test-key, got %q", client.APIKey)
	}
	if client.BaseURL != "https://mockapi.test" {
		t.Fatalf("expected custom base URL, got %q", client.BaseURL)
	}
}

func TestFetchAttractionsByCoords_Success(t *testing.T) {
	resp := &http.Response{
		StatusCode: 200,
		Body: io.NopCloser(strings.NewReader(`{
			"features": [{
				"properties": {
					"xid": "1",
					"name": "Park",
					"kinds": "historic",
					"dist": 3.0,
					"point": {"lon": 90.0, "lat": 23.0}
				}
			}]
		}`)),
	}
	client := &utils.OpenTripMapClient{BaseURL: "http://example.test", APIKey: "key", HTTPClient: &mockHTTPClient{resp: resp}}
	items, err := client.FetchAttractionsByCoords(23.7, 90.4, 1000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 || items[0].Name != "Park" {
		t.Fatalf("unexpected result: %#v", items)
	}
}
