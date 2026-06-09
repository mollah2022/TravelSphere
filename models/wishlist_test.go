package models_test

import (
	"TravelSphere/models"
	"testing"
)

// TestWishlistStatusIsValid status validation test করে
func TestWishlistStatusIsValid(t *testing.T) {
	tests := []struct {
		status   models.WishlistStatus
		expected bool
	}{
		{models.StatusPlanned, true},
		{models.StatusVisited, true},
		{"planned", false},   // lowercase invalid
		{"visited", false},   // lowercase invalid
		{"", false},          // empty invalid
		{"Random", false},    // random string invalid
	}

	for _, tt := range tests {
		result := tt.status.IsValid()
		if result != tt.expected {
			t.Errorf("IsValid(%q) = %v, want %v", tt.status, result, tt.expected)
		}
	}
}

// TestCreateWishlistRequestValidate create request validation test
func TestCreateWishlistRequestValidate(t *testing.T) {
	// country_name missing হলে error আসবে
	req := &models.CreateWishlistRequest{
		CountryName: "",
		Status:      "Planned",
	}
	if err := req.Validate(); err != models.ErrCountryNameRequired {
		t.Errorf("expected ErrCountryNameRequired, got %v", err)
	}

	// invalid status হলে error আসবে
	req2 := &models.CreateWishlistRequest{
		CountryName: "Bangladesh",
		Status:      "invalid",
	}
	if err := req2.Validate(); err != models.ErrInvalidStatus {
		t.Errorf("expected ErrInvalidStatus, got %v", err)
	}

	// empty status হলে default Planned set হবে
	req3 := &models.CreateWishlistRequest{
		CountryName: "Bangladesh",
		Status:      "",
	}
	if err := req3.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if req3.Status != string(models.StatusPlanned) {
		t.Errorf("expected default status Planned, got %v", req3.Status)
	}

	// valid request
	req4 := &models.CreateWishlistRequest{
		CountryName: "France",
		Status:      "Visited",
	}
	if err := req4.Validate(); err != nil {
		t.Errorf("unexpected error for valid request: %v", err)
	}
}

// TestUpdateWishlistRequestValidate update request validation test
func TestUpdateWishlistRequestValidate(t *testing.T) {
	// empty status
	req := &models.UpdateWishlistRequest{Status: ""}
	if err := req.Validate(); err != models.ErrStatusRequired {
		t.Errorf("expected ErrStatusRequired, got %v", err)
	}

	// invalid status
	req2 := &models.UpdateWishlistRequest{Status: "Maybe"}
	if err := req2.Validate(); err != models.ErrInvalidStatus {
		t.Errorf("expected ErrInvalidStatus, got %v", err)
	}

	// valid
	req3 := &models.UpdateWishlistRequest{Status: "Visited"}
	if err := req3.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestNewSuccessResponse success response struct test
func TestNewSuccessResponse(t *testing.T) {
	data := map[string]string{"key": "value"}
	resp := models.NewSuccessResponse(data, "ok")

	if !resp.Success {
		t.Error("expected Success to be true")
	}
	if resp.Message != "ok" {
		t.Errorf("expected message 'ok', got %q", resp.Message)
	}
}

// TestNewErrorResponse error response struct test
func TestNewErrorResponse(t *testing.T) {
	resp := models.NewErrorResponse("not found", 404)

	if resp.Success {
		t.Error("expected Success to be false")
	}
	if resp.Code != 404 {
		t.Errorf("expected code 404, got %d", resp.Code)
	}
	if resp.Error != "not found" {
		t.Errorf("expected error 'not found', got %q", resp.Error)
	}
}