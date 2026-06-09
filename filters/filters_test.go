package filters_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"TravelSphere/filters"

	ctxpkg "github.com/beego/beego/v2/server/web/context"
)

// mockStore implements session.Store for tests.
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
func (m *mockStore) SessionID(ctx context.Context) string                               { return "mock" }
func (m *mockStore) SessionReleaseIfPresent(ctx context.Context, w http.ResponseWriter) {}
func (m *mockStore) SessionRelease(ctx context.Context, w http.ResponseWriter)          {}
func (m *mockStore) Flush(ctx context.Context) error {
	m.vals = map[interface{}]interface{}{}
	return nil
}

func TestAuthFilter_RedirectsWhenNoSession(t *testing.T) {
	req := httptest.NewRequest("GET", "/protected", nil)
	rw := httptest.NewRecorder()
	c := ctxpkg.NewContext()
	c.Reset(rw, req)

	// No session set
	filters.AuthFilter(c)

	res := rw.Result()
	if res.StatusCode != 302 {
		t.Fatalf("expected redirect status 302, got %d", res.StatusCode)
	}
	loc := res.Header.Get("Location")
	if loc == "" || loc != "/login?redirect=/protected" {
		t.Fatalf("unexpected Location header: %q", loc)
	}
}

func TestAuthFilter_AllowsWhenUsernamePresent(t *testing.T) {
	req := httptest.NewRequest("GET", "/protected", nil)
	rw := httptest.NewRecorder()
	c := ctxpkg.NewContext()
	c.Reset(rw, req)

	// install a mock session with username
	m := &mockStore{vals: map[interface{}]interface{}{"username": "alice"}}
	c.Input.CruSession = m

	filters.AuthFilter(c)

	res := rw.Result()
	// should not redirect
	if res.Header.Get("Location") != "" {
		t.Fatalf("did not expect redirect, got Location=%q", res.Header.Get("Location"))
	}
}

func TestLoggingFilter_NoPanic(t *testing.T) {
	req := httptest.NewRequest("GET", "/ping", nil)
	rw := httptest.NewRecorder()
	c := ctxpkg.NewContext()
	c.Reset(rw, req)

	// Should not panic and returns quickly
	done := make(chan struct{})
	go func() {
		filters.LoggingFilter(c)
		close(done)
	}()

	select {
	case <-done:
		// ok
	case <-time.After(50 * time.Millisecond):
		t.Fatal("LoggingFilter did not return in time")
	}
	// allow goroutine to finish its small sleep
	time.Sleep(5 * time.Millisecond)
}
