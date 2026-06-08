package services

import (
	"TravelSphere/store"
	"TravelSphere/utils"
)

// ServiceContainer সব service একসাথে ধরে রাখে
// main.go তে একবার initialize হয়, সব controller এ share হয়
type ServiceContainer struct {
	CountryService    *CountryService
	AttractionService *AttractionService
	WishlistService   *WishlistService
	DashboardService  *DashboardService
	WeatherService    *WeatherService
}

// Global service container — সব controller এ ব্যবহার হবে
var Container *ServiceContainer

// InitServices সব service initialize করে
// main.go তে web.Run() এর আগে call করতে হবে
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