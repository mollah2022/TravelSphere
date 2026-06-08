package controllers

import (
	"TravelSphere/utils"
	"log"
	"strings"
)

// CountryController SSR country routes handle করে
type CountryController struct {
	BaseController
}

// List GET /countries — Country Explorer page render করে
// Initial load এ সব দেশ SSR এ দেখায়
func (c *CountryController) List() {
	countries, err := svc().CountryService.GetAllCountries()
	if err != nil {
		log.Printf("[ERROR] CountryController.List: %v", err)
		c.Data["Error"] = "Failed to load countries. Please try again."
		countries = nil
	}

	// Region list for filter dropdown
	regions := []string{"All Regions", "Africa", "Americas", "Asia", "Europe", "Oceania"}

	c.Data["Countries"] = countries
	c.Data["Regions"] = regions
	c.Data["TotalCount"] = len(countries)
	c.Data["PageTitle"] = "Country Explorer"
	c.TplName = "countries.tpl"
}

// Detail GET /countries/:slug — Destination detail page render করে
// যেমন: /countries/bangladesh → Bangladesh এর full page
func (c *CountryController) Detail() {
	slug := c.Ctx.Input.Param(":slug")
	slug = strings.ToLower(strings.TrimSpace(slug))
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// Slug validate করো
	if !utils.IsValidSlug(slug) {
		c.renderNotFound("Invalid country URL.")
		return
	}

	// Country খোঁজো
	country, err := svc().CountryService.GetCountryBySlug(slug)
	if err != nil {
		log.Printf("[ERROR] CountryController.Detail: slug=%s err=%v", slug, err)
		c.renderNotFound("Country not found.")
		return
	}

	// Attractions আনো (lat/lon দিয়ে)
	attractions, err := svc().AttractionService.GetAttractionsByCountry(
		country.Latitude,
		country.Longitude,
	)
	if err != nil {
		log.Printf("[WARN] failed to get attractions for %s: %v", slug, err)
		attractions = nil
	}

	// Weather আনো (bonus — capital city দিয়ে)
	weather := svc().WeatherService.GetWeather(country.Capital)

	// Population formatted
	formattedPop := utils.FormatPopulation(country.Population)

	c.Data["Country"] = country
	c.Data["Attractions"] = attractions
	c.Data["Weather"] = weather
	c.Data["FormattedPopulation"] = formattedPop
	c.Data["PageTitle"] = country.Name
	c.TplName = "destination.tpl"
}

// renderNotFound user-friendly 404 page দেখায়
func (c *CountryController) renderNotFound(msg string) {
	c.Ctx.ResponseWriter.WriteHeader(404)
	c.Data["ErrorMessage"] = msg
	c.Data["PageTitle"] = "Not Found"
	c.TplName = "error.tpl"
}
