package controllers

import "log"

// HomeController controls the home page
type HomeController struct {
	BaseController
}

// Get runs when user opens home page (GET request)
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
