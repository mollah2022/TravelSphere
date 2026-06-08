package services

import (
	"TravelSphere/models"
	"TravelSphere/utils"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

// CountryService handles all country-related business logic
// It also includes caching to improve performance
type CountryService struct {
	client    *utils.CountriesClient
	cache     []models.CountryDTO
	cacheTime time.Time
	mu        sync.RWMutex
}

// NewCountryService creates a new CountryService instance
func NewCountryService(client *utils.CountriesClient) *CountryService {
	return &CountryService{client: client}
}

// GetAllCountries returns all countries
// It uses cache to avoid calling API again and again
func (s *CountryService) GetAllCountries() ([]models.CountryDTO, error) {
	s.mu.RLock()
	if s.cache != nil && time.Since(s.cacheTime) < 10*time.Minute {
		defer s.mu.RUnlock()
		return s.cache, nil
	}
	s.mu.RUnlock()

	raw, err := s.client.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch countries: %w", err)
	}

	var dtos []models.CountryDTO
	for _, r := range raw {
		dtos = append(dtos, mapToDTO(r))
	}

	sort.Slice(dtos, func(i, j int) bool {
		return dtos[i].Name < dtos[j].Name
	})

	s.mu.Lock()
	s.cache = dtos
	s.cacheTime = time.Now()
	s.mu.Unlock()

	return dtos, nil
}

// GetCountryBySlug returns a single country by its slug (URL-friendly name)
func (s *CountryService) GetCountryBySlug(slug string) (*models.CountryDTO, error) {
	countries, err := s.GetAllCountries()
	if err != nil {
		return nil, err
	}

	for _, c := range countries {
		if c.Slug == slug {
			return &c, nil
		}
	}
	return nil, models.ErrNotFound
}

// GetFeaturedCountries returns a predefined list of popular countries
func (s *CountryService) GetFeaturedCountries() ([]models.CountryDTO, error) {
	featured := []string{
		"france", "japan", "brazil", "egypt",
		"australia", "canada", "italy", "india",
	}

	countries, err := s.GetAllCountries()
	if err != nil {
		return nil, err
	}

	slugMap := make(map[string]models.CountryDTO)
	for _, c := range countries {
		slugMap[c.Slug] = c
	}

	var result []models.CountryDTO
	for _, slug := range featured {
		if c, ok := slugMap[slug]; ok {
			result = append(result, c)
		}
	}
	return result, nil
}

// SearchCountries searches countries by name or capital
func (s *CountryService) SearchCountries(query string) []models.CountryDTO {
	countries, err := s.GetAllCountries()
	if err != nil {
		return nil
	}

	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return nil
	}

	var result []models.CountryDTO
	for _, c := range countries {
		if strings.Contains(strings.ToLower(c.Name), q) ||
			strings.Contains(strings.ToLower(c.Capital), q) {
			result = append(result, c)
			if len(result) >= 8 {
				break
			}
		}
	}
	return result
}

// mapToDTO converts raw API response to clean DTO format
func mapToDTO(r models.CountryResponse) models.CountryDTO {
	capital := ""
	if len(r.Capital) > 0 {
		capital = r.Capital[0]
	}

	var currencies []string
	for code, cur := range r.Currencies {
		currencies = append(currencies, fmt.Sprintf("%s (%s)", code, cur.Name))
	}

	var languages []string
	for _, lang := range r.Languages {
		languages = append(languages, lang)
	}
	sort.Strings(languages)

	lat, lon := 0.0, 0.0
	if len(r.LatLng) >= 2 {
		lat = r.LatLng[0]
		lon = r.LatLng[1]
	}

	slug := strings.ToLower(strings.ReplaceAll(r.Name.Common, " ", "-"))
	slug = strings.ReplaceAll(slug, "'", "")
	slug = strings.ReplaceAll(slug, ".", "")
	slug = strings.ReplaceAll(slug, ",", "")

	return models.CountryDTO{
		Slug:         slug,
		Name:         r.Name.Common,
		OfficialName: r.Name.Official,
		Capital:      capital,
		Region:       r.Region,
		Subregion:    r.Subregion,
		Population:   r.Population,
		FlagURL:      r.Flags.PNG,
		FlagAlt:      r.Flags.Alt,
		Currencies:   strings.Join(currencies, ", "),
		Languages:    strings.Join(languages, ", "),
		CCA2:         r.CCA2,
		CCA3:         r.CCA3,
		Latitude:     lat,
		Longitude:    lon,
	}
}
