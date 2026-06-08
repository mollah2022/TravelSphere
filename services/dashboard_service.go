package services

import (
	"TravelSphere/models"
	"TravelSphere/store"
)

// DashboardService handles dashboard-related business logic
// It works with wishlist data to generate summary and user overview
type DashboardService struct {
	store *store.WishlistStore
}

// NewDashboardService creates a new DashboardService instance
func NewDashboardService(store *store.WishlistStore) *DashboardService {
	return &DashboardService{store: store}
}

// GetSummary returns dashboard summary for a specific user
// It shows total, planned, and visited destinations
func (s *DashboardService) GetSummary(username string) models.DashboardSummary {
	total, planned, visited := s.store.CountByUsername(username)
	return models.DashboardSummary{
		Total:   total,
		Planned: planned,
		Visited: visited,
	}
}

// GetSavedDestinations returns all saved wishlist items for a user
func (s *DashboardService) GetSavedDestinations(username string) []*models.WishlistItem {
	return s.store.GetByUsername(username)
}
