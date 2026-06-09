package services

import (
	"TravelSphere/store"
	"TravelSphere/utils"
)

// ServiceContainer holds all services together.
// It is initialized once in main.go and shared across all controllers.
type ServiceContainer struct {
	CountryService    *CountryService
	AttractionService *AttractionService
	WishlistService   *WishlistService
	DashboardService  *DashboardService
	WeatherService    *WeatherService
}

var Container *ServiceContainer

// InitServices initializes all services.
// It should be called in main.go before web.Run().
func InitServices() {
	// Store
	wishlistStore := store.NewWishlistStore()

	// API Clients
	countriesClient := utils.NewCountriesClient()
	attractionClient := utils.NewOpenTripMapClient()
	weatherClient := utils.NewWeatherClient()

	// Services
	Container = &ServiceContainer{
		CountryService:    NewCountryService(countriesClient),
		AttractionService: NewAttractionService(attractionClient),
		WishlistService:   NewWishlistService(wishlistStore),
		DashboardService:  NewDashboardService(wishlistStore),
		WeatherService:    NewWeatherService(weatherClient),
	}
}
