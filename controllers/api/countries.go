package apicontrollers

import (
	"TravelSphere/services"
	"TravelSphere/utils"

	"github.com/beego/beego/v2/server/web"
)

// CountriesAPIController provides a JSON API for countries.
// All responses are JSON; no HTML is returned.
type CountriesAPIController struct {
	web.Controller
}

// svc accesses the service container.
func svc() *services.ServiceContainer {
	return services.Container
}

// List handles GET /api/countries.
// Supports query parameters: search, region, limit, offset.
// Defaults: limit=25, offset=0
// Max limit: 100
// Used for AJAX-based country search and filtering.
func (c *CountriesAPIController) List() {
	search := c.GetString("search")
	region := c.GetString("region")
	q := c.GetString("q") // Alternative search parameter
	limit, _ := c.GetInt("limit")
	offset, _ := c.GetInt("offset")

	// Use 'q' if 'search' is not provided (REST Countries v5 API uses 'q')
	if search == "" && q != "" {
		search = q
	}

	if !utils.IsValidSearch(search) {
		utils.SendError(&c.Controller, "Search query too long", 400)
		return
	}

	if !utils.IsValidRegion(region) {
		utils.SendError(&c.Controller, "Invalid region", 400)
		return
	}

	// Validate pagination parameters
	if limit < 1 || limit > 100 {
		limit = 25
	}
	if offset < 0 {
		offset = 0
	}

	countries, err := svc().CountryService.SearchCountriesWithPagination(search, region, limit, offset)
	if err != nil {
		utils.SendError(&c.Controller, "Failed to fetch countries", 500)
		return
	}

	utils.SendSuccess(&c.Controller, countries, "", 200)
}

// Detail handles GET /api/countries/:slug.
// Returns a single country as JSON.
func (c *CountriesAPIController) Detail() {
	slug := c.Ctx.Input.Param(":slug")

	if !utils.IsValidSlug(slug) {
		utils.SendError(&c.Controller, "Invalid country slug", 400)
		return
	}

	country, err := svc().CountryService.GetCountryBySlug(slug)
	if err != nil {
		utils.SendError(&c.Controller, "Country not found", 404)
		return
	}

	utils.SendSuccess(&c.Controller, country, "", 200)
}

// Attractions handles GET /api/attractions.
// Supports query parameters: lat, lon.
// Used for AJAX loading of attractions on the destination page.
func (c *CountriesAPIController) Attractions() {
	lat, err1 := c.GetFloat("lat")
	lon, err2 := c.GetFloat("lon")

	if err1 != nil || err2 != nil {
		utils.SendError(&c.Controller, "Invalid coordinates", 400)
		return
	}

	attractions, err := svc().AttractionService.GetAttractionsByCountry(lat, lon)
	if err != nil {
		utils.SendError(&c.Controller, "Failed to fetch attractions", 500)
		return
	}

	utils.SendSuccess(&c.Controller, attractions, "", 200)
}

// Suggestions handles GET /api/suggestions.
// Supports query parameter: q.
// Used for home page search autocomplete.
func (c *CountriesAPIController) Suggestions() {
	query := c.GetString("q")

	if query == "" {
		utils.SendSuccess(&c.Controller, []interface{}{}, "", 200)
		return
	}

	if !utils.IsValidSearch(query) {
		utils.SendError(&c.Controller, "Query too long", 400)
		return
	}

	suggestions, err := svc().CountryService.GetSearchSuggestions(query)
	if err != nil {
		utils.SendError(&c.Controller, "Failed to fetch suggestions", 500)
		return
	}

	utils.SendSuccess(&c.Controller, suggestions, "", 200)
}
