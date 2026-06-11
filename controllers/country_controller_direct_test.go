package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"TravelSphere/models"
	"TravelSphere/services"
	"TravelSphere/utils"

	"github.com/beego/beego/v2/server/web"
	ctxpkg "github.com/beego/beego/v2/server/web/context"
)

type failingCountryService struct {
	all []models.CountryDTO
}

func (m *failingCountryService) GetAllCountries() ([]models.CountryDTO, error) {
	return nil, errors.New("fail")
}

func (m *failingCountryService) SearchCountries(search, region string) ([]models.CountryDTO, error) {
	return nil, errors.New("fail")
}

func (m *failingCountryService) GetCountryBySlug(slug string) (*models.CountryDTO, error) {
	return nil, errors.New("fail")
}

func (m *failingCountryService) GetFeaturedCountries() ([]models.CountryDTO, error) {
	return nil, errors.New("fail")
}

func newCountryController(path string) *CountryController {
	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", path, nil)
	ctx := ctxpkg.NewContext()
	ctx.Reset(rw, req)
	return &CountryController{BaseController: BaseController{Controller: web.Controller{Ctx: ctx, Data: make(map[interface{}]interface{})}}}
}

func TestCountryController_List_ShowsErrorOnServiceFailure(t *testing.T) {
	services.Container = &services.ServiceContainer{CountryService: services.NewCountryService(&failingCountryAPIClient{})}
	c := newCountryController("/countries")
	c.List()

	if c.TplName != "countries.tpl" {
		t.Fatalf("expected countries.tpl, got %q", c.TplName)
	}
	if c.Data["Error"] == nil {
		t.Fatal("expected error message when countries fail to load")
	}
}

func TestCountryController_Detail_InvalidSlug(t *testing.T) {
	services.Container = &services.ServiceContainer{
		CountryService:    services.NewCountryService(&failingCountryAPIClient{}),
		AttractionService: services.NewAttractionService(&fakeAttractionAPIClient{}),
		WeatherService:    services.NewWeatherService(&utils.WeatherClient{HTTPClient: &mockHTTPClient{}}),
	}
	c := newCountryController("/countries/INVALID!")
	c.Ctx.Input.SetParam(":slug", "INVALID!")
	c.Detail()

	if c.TplName != "error.tpl" {
		t.Fatalf("expected error.tpl for invalid slug, got %q", c.TplName)
	}
}

func TestCountryController_Detail_NotFound(t *testing.T) {
	services.Container = &services.ServiceContainer{
		CountryService:    services.NewCountryService(&emptyCountryAPIClient{}),
		AttractionService: services.NewAttractionService(&fakeAttractionAPIClient{}),
		WeatherService:    services.NewWeatherService(&utils.WeatherClient{HTTPClient: &mockHTTPClient{}}),
	}
	c := newCountryController("/countries/unknown")
	c.Ctx.Input.SetParam(":slug", "unknown")
	c.Detail()

	if c.TplName != "error.tpl" {
		t.Fatalf("expected error.tpl for missing country, got %q", c.TplName)
	}
}

func TestHomeController_Get_FeaturedServiceFailure(t *testing.T) {
	services.Container = &services.ServiceContainer{
		CountryService:    services.NewCountryService(&failingCountryAPIClient{}),
		AttractionService: services.NewAttractionService(&fakeAttractionAPIClient{}),
		WeatherService:    services.NewWeatherService(&utils.WeatherClient{HTTPClient: &mockHTTPClient{}}),
	}
	c := &HomeController{BaseController: BaseController{Controller: web.Controller{Ctx: ctxpkg.NewContext(), Data: make(map[interface{}]interface{})}}}
	c.Get()

	featured, ok := c.Data["FeaturedCountries"].([]models.CountryDTO)
	if !ok || len(featured) != 0 {
		t.Fatalf("expected empty featured countries slice, got %#v", c.Data["FeaturedCountries"])
	}
	if c.Data["PopularAttractions"] == nil {
		t.Fatal("expected PopularAttractions to be set")
	}
}

type fakeCountryService struct {
	all []models.CountryDTO
}

func (m *fakeCountryService) GetAllCountries() ([]models.CountryDTO, error) {
	return m.all, nil
}

func (m *fakeCountryService) SearchCountries(search, region string) ([]models.CountryDTO, error) {
	return m.all, nil
}

func (m *fakeCountryService) GetCountryBySlug(slug string) (*models.CountryDTO, error) {
	for _, c := range m.all {
		if c.Slug == slug {
			return &c, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *fakeCountryService) GetFeaturedCountries() ([]models.CountryDTO, error) {
	return m.all, nil
}

type fakeAttractionClient struct{}

func (m *fakeAttractionClient) FetchAttractionsByCoords(lat, lon float64, radius int) ([]models.AttractionDTO, error) {
	return []models.AttractionDTO{{Name: "Park", Kinds: "historic"}}, nil
}

type fakeWeatherClient struct{}

func (m *fakeWeatherClient) FetchCurrentWeather(city string) (*models.WeatherDTO, error) {
	return &models.WeatherDTO{Available: false}, nil
}

type failingCountryAPIClient struct{}

type emptyCountryAPIClient struct{}

type fakeAttractionAPIClient struct{}

type mockHTTPClient struct {
	resp *http.Response
	err  error
}

func (f *mockHTTPClient) Get(url string) (*http.Response, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.resp, nil
}

func (f *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.resp, nil
}

func (f *failingCountryAPIClient) FetchAll() ([]models.CountryResponse, error) {
	return nil, errors.New("fail")
}

func (f *failingCountryAPIClient) FetchByName(name string) ([]models.CountryResponse, error) {
	return nil, errors.New("fail")
}

func (f *emptyCountryAPIClient) FetchAll() ([]models.CountryResponse, error) {
	return []models.CountryResponse{}, nil
}

func (f *emptyCountryAPIClient) FetchByName(name string) ([]models.CountryResponse, error) {
	return []models.CountryResponse{}, nil
}

func (f *fakeAttractionAPIClient) FetchAttractionsByCoords(lat, lon float64, radius int) ([]models.AttractionDTO, error) {
	return []models.AttractionDTO{{Name: "Park", Kinds: "historic"}}, nil
}
