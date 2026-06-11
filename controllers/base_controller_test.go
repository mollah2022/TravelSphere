package controllers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beego/beego/v2/server/web"
	ctxpkg "github.com/beego/beego/v2/server/web/context"
)

type mockStore struct {
	vals map[interface{}]interface{}
}

func (m *mockStore) Set(ctx context.Context, key, value interface{}) error {
	if m.vals == nil {
		m.vals = map[interface{}]interface{}{}
	}
	m.vals[key] = value
	return nil
}

func (m *mockStore) Get(ctx context.Context, key interface{}) interface{} {
	if m.vals == nil {
		return nil
	}
	return m.vals[key]
}

func (m *mockStore) Delete(ctx context.Context, key interface{}) error {
	delete(m.vals, key)
	return nil
}

func (m *mockStore) SessionID(ctx context.Context) string {
	return "mock-session"
}

func (m *mockStore) SessionReleaseIfPresent(ctx context.Context, w http.ResponseWriter) {}
func (m *mockStore) SessionRelease(ctx context.Context, w http.ResponseWriter)          {}
func (m *mockStore) Flush(ctx context.Context) error {
	m.vals = map[interface{}]interface{}{}
	return nil
}

func newBaseController(path string, username string) *BaseController {
	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", path, nil)
	if username != "" {
		req.AddCookie(&http.Cookie{Name: "travelsphere_user", Value: username})
	}
	ctx := ctxpkg.NewContext()
	ctx.Reset(rw, req)
	ctx.Input.CruSession = &mockStore{vals: map[interface{}]interface{}{}}

	return &BaseController{Controller: web.Controller{Ctx: ctx, Data: make(map[interface{}]interface{})}}
}

func TestPrepare_SetsNavigationAndLoginState(t *testing.T) {
	c := newBaseController("/countries", "alice")
	c.Prepare()

	if !c.IsLoggedIn {
		t.Fatal("expected user to be logged in")
	}
	if c.Username != "alice" {
		t.Fatalf("expected username alice, got %q", c.Username)
	}
	if c.Data["NavCountries"] != true {
		t.Fatalf("expected NavCountries true, got %#v", c.Data["NavCountries"])
	}
	if c.Data["AppName"] != "TravelSphere" {
		t.Fatalf("expected AppName TravelSphere, got %#v", c.Data["AppName"])
	}
}

func TestPrepare_NoSession_SetsDefaultValues(t *testing.T) {
	c := newBaseController("/unknown", "")
	c.Prepare()

	if c.IsLoggedIn {
		t.Fatal("expected not logged in")
	}
	if c.Data["Username"] != "" {
		t.Fatalf("expected empty username, got %#v", c.Data["Username"])
	}
	if c.Data["NavHome"] != false {
		t.Fatalf("expected NavHome false, got %#v", c.Data["NavHome"])
	}
}
