package controllers

import "log"

// HomeController handles the home page.
// Route: GET /
type HomeController struct {
	BaseController
}

// Get renders the home page.
// Displays featured countries and popular attractions.
func (c *HomeController) Get() {
	featured, err := svc().CountryService.GetFeaturedCountries()
	if err != nil {
		log.Printf("[ERROR] HomeController: failed to get featured countries: %v", err)
		featured = nil
	}

	attractions := svc().AttractionService.GetPopularAttractions()

	c.Data["FeaturedCountries"] = featured
	c.Data["PopularAttractions"] = attractions
	c.Data["PageTitle"] = "Discover Your Next Destination"
	c.TplName = "home.tpl"
}
