package apicontrollers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"TravelSphere/models"
	"TravelSphere/services"
	"TravelSphere/store"
	"TravelSphere/utils"

	"github.com/beego/beego/v2/server/web"
	ctxpkg "github.com/beego/beego/v2/server/web/context"
)

type mockSessionStore struct {
	vals map[interface{}]interface{}
}

func (m *mockSessionStore) Set(ctx context.Context, key, value interface{}) error {
	if m.vals == nil {
		m.vals = map[interface{}]interface{}{}
	}
	m.vals[key] = value
	return nil
}

func (m *mockSessionStore) Get(ctx context.Context, key interface{}) interface{} {
	if m.vals == nil {
		return nil
	}
	return m.vals[key]
}

func (m *mockSessionStore) Delete(ctx context.Context, key interface{}) error {
	delete(m.vals, key)
	return nil
}

func (m *mockSessionStore) SessionID(ctx context.Context) string {
	return "mock-session"
}

func (m *mockSessionStore) SessionReleaseIfPresent(ctx context.Context, w http.ResponseWriter) {}
func (m *mockSessionStore) SessionRelease(ctx context.Context, w http.ResponseWriter)          {}
func (m *mockSessionStore) Flush(ctx context.Context) error {
	m.vals = map[interface{}]interface{}{}
	return nil
}

type mockCountryClient struct {
	data []models.CountryResponse
	err  error
}

func (m *mockCountryClient) FetchAll() ([]models.CountryResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.data, nil
}

func (m *mockCountryClient) FetchByName(name string) ([]models.CountryResponse, error) {
	return m.FetchAll()
}

type mockAttractionAPIClient struct {
	data []models.AttractionDTO
	err  error
}

func (m *mockAttractionAPIClient) FetchAttractionsByCoords(lat, lon float64, radius int) ([]models.AttractionDTO, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.data, nil
}

type mockWeatherHTTPClient struct{}

func (m *mockWeatherHTTPClient) Get(url string) (*http.Response, error) {
	return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
}

func (m *mockWeatherHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
}

