package controllers

import (
	"TravelSphere/utils"
	"log"
	"strings"
)

// CountryController handles server-side rendered country routes.
type CountryController struct {
	BaseController
}

// List handles GET /countries and renders the Country Explorer page.
// Initially loads and displays all countries using SSR.
func (c *CountryController) List() {
	countries, err := svc().CountryService.GetAllCountries()
	if err != nil {
		log.Printf("[ERROR] CountryController.List: %v", err)
		c.Data["Error"] = "Failed to load countries. Please try again."
		countries = nil
	}

	regions := []string{"All Regions", "Africa", "Americas", "Asia", "Europe", "Oceania"}

	c.Data["Countries"] = countries
	c.Data["Regions"] = regions
	c.Data["TotalCount"] = len(countries)
	c.Data["PageTitle"] = "Country Explorer"
	c.TplName = "countries.tpl"
}

// Detail handles GET /countries/:slug and renders the destination detail page.
// Example: /countries/bangladesh → full page for Bangladesh.
func (c *CountryController) Detail() {
	slug := c.Ctx.Input.Param(":slug")
	slug = strings.ToLower(strings.TrimSpace(slug))
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	if !utils.IsValidSlug(slug) {
		c.renderNotFound("Invalid country URL.")
		return
	}

	country, err := svc().CountryService.GetCountryBySlug(slug)
	if err != nil {
		log.Printf("[ERROR] CountryController.Detail: slug=%s err=%v", slug, err)
		c.renderNotFound("Country not found.")
		return
	}

	attractions, err := svc().AttractionService.GetAttractionsByCountry(
		country.Latitude,
		country.Longitude,
	)
	if err != nil {
		log.Printf("[WARN] failed to get attractions for %s: %v", slug, err)
		attractions = nil
	}

	weather := svc().WeatherService.GetWeather(country.Capital)

	formattedPop := utils.FormatPopulation(country.Population)

	c.Data["Country"] = country
	c.Data["Attractions"] = attractions
	c.Data["Weather"] = weather
	c.Data["FormattedPopulation"] = formattedPop
	c.Data["PageTitle"] = country.Name
	c.TplName = "destination.tpl"
}

// renderNotFound displays a user-friendly 404 page.
func (c *CountryController) renderNotFound(msg string) {
	c.Ctx.ResponseWriter.WriteHeader(404)
	c.Data["ErrorMessage"] = msg
	c.Data["PageTitle"] = "Not Found"
	c.TplName = "error.tpl"
}
