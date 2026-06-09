package apicontrollers_test

import (
	"TravelSphere/models"
	"TravelSphere/services"
	"TravelSphere/store"
	"TravelSphere/utils"
	"os"
	"path/filepath"

	"github.com/beego/beego/v2/server/web"
)

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

type MockAttractionClient struct{}

func (m *MockAttractionClient) FetchAttractionsByCoords(
	lat, lon float64, radius int,
) ([]models.AttractionDTO, error) {
	return []models.AttractionDTO{
		{Name: "Lalbagh Fort", Kinds: "historic"},
	}, nil
}
