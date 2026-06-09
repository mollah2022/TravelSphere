package services_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"TravelSphere/services"
	"TravelSphere/utils"
)

type mockHTTPClient struct {
	resp *http.Response
	err  error
}

func (m *mockHTTPClient) Get(url string) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.resp, nil
}

func TestGetWeather_EmptyCity(t *testing.T) {
	svc := services.NewWeatherService(&utils.WeatherClient{APIKey: "", HTTPClient: &mockHTTPClient{}})
	weather := svc.GetWeather("")
	if weather.Available {
		t.Fatal("expected unavailable when city is empty")
	}
}

func TestGetWeather_FallbackOnError(t *testing.T) {
	client := &utils.WeatherClient{APIKey: "test-key", HTTPClient: &mockHTTPClient{err: errors.New("failed")}}
	svc := services.NewWeatherService(client)
	weather := svc.GetWeather("Dhaka")
	if weather.Available {
		t.Fatal("expected unavailable when client returns error")
	}
}

func TestGetWeather_ReturnsAvailable(t *testing.T) {
	jsonBody := `{"location":{"name":"Dhaka","country":"Bangladesh"},"current":{"temp_c":25.0,"condition":{"text":"Sunny","icon":"//cdn.test/image.png"},"humidity":78,"wind_kph":10.0,"feelslike_c":26.0}}`
	resp := &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(jsonBody))}
	client := &utils.WeatherClient{APIKey: "test-key", HTTPClient: &mockHTTPClient{resp: resp}}
	svc := services.NewWeatherService(client)
	weather := svc.GetWeather("Dhaka")
	if !weather.Available {
		t.Fatal("expected available weather")
	}
	if weather.Location != "Dhaka" || weather.Country != "Bangladesh" {
		t.Fatalf("unexpected weather returned: %#v", weather)
	}
}
