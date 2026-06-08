package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "TravelSphere/routers"

	"github.com/beego/beego/v2/server/web"
)

func TestAPICountriesList(t *testing.T) {
	setupTestServices()

	r, _ := http.NewRequest("GET", "/api/countries", nil)
	w := httptest.NewRecorder()
	web.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}

	// JSON response validate করো
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if resp["success"] != true {
		t.Error("expected success: true")
	}
}

func TestAPICountriesSearch(t *testing.T) {
	setupTestServices()

	r, _ := http.NewRequest("GET", "/api/countries?search=bangladesh", nil)
	w := httptest.NewRecorder()
	web.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAPICountriesInvalidRegion(t *testing.T) {
	setupTestServices()

	r, _ := http.NewRequest("GET", "/api/countries?region=Mars", nil)
	w := httptest.NewRecorder()
	web.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("expected 400 for invalid region, got %d", w.Code)
	}
}

func TestAPICountryDetail(t *testing.T) {
	setupTestServices()

	r, _ := http.NewRequest("GET", "/api/countries/bangladesh", nil)
	w := httptest.NewRecorder()
	web.BeeApp.Handlers.ServeHTTP(w, r)

	// Mock client এ Bangladesh আছে → 200
	if w.Code != 200 && w.Code != 404 {
		t.Errorf("unexpected status %d", w.Code)
	}
}

func TestAPICountryInvalidSlug(t *testing.T) {
	setupTestServices()

	r, _ := http.NewRequest("GET", "/api/countries/INVALID!", nil)
	w := httptest.NewRecorder()
	web.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("expected 400 for invalid slug, got %d", w.Code)
	}
}

func TestAPISuggestions(t *testing.T) {
	setupTestServices()

	r, _ := http.NewRequest("GET", "/api/suggestions?q=bang", nil)
	w := httptest.NewRecorder()
	web.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAPISuggestionsEmpty(t *testing.T) {
	setupTestServices()

	r, _ := http.NewRequest("GET", "/api/suggestions?q=", nil)
	w := httptest.NewRecorder()
	web.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("expected 200 for empty query, got %d", w.Code)
	}
}