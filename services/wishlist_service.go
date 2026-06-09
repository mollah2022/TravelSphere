package services

import (
	"TravelSphere/models"
	"TravelSphere/store"
)

// WishlistServiceInterface is used to mock the wishlist service in tests.
type WishlistServiceInterface interface {
	GetWishlist(username string) []*models.WishlistItem
	AddToWishlist(username string, req models.CreateWishlistRequest) (*models.WishlistItem, error)
	UpdateWishlistItem(username, id string, req models.UpdateWishlistRequest) (*models.WishlistItem, error)
	DeleteWishlistItem(username, id string) error
}

// WishlistService manages wishlist operations and business logic.
type WishlistService struct {
	store *store.WishlistStore
}

// NewWishlistService creates a new instance of WishlistService.
func NewWishlistService(s *store.WishlistStore) *WishlistService {
	return &WishlistService{store: s}
}

// GetWishlist returns all wishlist items of a user.
func (s *WishlistService) GetWishlist(username string) []*models.WishlistItem {
	return s.store.GetByUsername(username)
}

// AddToWishlist adds a new item to the wishlist.
func (s *WishlistService) AddToWishlist(username string, req models.CreateWishlistRequest) (*models.WishlistItem, error) {

	if err := req.Validate(); err != nil {
		return nil, err
	}

	if req.Status == "" {
		req.Status = string(models.StatusPlanned)
	}

	item := s.store.Create(
		username,
		req.CountryName,
		req.Note,
		req.Status,
	)
	return item, nil
}

// UpdateWishlistItem updates an existing wishlist item.
func (s *WishlistService) UpdateWishlistItem(
	username, id string,
	req models.UpdateWishlistRequest,
) (*models.WishlistItem, error) {

	if err := req.Validate(); err != nil {
		return nil, err
	}

	if !s.store.IsOwner(id, username) {
		return nil, models.ErrUnauthorized
	}

	item, exists := s.store.Update(id, req.Note, req.Status)
	if !exists {
		return nil, models.ErrNotFound
	}
	return item, nil
}

// DeleteWishlistItem deletes a wishlist item.
func (s *WishlistService) DeleteWishlistItem(username, id string) error {

	if !s.store.IsOwner(id, username) {
		return models.ErrUnauthorized
	}

	if !s.store.Delete(id) {
		return models.ErrNotFound
	}
	return nil
}
