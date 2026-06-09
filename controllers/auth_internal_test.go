package controllers

import (
	"os"
	"testing"
)

func TestSanitizeUsername(t *testing.T) {
	tests := []struct{ in, want string }{
		{"alice", "alice"},
		{"alice bob", "alicebob"},
		{"a!l@i#c$e%", "alice"},
		{"john_doe", "john_doe"},
		{"mary-jane", "mary-jane"},
		{"A_b-1!2", "A_b-12"},
	}
	for _, tt := range tests {
		if got := sanitizeUsername(tt.in); got != tt.want {
			t.Fatalf("sanitizeUsername(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestGetDemoCredentials(t *testing.T) {
	origUser := os.Getenv("DEMO_USERNAME")
	origPass := os.Getenv("DEMO_PASSWORD")
	defer func() {
		os.Setenv("DEMO_USERNAME", origUser)
		os.Setenv("DEMO_PASSWORD", origPass)
	}()

	os.Setenv("DEMO_USERNAME", "demoUser")
	os.Setenv("DEMO_PASSWORD", "demoPass")
	u, p := getDemoCredentials()
	if u != "demoUser" || p != "demoPass" {
		t.Fatalf("expected demoUser/demoPass, got %q/%q", u, p)
	}
}
