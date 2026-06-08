package store

import (
	"TravelSphere/models"
	"fmt"
	"sort"
	"sync"
	"time"
)

// WishlistStore is an in-memory storage for wishlist items
// It safely handles concurrent access using mutex
type WishlistStore struct {
	mu    sync.RWMutex
	items map[string]*models.WishlistItem
}

// NewWishlistStore creates and returns a new empty wishlist store
func NewWishlistStore() *WishlistStore {
	return &WishlistStore{
		items: make(map[string]*models.WishlistItem),
	}
}

// GetByUsername returns all wishlist items for a specific user
func (s *WishlistStore) GetByUsername(username string) []*models.WishlistItem {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*models.WishlistItem
	for _, item := range s.items {
		if item.Username == username {
			result = append(result, item)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	return result
}

// GetByID returns a single wishlist item by its ID
func (s *WishlistStore) GetByID(id string) (*models.WishlistItem, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[id]
	return item, ok
}

// Create adds a new wishlist item to the store
func (s *WishlistStore) Create(username, countryName, note, status string) *models.WishlistItem {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := fmt.Sprintf("%s-%d", username, time.Now().UnixNano())

	item := &models.WishlistItem{
		ID:          id,
		Username:    username,
		CountryName: countryName,
		Note:        note,
		Status:      models.WishlistStatus(status),
		CreatedAt:   time.Now(),
	}

	s.items[id] = item
	return item
}

// Update modifies an existing wishlist item
func (s *WishlistStore) Update(id, note, status string) (*models.WishlistItem, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[id]
	if !ok {
		return nil, false
	}

	item.Note = note
	item.Status = models.WishlistStatus(status)
	return item, true
}

// Delete removes a wishlist item from the store
func (s *WishlistStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return false
	}

	delete(s.items, id)
	return true
}

// IsOwner returns whether the given username owns the wishlist item
func (s *WishlistStore) IsOwner(id, username string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[id]
	return ok && item.Username == username
}

// CountByUsername returns summary statistics of a user's wishlist
// total = all items
// planned = planned visits
// visited = completed visits
func (s *WishlistStore) CountByUsername(username string) (total, planned, visited int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, item := range s.items {
		if item.Username != username {
			continue
		}
		total++
		switch item.Status {
		case models.StatusPlanned:
			planned++
		case models.StatusVisited:
			visited++
		}
	}
	return
}
