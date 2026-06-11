package controllers

import (
	"TravelSphere/services"

	"github.com/beego/beego/v2/server/web"
)

// BaseController is the parent controller for all SSR controllers.
type BaseController struct {
	web.Controller
	IsLoggedIn bool
	Username   string
}

// Prepare runs before every request.
// It sets common data for all pages such as navigation and login state.
func (c *BaseController) Prepare() {

	// Get username from cookie (more reliable than session)
	defer func() {
		if r := recover(); r != nil {
			// Cookie access failed, continue without user data
		}
	}()

	cookie, err := c.Ctx.Request.Cookie("travelsphere_user")
	if err == nil && cookie.Value != "" {
		c.IsLoggedIn = true
		c.Username = cookie.Value
	}

	// Ensure template data map is initialized
	if c.Data == nil {
		c.Data = make(map[interface{}]interface{})
	}

	c.Data["IsLoggedIn"] = c.IsLoggedIn
	c.Data["Username"] = c.Username
	c.Data["AppName"] = "TravelSphere"

	// Set active navigation item.
	// For example, on /countries page, the Countries nav item will be active.
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

// svc provides access to the service container.
func svc() *services.ServiceContainer {
	return services.Container
}
