package store_test

import (
	"testing"
	"time"
	"TravelSphere/models"
	"TravelSphere/store"
)

// ── Setup helper ──

func newStore() *store.WishlistStore {
	return store.NewWishlistStore()
}

// ── Create Tests ──

func TestCreate_ReturnsItemWithID(t *testing.T) {
	s := newStore()
	item := s.Create("beta", "Bangladesh", "Visit Dhaka", "Planned")

	if item.ID == "" {
		t.Error("ID should not be empty")
	}
	if item.Username != "beta" {
		t.Errorf("expected username 'beta', got %q", item.Username)
	}
	if item.CountryName != "Bangladesh" {
		t.Errorf("expected 'Bangladesh', got %q", item.CountryName)
	}
	if item.Note != "Visit Dhaka" {
		t.Errorf("expected note 'Visit Dhaka', got %q", item.Note)
	}
	if item.Status != models.StatusPlanned {
		t.Errorf("expected Planned, got %q", item.Status)
	}
}

func TestCreate_SetsCreatedAt(t *testing.T) {
	s := newStore()
	before := time.Now()
	item := s.Create("beta", "France", "", "Planned")
	after := time.Now()

	if item.CreatedAt.Before(before) || item.CreatedAt.After(after) {
		t.Error("CreatedAt should be set to current time")
	}
}

func TestCreate_UniqueIDs(t *testing.T) {
	s := newStore()
	item1 := s.Create("beta", "France", "", "Planned")
	// Small delay to ensure unique nanosecond timestamp
	time.Sleep(1 * time.Millisecond)
	item2 := s.Create("beta", "Japan", "", "Planned")

	if item1.ID == item2.ID {
		t.Error("each item should have a unique ID")
	}
}

// ── GetByUsername Tests ──

func TestGetByUsername_ReturnsOnlyUserItems(t *testing.T) {
	s := newStore()

	s.Create("beta", "France", "", "Planned")
	s.Create("beta", "Japan", "", "Visited")
	s.Create("john", "USA", "", "Planned")

	betaItems := s.GetByUsername("beta")
	if len(betaItems) != 2 {
		t.Errorf("expected 2 items for beta, got %d", len(betaItems))
	}

	johnItems := s.GetByUsername("john")
	if len(johnItems) != 1 {
		t.Errorf("expected 1 item for john, got %d", len(johnItems))
	}
}

func TestGetByUsername_EmptyForNewUser(t *testing.T) {
	s := newStore()
	items := s.GetByUsername("nobody")
	if len(items) != 0 {
		t.Errorf("new user should have 0 items, got %d", len(items))
	}
}

func TestGetByUsername_NewestFirst(t *testing.T) {
	s := newStore()
	s.Create("beta", "France", "", "Planned")
	time.Sleep(2 * time.Millisecond)
	s.Create("beta", "Japan", "", "Planned")

	items := s.GetByUsername("beta")
	if len(items) < 2 {
		t.Fatal("expected 2 items")
	}
	// Japan (newer) should come first
	if items[0].CountryName != "Japan" {
		t.Errorf("expected newest first (Japan), got %q", items[0].CountryName)
	}
}

// ── GetByID Tests ──

func TestGetByID_Found(t *testing.T) {
	s := newStore()
	created := s.Create("beta", "France", "", "Planned")

	item, exists := s.GetByID(created.ID)
	if !exists {
		t.Error("expected item to exist")
	}
	if item.ID != created.ID {
		t.Errorf("expected ID %q, got %q", created.ID, item.ID)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	s := newStore()
	_, exists := s.GetByID("nonexistent-id")
	if exists {
		t.Error("expected item not to exist")
	}
}

// ── Update Tests ──

func TestUpdate_UpdatesNoteAndStatus(t *testing.T) {
	s := newStore()
	created := s.Create("beta", "France", "", "Planned")

	updated, ok := s.Update(created.ID, "Visit Eiffel Tower", "Visited")
	if !ok {
		t.Fatal("expected update to succeed")
	}
	if updated.Note != "Visit Eiffel Tower" {
		t.Errorf("expected updated note, got %q", updated.Note)
	}
	if updated.Status != models.StatusVisited {
		t.Errorf("expected Visited, got %q", updated.Status)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	s := newStore()
	_, ok := s.Update("nonexistent", "note", "Planned")
	if ok {
		t.Error("expected update to fail for nonexistent ID")
	}
}

func TestUpdate_PreservesOtherFields(t *testing.T) {
	s := newStore()
	created := s.Create("beta", "France", "", "Planned")

	updated, _ := s.Update(created.ID, "new note", "Visited")
	// Country name আর username same থাকবে
	if updated.CountryName != "France" {
		t.Errorf("country name should not change, got %q", updated.CountryName)
	}
	if updated.Username != "beta" {
		t.Errorf("username should not change, got %q", updated.Username)
	}
}

// ── Delete Tests ──

func TestDelete_Success(t *testing.T) {
	s := newStore()
	created := s.Create("beta", "France", "", "Planned")

	ok := s.Delete(created.ID)
	if !ok {
		t.Error("expected delete to succeed")
	}

	// Item আর নেই
	_, exists := s.GetByID(created.ID)
	if exists {
		t.Error("item should not exist after delete")
	}
}

func TestDelete_NotFound(t *testing.T) {
	s := newStore()
	ok := s.Delete("nonexistent")
	if ok {
		t.Error("expected delete to fail for nonexistent ID")
	}
}

// ── IsOwner Tests ──

func TestIsOwner_True(t *testing.T) {
	s := newStore()
	created := s.Create("beta", "France", "", "Planned")

	if !s.IsOwner(created.ID, "beta") {
		t.Error("beta should be the owner")
	}
}

func TestIsOwner_False(t *testing.T) {
	s := newStore()
	created := s.Create("beta", "France", "", "Planned")

	if s.IsOwner(created.ID, "john") {
		t.Error("john should not be the owner")
	}
}

func TestIsOwner_NonexistentID(t *testing.T) {
	s := newStore()
	if s.IsOwner("nonexistent", "beta") {
		t.Error("nonexistent item should return false")
	}
}

// ── CountByUsername Tests ──

func TestCountByUsername(t *testing.T) {
	s := newStore()
	s.Create("beta", "France", "", "Planned")
	s.Create("beta", "Japan", "", "Visited")
	s.Create("beta", "USA", "", "Planned")
	// john এর item — beta count এ আসবে না
	s.Create("john", "Germany", "", "Planned")

	total, planned, visited := s.CountByUsername("beta")

	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if planned != 2 {
		t.Errorf("expected planned 2, got %d", planned)
	}
	if visited != 1 {
		t.Errorf("expected visited 1, got %d", visited)
	}
}

func TestCountByUsername_Empty(t *testing.T) {
	s := newStore()
	total, planned, visited := s.CountByUsername("nobody")

	if total != 0 || planned != 0 || visited != 0 {
		t.Errorf("new user should have 0 counts, got %d/%d/%d",
			total, planned, visited)
	}
}

// ── Concurrent Access Test ──

func TestConcurrentAccess(t *testing.T) {
	s := newStore()

	// 10টা goroutine একসাথে write করবে
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(n int) {
			s.Create("user", "Country", "", "Planned")
			done <- true
		}(i)
	}

	// সব goroutine শেষ হওয়া পর্যন্ত wait করো
	for i := 0; i < 10; i++ {
		<-done
	}

	items := s.GetByUsername("user")
	if len(items) != 10 {
		t.Errorf("expected 10 items after concurrent writes, got %d",
			len(items))
	}
}