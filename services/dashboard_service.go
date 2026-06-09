package services

import (
	"TravelSphere/models"
	"TravelSphere/store"
)

// DashboardServiceInterface is used to mock the dashboard service in tests.
type DashboardServiceInterface interface {
	GetSummary(username string) models.DashboardSummary
	GetSavedDestinations(username string) []*models.WishlistItem
}

// DashboardService dashboard related business logic
type DashboardService struct {
	store *store.WishlistStore
}

// NewDashboardService creates a new instance of DashboardService
func NewDashboardService(s *store.WishlistStore) *DashboardService {
	return &DashboardService{store: s}
}

// GetSummary returns dashboard statistics for a user.
func (s *DashboardService) GetSummary(username string) models.DashboardSummary {
	total, planned, visited := s.store.CountByUsername(username)
	return models.DashboardSummary{
		Total:   total,
		Planned: planned,
		Visited: visited,
	}
}

// GetSavedDestinations returns all saved destinations for a user.
func (s *DashboardService) GetSavedDestinations(username string) []*models.WishlistItem {
	return s.store.GetByUsername(username)
}
