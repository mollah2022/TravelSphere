package services

import (
	"TravelSphere/models"
	"TravelSphere/utils"
	"fmt"
	"strings"
	"sync"
)

// CountryServiceInterface mock করার জন্য interface
type CountryServiceInterface interface {
	GetAllCountries() ([]models.CountryDTO, error)
	SearchCountries(search, region string) ([]models.CountryDTO, error)
	GetCountryBySlug(slug string) (*models.CountryDTO, error)
	GetFeaturedCountries() ([]models.CountryDTO, error)
}

// CountryService country related সব business logic
type CountryService struct {
	client utils.CountriesClientInterface

	// In-memory cache
	cacheMu    sync.RWMutex
	cachedData []models.CountryDTO
	cacheReady bool
}

// NewCountryService নতুন CountryService তৈরি করে
func NewCountryService(client utils.CountriesClientInterface) *CountryService {
	return &CountryService{
		client: client,
	}
}

// getAllCached cache থেকে দেশ আনে, না থাকলে API call করে
func (s *CountryService) getAllCached() ([]models.CountryDTO, error) {
	// Read lock দিয়ে check করো cache আছে কিনা
	s.cacheMu.RLock()
	if s.cacheReady {
		data := s.cachedData
		s.cacheMu.RUnlock()
		return data, nil
	}
	s.cacheMu.RUnlock()

	// Cache নেই — API থেকে আনো
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()

	// Double check (অন্য goroutine এর মধ্যে cache হয়ে গেছে কিনা)
	if s.cacheReady {
		return s.cachedData, nil
	}

	raw, err := s.client.FetchAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch countries: %w", err)
	}

	// Transform raw → DTO
	dtos := make([]models.CountryDTO, 0, len(raw))
	for _, c := range raw {
		dtos = append(dtos, utils.TransformCountryToDTO(c))
	}

	s.cachedData = dtos
	s.cacheReady = true
	return dtos, nil
}

// GetAllCountries সব দেশ return করে
func (s *CountryService) GetAllCountries() ([]models.CountryDTO, error) {
	return s.getAllCached()
}

// SearchCountries search ও region দিয়ে filter করে
func (s *CountryService) SearchCountries(search, region string) ([]models.CountryDTO, error) {
	all, err := s.getAllCached()
	if err != nil {
		return nil, err
	}
	return utils.FilterCountries(all, search, region), nil
}

// GetCountryBySlug slug দিয়ে একটা দেশ খোঁজে
// যেমন: "bangladesh" → Bangladesh এর DTO
func (s *CountryService) GetCountryBySlug(slug string) (*models.CountryDTO, error) {
	all, err := s.getAllCached()
	if err != nil {
		return nil, err
	}

	slug = strings.ToLower(strings.TrimSpace(slug))
	for _, c := range all {
		if c.Slug == slug {
			return &c, nil
		}
		// CCA2 code দিয়েও খোঁজা যাবে যেমন "us", "bd"
		if strings.ToLower(c.CCA2) == slug {
			return &c, nil
		}
		// CCA3 code দিয়েও যেমন "usa", "bgd"
		if strings.ToLower(c.CCA3) == slug {
			return &c, nil
		}
	}

	// Fallback: allow typed country names or hyphenated names
	if normalizedName := strings.TrimSpace(utils.SlugToName(slug)); normalizedName != "" {
		for _, c := range all {
			if strings.EqualFold(c.Name, normalizedName) || strings.EqualFold(c.OfficialName, normalizedName) {
				return &c, nil
			}
		}
	}

	return nil, models.ErrNotFound
}

// GetFeaturedCountries home page এর জন্য featured দেশ return করে
func (s *CountryService) GetFeaturedCountries() ([]models.CountryDTO, error) {
	// স্ক্রিনশট অনুযায়ী: USA, France, Japan, Australia, Brazil, Bangladesh
	featuredNames := []string{
		"united-states",
		"france",
		"japan",
		"australia",
		"brazil",
		"bangladesh",
	}

	all, err := s.getAllCached()
	if err != nil {
		return nil, err
	}

	// slug দিয়ে featured দেশ গুলো খুঁজে বের করো
	slugMap := make(map[string]models.CountryDTO)
	for _, c := range all {
		slugMap[c.Slug] = c
	}

	result := make([]models.CountryDTO, 0, len(featuredNames))
	for _, slug := range featuredNames {
		if c, ok := slugMap[slug]; ok {
			result = append(result, c)
		}
	}
	return result, nil
}

// GetSearchSuggestions home page search autocomplete এর জন্য
func (s *CountryService) GetSearchSuggestions(query string) ([]models.CountryDTO, error) {
	if query == "" {
		return []models.CountryDTO{}, nil
	}
	results, err := s.SearchCountries(query, "")
	if err != nil {
		return nil, err
	}
	// সর্বোচ্চ 8টা suggestion দেখাও
	if len(results) > 8 {
		return results[:8], nil
	}
	return results, nil
}
