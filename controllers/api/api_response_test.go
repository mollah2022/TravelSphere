package apicontrollers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "TravelSphere/routers"

	"github.com/beego/beego/v2/server/web"
)

// TestAPIResponseShape সব API response এর shape validate করে
func TestAPIResponseShape(t *testing.T) {
	setupTestServices()

	endpoints := []struct {
		method string
		url    string
	}{
		{"GET", "/api/countries"},
		{"GET", "/api/suggestions?q="},
		{"GET", "/api/attractions?lat=23.7&lon=90.4"},
	}

	for _, ep := range endpoints {
		r, _ := http.NewRequest(ep.method, ep.url, nil)
		w := httptest.NewRecorder()
		web.BeeApp.Handlers.ServeHTTP(w, r)

		// Content-Type JSON হতে হবে
		contentType := w.Header().Get("Content-Type")
		if contentType == "" {
			t.Logf("[WARN] %s %s: no Content-Type header", ep.method, ep.url)
		}

		// Valid JSON হতে হবে
		var resp map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Errorf("%s %s: invalid JSON: %v", ep.method, ep.url, err)
			continue
		}

		// success field থাকতে হবে
		if _, ok := resp["success"]; !ok {
			t.Errorf("%s %s: missing 'success' field", ep.method, ep.url)
		}
	}
}

// TestAPIErrorResponseShape error response এর shape validate করে
func TestAPIErrorResponseShape(t *testing.T) {
	setupTestServices()

	r, _ := http.NewRequest("GET", "/api/countries?region=InvalidRegion", nil)
	w := httptest.NewRecorder()
	web.BeeApp.Handlers.ServeHTTP(w, r)

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// Error response এ success: false থাকতে হবে
	if resp["success"] != false {
		t.Error("expected success: false for error response")
	}

	// error field থাকতে হবে
	if _, ok := resp["error"]; !ok {
		t.Error("missing 'error' field in error response")
	}

	// code field থাকতে হবে
	if _, ok := resp["code"]; !ok {
		t.Error("missing 'code' field in error response")
	}
}
