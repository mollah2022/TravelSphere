package routers

import (
	"TravelSphere/controllers"
	"TravelSphere/filters"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	web.InsertFilter("/*", web.BeforeRouter, filters.LoggingFilter)

	web.InsertFilter("/wishlist", web.BeforeRouter, filters.AuthFilter)
	web.InsertFilter("/dashboard", web.BeforeRouter, filters.AuthFilter)

	web.Router("/", &controllers.HomeController{})
	web.Router("/countries", &controllers.CountryController{}, "get:List")
	web.Router("/countries/:slug", &controllers.CountryController{}, "get:Detail")
	web.Router("/wishlist", &controllers.WishlistController{}, "get:Get")
	web.Router("/dashboard", &controllers.DashboardController{}, "get:Get")

	web.Router("/login", &controllers.AuthController{},
		"get:ShowLogin;post:DoLogin")
	web.Router("/logout", &controllers.AuthController{}, "get:DoLogout")
}
