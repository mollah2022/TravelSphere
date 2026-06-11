package utils_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"TravelSphere/utils"
)

type mockHTTPClient struct {
	resp *http.Response
	err  error
}

func (m *mockHTTPClient) Get(url string) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.resp, nil
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.resp, nil
}

type errReader struct{}

func (errReader) Read([]byte) (int, error) {
	return 0, errors.New("read failure")
}

func (errReader) Close() error {
	return nil
}

func TestFetchJSON_HTTPError(t *testing.T) {
	client := &mockHTTPClient{err: errors.New("network down")}
	var target []interface{}
	if err := utils.FetchJSON(client, "http://example.test", &target); err == nil {
		t.Fatal("expected error for http failure")
	}
}

func TestFetchJSON_StatusCodeNotOK(t *testing.T) {
	client := &mockHTTPClient{resp: &http.Response{StatusCode: 500, Body: io.NopCloser(strings.NewReader("error"))}}
	var target []interface{}
	if err := utils.FetchJSON(client, "http://example.test", &target); err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestFetchJSON_ReadError(t *testing.T) {
	client := &mockHTTPClient{resp: &http.Response{StatusCode: 200, Body: errReader{}}}
	var target []interface{}
	if err := utils.FetchJSON(client, "http://example.test", &target); err == nil {
		t.Fatal("expected error for read failure")
	}
}

func TestFetchJSON_InvalidJSON(t *testing.T) {
	client := &mockHTTPClient{resp: &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("not-json"))}}
	var target []interface{}
	if err := utils.FetchJSON(client, "http://example.test", &target); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestNewHTTPClient(t *testing.T) {
	client := utils.NewHTTPClient(1)
	if client == nil {
		t.Fatal("expected non-nil HTTP client")
	}
}

func TestHTTPClient_Get(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("hello"))
	}))
	defer testServer.Close()

	client := utils.NewHTTPClient(1)
	resp, err := client.Get(testServer.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}
