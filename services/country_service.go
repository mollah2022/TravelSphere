package services

import (
	"TravelSphere/models"
	"TravelSphere/utils"
	"fmt"
	"strings"
	"sync"
)

// CountryServiceInterface is used to mock the country service in tests.
type CountryServiceInterface interface {
	GetAllCountries() ([]models.CountryDTO, error)
	SearchCountries(search, region string) ([]models.CountryDTO, error)
	GetCountryBySlug(slug string) (*models.CountryDTO, error)
	GetFeaturedCountries() ([]models.CountryDTO, error)
}

// CountryService country related business logic
type CountryService struct {
	client utils.CountriesClientInterface

	// In-memory cache
	cacheMu    sync.RWMutex
	cachedData []models.CountryDTO
	cacheReady bool
}

// NewCountryService creates a new instance of CountryService
func NewCountryService(client utils.CountriesClientInterface) *CountryService {
	return &CountryService{
		client: client,
	}
}

// getAllCached fetches countries from cache; if not available, it calls the API.
func (s *CountryService) getAllCached() ([]models.CountryDTO, error) {

	s.cacheMu.RLock()
	if s.cacheReady {
		data := s.cachedData
		s.cacheMu.RUnlock()
		return data, nil
	}
	s.cacheMu.RUnlock()

	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()

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

// GetAllCountries returns all countries.
func (s *CountryService) GetAllCountries() ([]models.CountryDTO, error) {
	return s.getAllCached()
}

// SearchCountries search and region parameters to filter countries.
func (s *CountryService) SearchCountries(search, region string) ([]models.CountryDTO, error) {
	all, err := s.getAllCached()
	if err != nil {
		return nil, err
	}
	return utils.FilterCountries(all, search, region), nil
}

// GetCountryBySlug finds a country using its slug (e.g., "bangladesh") and returns its DTO.
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

		if strings.ToLower(c.CCA2) == slug {
			return &c, nil
		}

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

// GetFeaturedCountries returns featured countries for the home page.
func (s *CountryService) GetFeaturedCountries() ([]models.CountryDTO, error) {

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

// GetSearchSuggestions returns search autocomplete suggestions for the home page.
func (s *CountryService) GetSearchSuggestions(query string) ([]models.CountryDTO, error) {
	if query == "" {
		return []models.CountryDTO{}, nil
	}
	results, err := s.SearchCountries(query, "")
	if err != nil {
		return nil, err
	}

	if len(results) > 8 {
		return results[:8], nil
	}
	return results, nil
}
