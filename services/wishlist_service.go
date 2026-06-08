package services

import (
	"TravelSphere/models"
	"TravelSphere/store"
)

// WishlistServiceInterface mock করার জন্য interface
type WishlistServiceInterface interface {
	GetWishlist(username string) []*models.WishlistItem
	AddToWishlist(username string, req models.CreateWishlistRequest) (*models.WishlistItem, error)
	UpdateWishlistItem(username, id string, req models.UpdateWishlistRequest) (*models.WishlistItem, error)
	DeleteWishlistItem(username, id string) error
}

// WishlistService wishlist related সব business logic
type WishlistService struct {
	store *store.WishlistStore
}

// NewWishlistService নতুন WishlistService তৈরি করে
func NewWishlistService(s *store.WishlistStore) *WishlistService {
	return &WishlistService{store: s}
}

// GetWishlist user এর সব wishlist items return করে
func (s *WishlistService) GetWishlist(username string) []*models.WishlistItem {
	return s.store.GetByUsername(username)
}

// AddToWishlist নতুন item wishlist এ যোগ করে
func (s *WishlistService) AddToWishlist(username string, req models.CreateWishlistRequest) (*models.WishlistItem, error) {
	// Validate করো
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Default status set করো
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

// UpdateWishlistItem existing item update করে
func (s *WishlistService) UpdateWishlistItem(
	username, id string,
	req models.UpdateWishlistRequest,
) (*models.WishlistItem, error) {
	// Validate করো
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Ownership check — অন্য user এর item update করা যাবে না
	if !s.store.IsOwner(id, username) {
		return nil, models.ErrUnauthorized
	}

	item, exists := s.store.Update(id, req.Note, req.Status)
	if !exists {
		return nil, models.ErrNotFound
	}
	return item, nil
}

// DeleteWishlistItem item delete করে
func (s *WishlistService) DeleteWishlistItem(username, id string) error {
	// Ownership check
	if !s.store.IsOwner(id, username) {
		return models.ErrUnauthorized
	}

	if !s.store.Delete(id) {
		return models.ErrNotFound
	}
	return nil
}