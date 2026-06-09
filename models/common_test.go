package models

import "testing"

func TestNewSuccessResponse(t *testing.T) {
	resp := NewSuccessResponse(map[string]string{"name": "Bangladesh"}, "Loaded")
	if !resp.Success {
		t.Fatal("expected success to be true")
	}
	data, ok := resp.Data.(map[string]string)
	if !ok || data["name"] != "Bangladesh" {
		t.Fatalf("unexpected data %#v", resp.Data)
	}
	if resp.Message != "Loaded" {
		t.Fatalf("unexpected message %q", resp.Message)
	}
}

func TestNewErrorResponse(t *testing.T) {
	err := NewErrorResponse("Not found", 404)
	if err.Success {
		t.Fatal("expected success to be false")
	}
	if err.Error != "Not found" {
		t.Fatalf("unexpected error message %q", err.Error)
	}
	if err.Code != 404 {
		t.Fatalf("unexpected code %d", err.Code)
	}
}

func TestWishlistStatus_IsValid(t *testing.T) {
	if !StatusPlanned.IsValid() {
		t.Fatal("expected Planned to be valid")
	}
	if !StatusVisited.IsValid() {
		t.Fatal("expected Visited to be valid")
	}
	if WishlistStatus("unknown").IsValid() {
		t.Fatal("expected unknown status to be invalid")
	}
}

func TestCreateWishlistRequestValidate(t *testing.T) {
	req := CreateWishlistRequest{CountryName: "France", Status: ""}
	if err := req.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Status != string(StatusPlanned) {
		t.Fatalf("expected default status Planned, got %q", req.Status)
	}

	req2 := CreateWishlistRequest{CountryName: "", Status: "Planned"}
	if err := req2.Validate(); err != ErrCountryNameRequired {
		t.Fatalf("expected ErrCountryNameRequired, got %v", err)
	}

	req3 := CreateWishlistRequest{CountryName: "Japan", Status: "Unknown"}
	if err := req3.Validate(); err != ErrInvalidStatus {
		t.Fatalf("expected ErrInvalidStatus, got %v", err)
	}
}

func TestUpdateWishlistRequestValidate(t *testing.T) {
	req := UpdateWishlistRequest{Status: "Visited"}
	if err := req.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req2 := UpdateWishlistRequest{Status: ""}
	if err := req2.Validate(); err != ErrStatusRequired {
		t.Fatalf("expected ErrStatusRequired, got %v", err)
	}

	req3 := UpdateWishlistRequest{Status: "Unknown"}
	if err := req3.Validate(); err != ErrInvalidStatus {
		t.Fatalf("expected ErrInvalidStatus, got %v", err)
	}
}
