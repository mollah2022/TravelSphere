package apicontrollers

import (
	"TravelSphere/services"
	"TravelSphere/utils"

	"github.com/beego/beego/v2/server/web"
)

type CountriesAPIController struct {
	web.Controller
}

// List returns all countries from the service in JSON format.
func (c *CountriesAPIController) List() {
	countries, err := services.Container.CountryService.GetAllCountries()
	if err != nil {
		utils.JSONError(c.Ctx, "Failed to fetch countries", 500)
		return
	}
	utils.JSONSuccess(c.Ctx, countries, "")
}

// Detail returns a single country based on slug validation and lookup.
func (c *CountriesAPIController) Detail() {
	slug := c.Ctx.Input.Param(":slug")
	if !utils.IsValidSlug(slug) {
		utils.JSONError(c.Ctx, "Invalid country slug", 400)
		return
	}

	country, err := services.Container.CountryService.GetCountryBySlug(slug)
	if err != nil {
		utils.JSONError(c.Ctx, "Country not found", 404)
		return
	}
	utils.JSONSuccess(c.Ctx, country, "")
}

// Attractions returns attractions for a specific country using its coordinates.
func (c *CountriesAPIController) Attractions() {
	slug := c.GetString("slug")
	if !utils.IsValidSlug(slug) {
		utils.JSONError(c.Ctx, "Invalid slug", 400)
		return
	}

	country, err := services.Container.CountryService.GetCountryBySlug(slug)
	if err != nil {
		utils.JSONError(c.Ctx, "Country not found", 404)
		return
	}

	attractions, err := services.Container.AttractionService.GetAttractionsByCountry(
		country.Latitude,
		country.Longitude,
	)
	if err != nil {
		utils.JSONError(c.Ctx, "Failed to fetch attractions", 500)
		return
	}
	utils.JSONSuccess(c.Ctx, attractions, "")
}

// Suggestions returns country search results based on a query string.
func (c *CountriesAPIController) Suggestions() {
	query := c.GetString("q")
	if query == "" {
		utils.JSONSuccess(c.Ctx, []interface{}{}, "")
		return
	}

	results := services.Container.CountryService.SearchCountries(query)
	utils.JSONSuccess(c.Ctx, results, "")
}
