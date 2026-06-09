package utils_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"TravelSphere/utils"
)

type mockWeatherHTTPClient struct {
	resp *http.Response
	err  error
}

func (m *mockWeatherHTTPClient) Get(url string) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.resp, nil
}

func TestWeatherClient_IsAvailable(t *testing.T) {
	c := &utils.WeatherClient{APIKey: "", HTTPClient: &mockWeatherHTTPClient{}}
	if c.IsAvailable() {
		t.Fatal("expected unavailable when API key is empty")
	}

	c.APIKey = "key"
	if !c.IsAvailable() {
		t.Fatal("expected available when API key is set")
	}
}

func TestFetchCurrentWeather_NoAPIKey(t *testing.T) {
	c := &utils.WeatherClient{APIKey: "", HTTPClient: &mockWeatherHTTPClient{}}
	weather, err := c.FetchCurrentWeather("Dhaka")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if weather.Available {
		t.Fatal("expected unavailable weather when no API key")
	}
}

func TestFetchCurrentWeather_HTTPFailureReturnsUnavailable(t *testing.T) {
	c := &utils.WeatherClient{APIKey: "key", HTTPClient: &mockWeatherHTTPClient{err: errors.New("boom")}}
	weather, err := c.FetchCurrentWeather("Dhaka")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if weather.Available {
		t.Fatal("expected unavailable weather on HTTP failure")
	}
}

func TestFetchCurrentWeather_Success(t *testing.T) {
	jsonBody := `{ "location": {"name":"Dhaka","country":"Bangladesh"}, "current": {"temp_c": 25.0, "condition": {"text":"Sunny","icon":"//icon.png"}, "humidity": 60, "wind_kph": 10.5, "feelslike_c": 26.0 } }`
	c := &utils.WeatherClient{
		APIKey: "key",
		HTTPClient: &mockWeatherHTTPClient{resp: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(jsonBody)),
		}},
	}
	weather, err := c.FetchCurrentWeather("Dhaka")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !weather.Available {
		t.Fatal("expected weather available")
	}
	if weather.Location != "Dhaka" || weather.Country != "Bangladesh" {
		t.Fatalf("unexpected weather payload: %#v", weather)
	}
	if weather.Icon != "https://icon.png" {
		t.Fatalf("expected full icon URL, got %q", weather.Icon)
	}
}
