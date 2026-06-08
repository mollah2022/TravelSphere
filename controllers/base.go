package controllers

import (
	"TravelSphere/services"

	"github.com/beego/beego/v2/server/web"
)

// BaseController সব SSR controller এর parent
type BaseController struct {
	web.Controller
	IsLoggedIn bool
	Username   string
}

// Prepare প্রতিটা request এর আগে চলে
// সব page এ common data set করে — navigation, login state ইত্যাদি
func (c *BaseController) Prepare() {
	// Session থেকে username নাও
	var username interface{}
	if c.Ctx != nil && c.Ctx.Input != nil && c.Ctx.Input.CruSession != nil {
		username = c.GetSession("username")
	}
	if username != nil {
		c.IsLoggedIn = true
		c.Username = username.(string)
	}

	// Ensure template data map is initialized
	if c.Data == nil {
		c.Data = make(map[interface{}]interface{})
	}

	// Template এ pass করো
	c.Data["IsLoggedIn"] = c.IsLoggedIn
	c.Data["Username"] = c.Username
	c.Data["AppName"] = "TravelSphere"

	// Active nav item set করো
	// যেমন /countries page এ থাকলে Countries nav active হবে
	path := ""
	if c.Ctx != nil && c.Ctx.Request != nil {
		path = c.Ctx.Request.URL.Path
	}
	c.Data["NavHome"] = path == "/"
	c.Data["NavCountries"] = len(path) >= 10 &&
		path[:10] == "/countries"
	c.Data["NavWishlist"] = path == "/wishlist"
	c.Data["NavDashboard"] = path == "/dashboard"

	// Layout set করো
	c.Layout = "layout/main.tpl"
	c.LayoutSections = map[string]string{
		"Header": "partial/header.tpl",
		"Footer": "partial/footer.tpl",
	}
}

// svc service container এ access দেয়
func svc() *services.ServiceContainer {
	return services.Container
}