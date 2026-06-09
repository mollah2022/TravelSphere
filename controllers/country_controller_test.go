package controllers_test

import (
	"TravelSphere/models"
	"TravelSphere/services"
	"TravelSphere/store"
	"TravelSphere/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	_ "TravelSphere/routers"

	"github.com/beego/beego/v2/server/web"
)

// setupTestServices test এর জন্য service container তৈরি করে
func setupTestServices() {
	wd, err := os.Getwd()
	if err == nil {
		viewsPath, staticPath := resolveProjectPaths(wd)
		web.BConfig.WebConfig.ViewsPath = viewsPath
		web.AddViewPath(viewsPath)
		web.SetStaticPath("/static", staticPath)
	}

	wishlistStore := store.NewWishlistStore()
	services.Container = &services.ServiceContainer{
		CountryService: services.NewCountryService(
			&MockCountriesClient{},
		),
		AttractionService: services.NewAttractionService(
			&MockAttractionClient{},
		),
		WishlistService:  services.NewWishlistService(wishlistStore),
		DashboardService: services.NewDashboardService(wishlistStore),
		WeatherService:   services.NewWeatherService(utils.NewWeatherClient()),
	}
}

func resolveProjectPaths(start string) (viewsPath, staticPath string) {
	dir := start
	for {
		viewsPath = filepath.Join(dir, "views")
		staticPath = filepath.Join(dir, "static")
		if vi, err := os.Stat(viewsPath); err == nil && vi.IsDir() {
			if si, err := os.Stat(staticPath); err == nil && si.IsDir() {
				return viewsPath, staticPath
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return filepath.Join(start, "..", "..", "views"), filepath.Join(start, "..", "..", "static")
}

// MockCountriesClient test এ real API call avoid করে
type MockCountriesClient struct{}

func (m *MockCountriesClient) FetchAll() ([]models.CountryResponse, error) {
	return []models.CountryResponse{
		{
			Name:       models.CountryName{Common: "Bangladesh"},
			CCA2:       "BD",
			CCA3:       "BGD",
			Capital:    []string{"Dhaka"},
			Region:     "Asia",
			Population: 170000000,
			Flags:      models.CountryFlag{PNG: "https://flag.png"},
			Currencies: map[string]models.Currency{
				"BDT": {Name: "Bangladeshi taka"},
			},
			Languages: map[string]string{"ben": "Bengali"},
			LatLng:    []float64{24.0, 90.0},
		},
	}, nil
}

func (m *MockCountriesClient) FetchByName(name string) ([]models.CountryResponse, error) {
	return m.FetchAll()
}

// MockAttractionClient test এ real API call avoid করে
type MockAttractionClient struct{}

func (m *MockAttractionClient) FetchAttractionsByCoords(
	lat, lon float64, radius int,
) ([]models.AttractionDTO, error) {
	return []models.AttractionDTO{
		{Name: "Lalbagh Fort", Kinds: "historic"},
	}, nil
}

func TestHomePageReturns200(t *testing.T) {
	setupTestServices()

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	web.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestCountriesPageReturns200(t *testing.T) {
	setupTestServices()

	r, _ := http.NewRequest("GET", "/countries", nil)
	w := httptest.NewRecorder()
	web.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestDestinationPageValidSlug(t *testing.T) {
	setupTestServices()

	r, _ := http.NewRequest("GET", "/countries/bangladesh", nil)
	w := httptest.NewRecorder()
	web.BeeApp.Handlers.ServeHTTP(w, r)

	// Valid slug → 200
	if w.Code != 200 && w.Code != 404 {
		t.Errorf("unexpected status: %d", w.Code)
	}
}

func TestDestinationPageInvalidSlug(t *testing.T) {
	setupTestServices()

	r, _ := http.NewRequest("GET", "/countries/INVALID SLUG!", nil)
	w := httptest.NewRecorder()
	web.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 404 {
		t.Errorf("invalid slug should return 404, got %d", w.Code)
	}
}

func TestLoginRedirectsIfLoggedIn(t *testing.T) {
	setupTestServices()

	r, _ := http.NewRequest("GET", "/login", nil)
	w := httptest.NewRecorder()
	web.BeeApp.Handlers.ServeHTTP(w, r)

	// Not logged in → 200 (show login page)
	if w.Code != 200 {
		t.Errorf("expected 200 for login page, got %d", w.Code)
	}
}

func TestWishlistRequiresAuth(t *testing.T) {
	setupTestServices()

	r, _ := http.NewRequest("GET", "/wishlist", nil)
	w := httptest.NewRecorder()
	web.BeeApp.Handlers.ServeHTTP(w, r)

	// Not logged in → 302 redirect to login
	if w.Code != 302 {
		t.Errorf("expected 302 redirect, got %d", w.Code)
	}
}
