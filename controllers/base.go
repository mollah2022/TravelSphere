package controllers

import (
	"TravelSphere/services"

	"github.com/beego/beego/v2/server/web"
)

// BaseController is the main controller used by all other controllers
// It contains common data and logic
type BaseController struct {
	web.Controller
	IsLoggedIn bool
	Username   string
}

// Prepare runs before every request
// It sets common data for all pages
func (c *BaseController) Prepare() {
	var username interface{}
	if c.Ctx != nil && c.Ctx.Input != nil && c.Ctx.Input.CruSession != nil {
		username = c.GetSession("username")
	}
	if username != nil {
		c.IsLoggedIn = true
		c.Username = username.(string)
	}

	if c.Data == nil {
		c.Data = make(map[interface{}]interface{})
	}

	c.Data["IsLoggedIn"] = c.IsLoggedIn
	c.Data["Username"] = c.Username
	c.Data["AppName"] = "TravelSphere"

	path := ""
	if c.Ctx != nil && c.Ctx.Request != nil {
		path = c.Ctx.Request.URL.Path
	}
	c.Data["NavHome"] = path == "/"
	c.Data["NavCountries"] = len(path) >= 10 &&
		path[:10] == "/countries"
	c.Data["NavWishlist"] = path == "/wishlist"
	c.Data["NavDashboard"] = path == "/dashboard"

	c.Layout = "layout/main.tpl"
	c.LayoutSections = map[string]string{
		"Header": "partial/header.tpl",
		"Footer": "partial/footer.tpl",
	}
}

// svc returns service container (used to access services easily)
func svc() *services.ServiceContainer {
	return services.Container
}
