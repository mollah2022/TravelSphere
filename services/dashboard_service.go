package services

import (
	"TravelSphere/models"
	"TravelSphere/store"
)

// DashboardServiceInterface mock করার জন্য interface
type DashboardServiceInterface interface {
	GetSummary(username string) models.DashboardSummary
	GetSavedDestinations(username string) []*models.WishlistItem
}

// DashboardService dashboard related business logic
type DashboardService struct {
	store *store.WishlistStore
}

// NewDashboardService নতুন DashboardService তৈরি করে
func NewDashboardService(s *store.WishlistStore) *DashboardService {
	return &DashboardService{store: s}
}

// GetSummary user এর dashboard stats return করে
func (s *DashboardService) GetSummary(username string) models.DashboardSummary {
	total, planned, visited := s.store.CountByUsername(username)
	return models.DashboardSummary{
		Total:   total,
		Planned: planned,
		Visited: visited,
	}
}

// GetSavedDestinations user এর সব saved destinations return করে
func (s *DashboardService) GetSavedDestinations(username string) []*models.WishlistItem {
	return s.store.GetByUsername(username)
}