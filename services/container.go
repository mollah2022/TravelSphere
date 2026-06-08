package services

import (
	"TravelSphere/store"
	"TravelSphere/utils"
)

// ServiceContainer holds all services of the application
// This helps to access all services from one place (central container)
type ServiceContainer struct {
	CountryService    *CountryService
	AttractionService *AttractionService
	WishlistService   *WishlistService
	DashboardService  *DashboardService
	WeatherService    *WeatherService
}

var Container *ServiceContainer

// InitServices initializes all services and injects required dependencies
func InitServices() {
	wishlistStore := store.NewWishlistStore()

	countriesClient := utils.NewCountriesClient()
	attractionClient := utils.NewOpenTripMapClient()
	weatherClient := utils.NewWeatherClient()

	Container = &ServiceContainer{
		CountryService:    NewCountryService(countriesClient),
		AttractionService: NewAttractionService(attractionClient),
		WishlistService:   NewWishlistService(wishlistStore),
		DashboardService:  NewDashboardService(wishlistStore),
		WeatherService:    NewWeatherService(weatherClient),
	}
}
