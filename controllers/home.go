package controllers

import "log"

// HomeController home page handle করে
// Route: GET /
type HomeController struct {
	BaseController
}

// Get home page render করে
// featured countries + popular attractions দেখায়
func (c *HomeController) Get() {
	// Featured countries আনো
	featured, err := svc().CountryService.GetFeaturedCountries()
	if err != nil {
		log.Printf("[ERROR] HomeController: failed to get featured countries: %v", err)
		// Error হলেও page render হবে, empty list দিয়ে
		featured = nil
	}

	// Popular attractions আনো (static list)
	attractions := svc().AttractionService.GetPopularAttractions()

	c.Data["FeaturedCountries"] = featured
	c.Data["PopularAttractions"] = attractions
	c.Data["PageTitle"] = "Discover Your Next Destination"
	c.TplName = "home.tpl"
}