func newAPIContext(method, url string, body []byte) (*web.Controller, *ctxpkg.Context, *httptest.ResponseRecorder, *http.Request) {
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(method, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := ctxpkg.NewContext()
	ctx.Reset(rw, req)
	if len(body) > 0 {
		ctx.Input.RequestBody = body
	}
	ctrl := &web.Controller{Ctx: ctx, Data: make(map[interface{}]interface{})}
	return ctrl, ctx, rw, req
}

func setCookie(req *http.Request, username string) {
	req.AddCookie(&http.Cookie{
		Name:  "travelsphere_user",
		Value: username,
	})
}

func setupAPIServices(countryData []models.CountryResponse, countryErr error, attractionErr error) {
	store := store.NewWishlistStore()
	services.Container = &services.ServiceContainer{
		CountryService:    services.NewCountryService(&mockCountryClient{data: countryData, err: countryErr}),
		AttractionService: services.NewAttractionService(&mockAttractionAPIClient{data: []models.AttractionDTO{{Name: "Lalbagh Fort", Kinds: "historic"}}, err: attractionErr}),
		WishlistService:   services.NewWishlistService(store),
		DashboardService:  services.NewDashboardService(store),
		WeatherService:    services.NewWeatherService(&utils.WeatherClient{APIKey: "", HTTPClient: &mockWeatherHTTPClient{}}),
	}
}

func decodeJSONResponse(t *testing.T, body *bytes.Buffer) map[string]interface{} {
	t.Helper()
	var resp map[string]interface{}
	if err := json.NewDecoder(body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}
	return resp
}

func TestCountriesAPIList_InvalidSearch(t *testing.T) {
	setupAPIServices([]models.CountryResponse{}, nil, nil)
	ctrl, _, _, _ := newAPIContext("GET", "/api/countries?search="+strings.Repeat("a", 101), nil)
	c := &CountriesAPIController{Controller: *ctrl}

	c.List()

	if ctrl.Ctx.ResponseWriter.Status != 400 {
		t.Fatalf("expected 400, got %d", ctrl.Ctx.ResponseWriter.Status)
	}
}

func TestCountriesAPIList_ServiceError(t *testing.T) {
	setupAPIServices(nil, errors.New("fail"), nil)
	ctrl, _, _, _ := newAPIContext("GET", "/api/countries?search=bang&region=Asia", nil)
	c := &CountriesAPIController{Controller: *ctrl}

	c.List()

	if ctrl.Ctx.ResponseWriter.Status != 500 {
		t.Fatalf("expected 500, got %d", ctrl.Ctx.ResponseWriter.Status)
	}
}

func TestCountriesAPIDetail_NotFound(t *testing.T) {
	setupAPIServices([]models.CountryResponse{}, nil, nil)
	ctrl, ctx, _, _ := newAPIContext("GET", "/api/countries/missing", nil)
	ctx.Input.SetParam(":slug", "missing")
	c := &CountriesAPIController{Controller: *ctrl}

	c.Detail()

	if ctrl.Ctx.ResponseWriter.Status != 404 {
		t.Fatalf("expected 404, got %d", ctrl.Ctx.ResponseWriter.Status)
	}
}

func TestCountriesAPIAttractions_InvalidCoordinates(t *testing.T) {
	setupAPIServices([]models.CountryResponse{}, nil, nil)
	ctrl, _, _, _ := newAPIContext("GET", "/api/attractions?lat=abc&lon=xyz", nil)
	c := &CountriesAPIController{Controller: *ctrl}

	c.Attractions()

	if ctrl.Ctx.ResponseWriter.Status != 400 {
		t.Fatalf("expected 400 for invalid coords, got %d", ctrl.Ctx.ResponseWriter.Status)
	}
}

func TestCountriesAPIAttractions_ServiceError(t *testing.T) {
	setupAPIServices([]models.CountryResponse{{Name: models.CountryName{Common: "Bangladesh"}}}, nil, errors.New("fail"))
	ctrl, _, _, _ := newAPIContext("GET", "/api/attractions?lat=23.7&lon=90.4", nil)
	c := &CountriesAPIController{Controller: *ctrl}

	c.Attractions()

	if ctrl.Ctx.ResponseWriter.Status != 200 {
		t.Fatalf("expected 200 when attraction service returns empty list on failure, got %d", ctrl.Ctx.ResponseWriter.Status)
	}
}

func TestCountriesAPIAttractions_Success(t *testing.T) {
	setupAPIServices([]models.CountryResponse{}, nil, nil)
	ctrl, _, _, _ := newAPIContext("GET", "/api/attractions?lat=23.7&lon=90.4", nil)
	c := &CountriesAPIController{Controller: *ctrl}

	c.Attractions()

	if ctrl.Ctx.ResponseWriter.Status != 200 {
		t.Fatalf("expected 200 for attractions success, got %d", ctrl.Ctx.ResponseWriter.Status)
	}
}

func TestCountriesAPISuggestions_QueryTooLong(t *testing.T) {
	setupAPIServices([]models.CountryResponse{}, nil, nil)
	ctrl, _, _, _ := newAPIContext("GET", "/api/suggestions?q="+strings.Repeat("a", 101), nil)
	c := &CountriesAPIController{Controller: *ctrl}

	c.Suggestions()

	if ctrl.Ctx.ResponseWriter.Status != 400 {
		t.Fatalf("expected 400 for long suggestion query, got %d", ctrl.Ctx.ResponseWriter.Status)
	}
}

func TestCountriesAPISuggestions_ServiceError(t *testing.T) {
	setupAPIServices(nil, errors.New("fail"), nil)
	ctrl, _, _, _ := newAPIContext("GET", "/api/suggestions?q=bang", nil)
	c := &CountriesAPIController{Controller: *ctrl}

	c.Suggestions()

	if ctrl.Ctx.ResponseWriter.Status != 500 {
		t.Fatalf("expected 500 on suggestions failure, got %d", ctrl.Ctx.ResponseWriter.Status)
	}
}

func TestDashboardAPISummary_Authorized(t *testing.T) {
	setupAPIServices([]models.CountryResponse{}, nil, nil)
	store := store.NewWishlistStore()
	store.Create("alice", "France", "Visit", string(models.StatusPlanned))
	services.Container.DashboardService = services.NewDashboardService(store)

	ctrl, _, rw, req := newAPIContext("GET", "/api/dashboard/summary", nil)
	setCookie(req, "alice")
	c := &DashboardAPIController{CountriesAPIController: CountriesAPIController{Controller: *ctrl}}

	c.Summary()

	if rw.Code != 200 {
		t.Fatalf("expected 200 authorized summary, got %d", rw.Code)
	}
	resp := decodeJSONResponse(t, rw.Body)
	if resp["success"] != true {
		t.Fatal("expected success true")
	}
}

func TestDashboardAPISummary_UsernameMissing(t *testing.T) {
	setupAPIServices([]models.CountryResponse{}, nil, nil)
	ctrl, _, rw, _ := newAPIContext("GET", "/api/dashboard/summary", nil)
	c := &DashboardAPIController{CountriesAPIController: CountriesAPIController{Controller: *ctrl}}

	c.Summary()

	if rw.Code != 401 {
		t.Fatalf("expected 401 when username missing, got %d", rw.Code)
	}
}

func TestWishlistAPIGetUsername_Authorized(t *testing.T) {
	setupAPIServices([]models.CountryResponse{}, nil, nil)
	ctrl, _, _, req := newAPIContext("GET", "/api/wishlist", nil)
	setCookie(req, "alice")
	c := &WishlistAPIController{CountriesAPIController: CountriesAPIController{Controller: *ctrl}}

	username, ok := c.getUsername()
	if !ok || username != "alice" {
		t.Fatalf("expected alice authorized, got %q ok=%v", username, ok)
	}
}

func TestWishlistAPIUpdate_StatusRequired(t *testing.T) {
	setupAPIServices([]models.CountryResponse{}, nil, nil)
	store := store.NewWishlistStore()
	item := store.Create("alice", "France", "Old note", string(models.StatusPlanned))
	services.Container.WishlistService = services.NewWishlistService(store)

	ctrl, ctx, _, req := newAPIContext("PUT", "/api/wishlist/"+item.ID, []byte(`{"note":"Updated note","status":""}`))
	setCookie(req, "alice")
	ctx.Input.SetParam(":id", item.ID)
	c := &WishlistAPIController{CountriesAPIController: CountriesAPIController{Controller: *ctrl}}

	c.Update()

	if ctrl.Ctx.ResponseWriter.Status != 400 {
		t.Fatalf("expected 400 status required, got %d", ctrl.Ctx.ResponseWriter.Status)
	}
}

func TestDashboardAPISummary_Unauthorized(t *testing.T) {
	setupAPIServices([]models.CountryResponse{}, nil, nil)
	ctrl, _, rw, _ := newAPIContext("GET", "/api/dashboard/summary", nil)
	c := &DashboardAPIController{CountriesAPIController: CountriesAPIController{Controller: *ctrl}}

	c.Summary()

	if rw.Code != 401 {
		t.Fatalf("expected 401 unauthorized, got %d", rw.Code)
	}
}

func TestWishlistAPIUpdate_InvalidBody(t *testing.T) {
	setupAPIServices([]models.CountryResponse{}, nil, nil)
	ctrl, ctx, _, req := newAPIContext("PUT", "/api/wishlist/id-123", []byte(`{"note":"New note",`))
	setCookie(req, "alice")
	ctx.Input.SetParam(":id", "id-123")
	c := &WishlistAPIController{CountriesAPIController: CountriesAPIController{Controller: *ctrl}}

	c.Update()

	if ctrl.Ctx.ResponseWriter.Status != 400 {
		t.Fatalf("expected 400 invalid body, got %d", ctrl.Ctx.ResponseWriter.Status)
	}
}

func TestWishlistAPIGetUsername_Unauthorized(t *testing.T) {
	setupAPIServices([]models.CountryResponse{}, nil, nil)
	ctrl, _, _, _ := newAPIContext("GET", "/api/wishlist", nil)
	c := &WishlistAPIController{CountriesAPIController: CountriesAPIController{Controller: *ctrl}}

	_, ok := c.getUsername()
	if ok {
		t.Fatal("expected unauthorized getUsername to return false")
	}
}

func TestWishlistAPICreate_InvalidStatus(t *testing.T) {
	setupAPIServices([]models.CountryResponse{}, nil, nil)
	ctrl, _, _, req := newAPIContext("POST", "/api/wishlist", []byte(`{"country_name":"France","status":"planned"}`))
	setCookie(req, "alice")
	c := &WishlistAPIController{CountriesAPIController: CountriesAPIController{Controller: *ctrl}}

	c.Create()

	if ctrl.Ctx.ResponseWriter.Status != 400 {
		t.Fatalf("expected 400 invalid status, got %d", ctrl.Ctx.ResponseWriter.Status)
	}
}

func TestWishlistAPIDelete_NotFound(t *testing.T) {
	setupAPIServices([]models.CountryResponse{}, nil, nil)
	ctrl, _, _, req := newAPIContext("DELETE", "/api/wishlist/missing", nil)
	setCookie(req, "alice")
	ctrl.Ctx.Input.SetParam(":id", "missing")
	c := &WishlistAPIController{CountriesAPIController: CountriesAPIController{Controller: *ctrl}}

	c.Delete()

	if ctrl.Ctx.ResponseWriter.Status != 403 {
		t.Fatalf("expected 403 unauthorized for missing item, got %d", ctrl.Ctx.ResponseWriter.Status)
	}
}

func TestWishlistAPIList_Authorized(t *testing.T) {
	setupAPIServices([]models.CountryResponse{}, nil, nil)
	store := store.NewWishlistStore()
	store.Create("alice", "France", "Visit", string(models.StatusPlanned))
	services.Container.WishlistService = services.NewWishlistService(store)

	ctrl, _, _, req := newAPIContext("GET", "/api/wishlist", nil)
	setCookie(req, "alice")
	c := &WishlistAPIController{CountriesAPIController: CountriesAPIController{Controller: *ctrl}}

	c.List()

	if ctrl.Ctx.ResponseWriter.Status != 200 {
		t.Fatalf("expected 200 for wishlist list, got %d", ctrl.Ctx.ResponseWriter.Status)
	}
}

func TestWishlistAPICreate_Success(t *testing.T) {
	setupAPIServices([]models.CountryResponse{}, nil, nil)
	ctrl, _, _, req := newAPIContext("POST", "/api/wishlist", []byte(`{"country_name":"France","note":"Visit Eiffel Tower","status":"Planned"}`))
	setCookie(req, "alice")
	c := &WishlistAPIController{CountriesAPIController: CountriesAPIController{Controller: *ctrl}}

	c.Create()

	if ctrl.Ctx.ResponseWriter.Status != 201 {
		t.Fatalf("expected 201 created, got %d", ctrl.Ctx.ResponseWriter.Status)
	}
}

func TestWishlistAPICreate_InvalidBody(t *testing.T) {
	setupAPIServices([]models.CountryResponse{}, nil, nil)
	ctrl, _, _, req := newAPIContext("POST", "/api/wishlist", []byte(`{"country_name":"France",`))
	setCookie(req, "alice")
	c := &WishlistAPIController{CountriesAPIController: CountriesAPIController{Controller: *ctrl}}

	c.Create()

	if ctrl.Ctx.ResponseWriter.Status != 400 {
		t.Fatalf("expected 400 invalid body, got %d", ctrl.Ctx.ResponseWriter.Status)
	}
}

func TestWishlistAPICreate_ValidationError(t *testing.T) {
	setupAPIServices([]models.CountryResponse{}, nil, nil)
	ctrl, _, _, req := newAPIContext("POST", "/api/wishlist", []byte(`{"country_name":"","status":"Planned"}`))
	setCookie(req, "alice")
	c := &WishlistAPIController{CountriesAPIController: CountriesAPIController{Controller: *ctrl}}

	c.Create()

	if ctrl.Ctx.ResponseWriter.Status != 400 {
		t.Fatalf("expected 400 for validation error, got %d", ctrl.Ctx.ResponseWriter.Status)
	}
}

func TestWishlistAPIUpdate_InvalidID(t *testing.T) {
	setupAPIServices([]models.CountryResponse{}, nil, nil)
	ctrl, _, _, req := newAPIContext("PUT", "/api/wishlist/", []byte(`{"note":"New note","status":"Visited"}`))
	setCookie(req, "alice")
	c := &WishlistAPIController{CountriesAPIController: CountriesAPIController{Controller: *ctrl}}

	c.Update()

	if ctrl.Ctx.ResponseWriter.Status != 400 {
		t.Fatalf("expected 400 invalid id, got %d", ctrl.Ctx.ResponseWriter.Status)
	}
}

func TestWishlistAPIUpdate_Unauthorized(t *testing.T) {
	setupAPIServices([]models.CountryResponse{}, nil, nil)
	store := store.NewWishlistStore()
	item := store.Create("bob", "France", "Old note", string(models.StatusPlanned))
	services.Container.WishlistService = services.NewWishlistService(store)

	ctrl, _, _, req := newAPIContext("PUT", "/api/wishlist/"+item.ID, []byte(`{"note":"New note","status":"Visited"}`))
	setCookie(req, "alice")
	ctrl.Ctx.Input.SetParam(":id", item.ID)
	c := &WishlistAPIController{CountriesAPIController: CountriesAPIController{Controller: *ctrl}}

	c.Update()

	if ctrl.Ctx.ResponseWriter.Status != 403 {
		t.Fatalf("expected 403 unauthorized update, got %d", ctrl.Ctx.ResponseWriter.Status)
	}
}

func TestWishlistAPIUpdate_Success(t *testing.T) {
	setupAPIServices([]models.CountryResponse{}, nil, nil)
	store := store.NewWishlistStore()
	item := store.Create("alice", "France", "Old note", string(models.StatusPlanned))
	services.Container.WishlistService = services.NewWishlistService(store)

	ctrl, _, _, req := newAPIContext("PUT", "/api/wishlist/"+item.ID, []byte(`{"note":"Updated note","status":"Visited"}`))
	setCookie(req, "alice")
	ctrl.Ctx.Input.SetParam(":id", item.ID)
	c := &WishlistAPIController{CountriesAPIController: CountriesAPIController{Controller: *ctrl}}

	c.Update()

	if ctrl.Ctx.ResponseWriter.Status != 200 {
		t.Fatalf("expected 200 successful update, got %d", ctrl.Ctx.ResponseWriter.Status)
	}
}

func TestWishlistAPIDelete_InvalidID(t *testing.T) {
	setupAPIServices([]models.CountryResponse{}, nil, nil)
	ctrl, _, _, req := newAPIContext("DELETE", "/api/wishlist/invalid", nil)
	setCookie(req, "alice")
	c := &WishlistAPIController{CountriesAPIController: CountriesAPIController{Controller: *ctrl}}

	c.Delete()

	if ctrl.Ctx.ResponseWriter.Status != 400 {
		t.Fatalf("expected 400 invalid id, got %d", ctrl.Ctx.ResponseWriter.Status)
	}
}

func TestWishlistAPIDelete_Unauthorized(t *testing.T) {
	setupAPIServices([]models.CountryResponse{}, nil, nil)
	store := store.NewWishlistStore()
	item := store.Create("bob", "France", "Note", string(models.StatusPlanned))
	services.Container.WishlistService = services.NewWishlistService(store)

	ctrl, _, _, req := newAPIContext("DELETE", "/api/wishlist/"+item.ID, nil)
	setCookie(req, "alice")
	ctrl.Ctx.Input.SetParam(":id", item.ID)
	c := &WishlistAPIController{CountriesAPIController: CountriesAPIController{Controller: *ctrl}}

	c.Delete()

	if ctrl.Ctx.ResponseWriter.Status != 403 {
		t.Fatalf("expected 403 unauthorized delete, got %d", ctrl.Ctx.ResponseWriter.Status)
	}
}

func TestWishlistAPIDelete_Success(t *testing.T) {
	setupAPIServices([]models.CountryResponse{}, nil, nil)
	store := store.NewWishlistStore()
	item := store.Create("alice", "France", "Note", string(models.StatusPlanned))
	services.Container.WishlistService = services.NewWishlistService(store)

	ctrl, _, _, req := newAPIContext("DELETE", "/api/wishlist/"+item.ID, nil)
	setCookie(req, "alice")
	ctrl.Ctx.Input.SetParam(":id", item.ID)
	c := &WishlistAPIController{CountriesAPIController: CountriesAPIController{Controller: *ctrl}}

	c.Delete()

	if ctrl.Ctx.ResponseWriter.Status != 200 {
		t.Fatalf("expected 200 successful delete, got %d", ctrl.Ctx.ResponseWriter.Status)
	}
}
