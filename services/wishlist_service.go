package services

import (
	"TravelSphere/models"
	"TravelSphere/store"
)

// WishlistService handles business logic for wishlist feature
// It communicates with WishlistStore (in-memory database)
type WishlistService struct {
	store *store.WishlistStore
}

// NewWishlistService creates a new WishlistService instance
func NewWishlistService(store *store.WishlistStore) *WishlistService {
	return &WishlistService{store: store}
}

// GetWishlist returns all wishlist items for a specific user
func (s *WishlistService) GetWishlist(username string) []*models.WishlistItem {
	return s.store.GetByUsername(username)
}

// Create adds a new wishlist item for a user
func (s *WishlistService) Create(username string, req *models.CreateWishlistRequest) (*models.WishlistItem, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return s.store.Create(username, req.CountryName, req.Note, req.Status)
}

// Update modifies an existing wishlist item
func (s *WishlistService) Update(username, id string, req *models.UpdateWishlistRequest) (*models.WishlistItem, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	item, err := s.store.GetByID(id)
	if err != nil {
		return nil, err
	}

	if item.Username != username {
		return nil, models.ErrUnauthorized
	}

	return s.store.Update(id, req.Note, req.Status)
}

// Delete removes a wishlist item
func (s *WishlistService) Delete(username, id string) error {
	item, err := s.store.GetByID(id)
	if err != nil {
		return err
	}

	if item.Username != username {
		return models.ErrUnauthorized
	}

	return s.store.Delete(id)
}
