package apicontrollers

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"TravelSphere/models"
	"TravelSphere/services"

	"github.com/beego/beego/v2/server/web"
	ctxpkg "github.com/beego/beego/v2/server/web/context"
)

type mockCountriesClient struct{}

func (m *mockCountriesClient) FetchAll() ([]models.CountryResponse, error) {
	return []models.CountryResponse{{Name: models.CountryName{Common: "Bangladesh"}, CCA2: "BD", CCA3: "BGD", Capital: []string{"Dhaka"}, Region: "Asia", Population: 170000000, Flags: models.CountryFlag{PNG: "https://flag.png"}, Currencies: map[string]models.Currency{"BDT": {Name: "Taka"}}, Languages: map[string]string{"ben": "Bengali"}, LatLng: []float64{24.0, 90.0}}}, nil
}

func (m *mockCountriesClient) FetchByName(name string) ([]models.CountryResponse, error) {
	return m.FetchAll()
}

type mockAttractionClient struct{}

func (m *mockAttractionClient) FetchAttractionsByCoords(lat, lon float64, radius int) ([]models.AttractionDTO, error) {
	return []models.AttractionDTO{{Name: "Park"}}, nil
}

func newAPIRequest(method, url string) (*CountriesAPIController, *httptest.ResponseRecorder) {
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(method, url, nil)
	ctx := ctxpkg.NewContext()
	ctx.Reset(rw, req)
	ctrl := &CountriesAPIController{Controller: web.Controller{Ctx: ctx}}
	ctrl.Data = make(map[interface{}]interface{})
	return ctrl, rw
}

func setupAPITestServices() {
	services.Container = &services.ServiceContainer{
		CountryService:    services.NewCountryService(&mockCountriesClient{}),
		AttractionService: services.NewAttractionService(&mockAttractionClient{}),
	}
}

func TestList_InvalidSearch(t *testing.T) {
	setupAPITestServices()
	ctrl, rw := newAPIRequest("GET", "/api/countries?search="+strings.Repeat("a", 101))
	ctrl.List()
	if rw.Code != 400 {
		t.Fatalf("expected 400, got %d", rw.Code)
	}
}

func TestList_InvalidRegion(t *testing.T) {
	setupAPITestServices()
	ctrl, rw := newAPIRequest("GET", "/api/countries?region=Mars")
	ctrl.List()
	if rw.Code != 400 {
		t.Fatalf("expected 400, got %d", rw.Code)
	}
}

func TestSuggestions_EmptyQuery(t *testing.T) {
	setupAPITestServices()
	ctrl, rw := newAPIRequest("GET", "/api/suggestions?q=")
	ctrl.Suggestions()
	if rw.Code != 200 {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(rw.Body.Bytes(), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestSuggestions_SearchQuery(t *testing.T) {
	setupAPITestServices()
	ctrl, rw := newAPIRequest("GET", "/api/suggestions?q=bang")
	ctrl.Suggestions()
	if rw.Code != 200 {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
}

func TestDetail_InvalidSlug(t *testing.T) {
	setupAPITestServices()
	ctrl, rw := newAPIRequest("GET", "/api/countries/INVALID!")
	ctrl.Ctx.Input.SetParam(":slug", "INVALID!")
	ctrl.Detail()
	if rw.Code != 400 {
		t.Fatalf("expected 400, got %d", rw.Code)
	}
}

func TestAttractions_InvalidCoordinates(t *testing.T) {
	setupAPITestServices()
	ctrl, rw := newAPIRequest("GET", "/api/attractions?lat=abc&lon=xyz")
	ctrl.Attractions()
	if rw.Code != 400 {
		t.Fatalf("expected 400, got %d", rw.Code)
	}
}
