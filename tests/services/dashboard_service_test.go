package services_test

import (
	"testing"
	"TravelSphere/models"
	"TravelSphere/services"
	"TravelSphere/store"
)

func newDashboardSetup() (*services.WishlistService, *services.DashboardService) {
	s := store.NewWishlistStore()
	return services.NewWishlistService(s), services.NewDashboardService(s)
}

func TestDashboardSummary_AllCounts(t *testing.T) {
	wishSvc, dashSvc := newDashboardSetup()

	wishSvc.AddToWishlist("beta", models.CreateWishlistRequest{
		CountryName: "France", Status: "Planned",
	})
	wishSvc.AddToWishlist("beta", models.CreateWishlistRequest{
		CountryName: "Japan", Status: "Visited",
	})
	wishSvc.AddToWishlist("beta", models.CreateWishlistRequest{
		CountryName: "USA", Status: "Planned",
	})

	summary := dashSvc.GetSummary("beta")

	if summary.Total != 3 {
		t.Errorf("expected total 3, got %d", summary.Total)
	}
	if summary.Planned != 2 {
		t.Errorf("expected planned 2, got %d", summary.Planned)
	}
	if summary.Visited != 1 {
		t.Errorf("expected visited 1, got %d", summary.Visited)
	}
}

func TestDashboardSummary_Empty(t *testing.T) {
	_, dashSvc := newDashboardSetup()

	summary := dashSvc.GetSummary("newuser")
	if summary.Total != 0 || summary.Planned != 0 || summary.Visited != 0 {
		t.Error("new user should have all zeros")
	}
}

func TestDashboardSummary_UserIsolation(t *testing.T) {
	wishSvc, dashSvc := newDashboardSetup()

	// beta এর items
	wishSvc.AddToWishlist("beta", models.CreateWishlistRequest{
		CountryName: "France", Status: "Planned",
	})

	// john এর items
	wishSvc.AddToWishlist("john", models.CreateWishlistRequest{
		CountryName: "USA", Status: "Visited",
	})
	wishSvc.AddToWishlist("john", models.CreateWishlistRequest{
		CountryName: "Germany", Status: "Visited",
	})

	betaSummary := dashSvc.GetSummary("beta")
	johnSummary := dashSvc.GetSummary("john")

	if betaSummary.Total != 1 {
		t.Errorf("beta total should be 1, got %d", betaSummary.Total)
	}
	if johnSummary.Total != 2 {
		t.Errorf("john total should be 2, got %d", johnSummary.Total)
	}
}

func TestGetSavedDestinations(t *testing.T) {
	wishSvc, dashSvc := newDashboardSetup()

	wishSvc.AddToWishlist("beta", models.CreateWishlistRequest{
		CountryName: "France", Status: "Planned",
	})
	wishSvc.AddToWishlist("beta", models.CreateWishlistRequest{
		CountryName: "Japan", Status: "Visited",
	})

	destinations := dashSvc.GetSavedDestinations("beta")
	if len(destinations) != 2 {
		t.Errorf("expected 2 destinations, got %d", len(destinations))
	}
}