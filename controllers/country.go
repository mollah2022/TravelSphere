package controllers

import (
	"TravelSphere/utils"
	"log"
	"strings"
)

type CountryController struct {
	BaseController
}

// List handles the country listing page.
// It fetches all countries from service and sends them to the template.
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

// Detail handles single country detail page.
// It validates slug, fetches country info, attractions, and weather data.
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

// renderNotFound shows a 404 error page with a custom message.
func (c *CountryController) renderNotFound(msg string) {
	c.Ctx.ResponseWriter.WriteHeader(404)
	c.Data["ErrorMessage"] = msg
	c.Data["PageTitle"] = "Not Found"
	c.TplName = "error.tpl"
}
