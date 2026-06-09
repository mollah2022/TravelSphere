package controllers

import (
	"testing"

	"TravelSphere/models"
	"TravelSphere/services"
	"TravelSphere/store"

	"github.com/beego/beego/v2/server/web"
)

type testCountryClient struct {
	data []models.CountryResponse
}

func (m *testCountryClient) FetchAll() ([]models.CountryResponse, error) {
	return m.data, nil
}

func (m *testCountryClient) FetchByName(name string) ([]models.CountryResponse, error) {
	return m.data, nil
}

type testAttractionClient struct{}

func (m *testAttractionClient) FetchAttractionsByCoords(lat, lon float64, radius int) ([]models.AttractionDTO, error) {
	return []models.AttractionDTO{{Name: "Park", Kinds: "historic"}}, nil
}

func newTestServiceContainer() *services.ServiceContainer {
	store := store.NewWishlistStore()
	return &services.ServiceContainer{
		CountryService:    services.NewCountryService(&testCountryClient{data: []models.CountryResponse{{Name: models.CountryName{Common: "Bangladesh"}, CCA2: "BD", CCA3: "BGD", Capital: []string{"Dhaka"}, Region: "Asia", Population: 170000000, Flags: models.CountryFlag{PNG: "https://flag.png"}, Currencies: map[string]models.Currency{"BDT": {Name: "Taka"}}, Languages: map[string]string{"ben": "Bengali"}, LatLng: []float64{24.0, 90.0}}}}),
		AttractionService: services.NewAttractionService(&testAttractionClient{}),
		WishlistService:   services.NewWishlistService(store),
		DashboardService:  services.NewDashboardService(store),
	}
}

func TestHomeController_Get_SetsFeaturedAndPopular(t *testing.T) {
	services.Container = newTestServiceContainer()
	c := &HomeController{BaseController: BaseController{Controller: web.Controller{Data: make(map[interface{}]interface{})}}}
	c.Get()

	if c.TplName != "home.tpl" {
		t.Fatalf("expected home.tpl, got %q", c.TplName)
	}
	if c.Data["PageTitle"] != "Discover Your Next Destination" {
		t.Fatalf("unexpected PageTitle: %#v", c.Data["PageTitle"])
	}
}

func TestDashboardController_Get_SetsSummary(t *testing.T) {
	container := newTestServiceContainer()
	services.Container = container
	container.WishlistService.AddToWishlist("alice", models.CreateWishlistRequest{CountryName: "France", Status: "Planned"})
	c := &DashboardController{BaseController: BaseController{Controller: web.Controller{Data: make(map[interface{}]interface{})}, Username: "alice"}}
	c.Get()

	if c.TplName != "dashboard.tpl" {
		t.Fatalf("expected dashboard.tpl, got %q", c.TplName)
	}
	if c.Data["Summary"] == nil {
		t.Fatal("expected summary data")
	}
}

func TestWishlistController_Get_SetsWishlistItems(t *testing.T) {
	container := newTestServiceContainer()
	services.Container = container
	container.WishlistService.AddToWishlist("alice", models.CreateWishlistRequest{CountryName: "Japan", Status: "Planned"})
	c := &WishlistController{BaseController: BaseController{Controller: web.Controller{Data: make(map[interface{}]interface{})}, Username: "alice"}}
	c.Get()

	if c.TplName != "wishlist.tpl" {
		t.Fatalf("expected wishlist.tpl, got %q", c.TplName)
	}
	items, ok := c.Data["WishlistItems"].([]*models.WishlistItem)
	if !ok || len(items) != 1 {
		t.Fatalf("expected one wishlist item, got %#v", c.Data["WishlistItems"])
	}
}
