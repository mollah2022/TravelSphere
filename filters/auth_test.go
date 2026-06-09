package filters

import (
	"context"
	"net/http"
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
