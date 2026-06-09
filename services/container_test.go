package services_test

import (
	"testing"

	"TravelSphere/services"
)

func TestInitServices_PopulatesContainer(t *testing.T) {
	services.InitServices()
	if services.Container == nil {
		t.Fatal("expected global service container to be initialized")
	}
	if services.Container.CountryService == nil || services.Container.AttractionService == nil || services.Container.WishlistService == nil || services.Container.DashboardService == nil || services.Container.WeatherService == nil {
		t.Fatal("expected all services to be initialized")
	}
}
