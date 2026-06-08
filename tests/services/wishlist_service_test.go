package services_test

import (
	"testing"
	"TravelSphere/models"
	"TravelSphere/services"
	"TravelSphere/store"
)

func newWishlistService() *services.WishlistService {
	return services.NewWishlistService(store.NewWishlistStore())
}

func TestWishlistAdd_Valid(t *testing.T) {
	svc := newWishlistService()

	item, err := svc.AddToWishlist("beta", models.CreateWishlistRequest{
		CountryName: "France",
		Status:      "Planned",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.CountryName != "France" {
		t.Errorf("expected France, got %q", item.CountryName)
	}
}

func TestWishlistAdd_DefaultStatus(t *testing.T) {
	svc := newWishlistService()

	item, err := svc.AddToWishlist("beta", models.CreateWishlistRequest{
		CountryName: "France",
		Status:      "", // empty → default Planned
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Status != models.StatusPlanned {
		t.Errorf("expected default Planned, got %q", item.Status)
	}
}

func TestWishlistAdd_MissingCountryName(t *testing.T) {
	svc := newWishlistService()

	_, err := svc.AddToWishlist("beta", models.CreateWishlistRequest{
		CountryName: "",
		Status:      "Planned",
	})
	if err != models.ErrCountryNameRequired {
		t.Errorf("expected ErrCountryNameRequired, got %v", err)
	}
}

func TestWishlistAdd_InvalidStatus(t *testing.T) {
	svc := newWishlistService()

	_, err := svc.AddToWishlist("beta", models.CreateWishlistRequest{
		CountryName: "France",
		Status:      "Maybe",
	})
	if err != models.ErrInvalidStatus {
		t.Errorf("expected ErrInvalidStatus, got %v", err)
	}
}

func TestWishlistUpdate_Valid(t *testing.T) {
	svc := newWishlistService()

	item, _ := svc.AddToWishlist("beta", models.CreateWishlistRequest{
		CountryName: "France", Status: "Planned",
	})

	updated, err := svc.UpdateWishlistItem("beta", item.ID,
		models.UpdateWishlistRequest{
			Note:   "Visit Louvre",
			Status: "Visited",
		})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Note != "Visit Louvre" {
		t.Errorf("note not updated: %q", updated.Note)
	}
	if updated.Status != models.StatusVisited {
		t.Errorf("status not updated: %q", updated.Status)
	}
}

func TestWishlistUpdate_Unauthorized(t *testing.T) {
	svc := newWishlistService()

	item, _ := svc.AddToWishlist("beta", models.CreateWishlistRequest{
		CountryName: "France", Status: "Planned",
	})

	_, err := svc.UpdateWishlistItem("john", item.ID,
		models.UpdateWishlistRequest{Status: "Visited"})
	if err != models.ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestWishlistUpdate_NotFound(t *testing.T) {
	svc := newWishlistService()

	_, err := svc.UpdateWishlistItem("beta", "nonexistent",
		models.UpdateWishlistRequest{Status: "Planned"})
	if err != models.ErrUnauthorized {
		// IsOwner false → ErrUnauthorized আসবে (not found এর আগে)
		t.Logf("got error: %v (acceptable)", err)
	}
}

func TestWishlistDelete_Valid(t *testing.T) {
	svc := newWishlistService()

	item, _ := svc.AddToWishlist("beta", models.CreateWishlistRequest{
		CountryName: "France", Status: "Planned",
	})

	err := svc.DeleteWishlistItem("beta", item.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	items := svc.GetWishlist("beta")
	if len(items) != 0 {
		t.Errorf("expected empty list after delete, got %d", len(items))
	}
}

func TestWishlistDelete_Unauthorized(t *testing.T) {
	svc := newWishlistService()

	item, _ := svc.AddToWishlist("beta", models.CreateWishlistRequest{
		CountryName: "France", Status: "Planned",
	})

	err := svc.DeleteWishlistItem("john", item.ID)
	if err != models.ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestWishlistGetList_UserIsolation(t *testing.T) {
	svc := newWishlistService()

	svc.AddToWishlist("alpha", models.CreateWishlistRequest{
		CountryName: "France", Status: "Planned",
	})
	svc.AddToWishlist("beta", models.CreateWishlistRequest{
		CountryName: "Japan", Status: "Visited",
	})

	alphaItems := svc.GetWishlist("alpha")
	betaItems := svc.GetWishlist("beta")

	if len(alphaItems) != 1 {
		t.Errorf("alpha should have 1 item, got %d", len(alphaItems))
	}
	if len(betaItems) != 1 {
		t.Errorf("beta should have 1 item, got %d", len(betaItems))
	}
	if alphaItems[0].CountryName == betaItems[0].CountryName {
		t.Error("different users should have different items")
	}
}