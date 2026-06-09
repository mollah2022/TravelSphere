package apicontrollers

import (
	"TravelSphere/services"
	"TravelSphere/utils"

	"github.com/beego/beego/v2/server/web"
)

// CountriesAPIController JSON API for countries
// সব response JSON, কোনো HTML নেই
type CountriesAPIController struct {
	web.Controller
}

// svc service container access করে
func svc() *services.ServiceContainer {
	return services.Container
}

// List GET /api/countries
// Query params: search, region
// AJAX country search + filter এর জন্য ব্যবহার হয়
func (c *CountriesAPIController) List() {
	search := c.GetString("search")
	region := c.GetString("region")

	// Input validate করো
	if !utils.IsValidSearch(search) {
		utils.SendError(&c.Controller, "Search query too long", 400)
		return
	}

	if !utils.IsValidRegion(region) {
		utils.SendError(&c.Controller, "Invalid region", 400)
		return
	}

	// Service call করো
	countries, err := svc().CountryService.SearchCountries(search, region)
	if err != nil {
		utils.SendError(&c.Controller, "Failed to fetch countries", 500)
		return
	}

	utils.SendSuccess(&c.Controller, countries, "", 200)
}

// Detail GET /api/countries/:slug
// Single country JSON return করে
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

// Attractions GET /api/attractions
// Query params: lat, lon
// Destination page এ AJAX attractions load এর জন্য
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

// Suggestions GET /api/suggestions
// Query params: q
// Home page search autocomplete এর জন্য
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
