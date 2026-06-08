package routers

import (
	apicontrollers "TravelSphere/controllers/api"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	web.Router("/api/countries",
		&apicontrollers.CountriesAPIController{},
		"get:List",
	)
	web.Router("/api/countries/:slug",
		&apicontrollers.CountriesAPIController{},
		"get:Detail",
	)
	web.Router("/api/attractions",
		&apicontrollers.CountriesAPIController{},
		"get:Attractions",
	)
	web.Router("/api/suggestions",
		&apicontrollers.CountriesAPIController{},
		"get:Suggestions",
	)
	web.Router("/api/wishlist",
		&apicontrollers.WishlistAPIController{},
		"get:List;post:Create",
	)
	web.Router("/api/wishlist/:id",
		&apicontrollers.WishlistAPIController{},
		"put:Update;delete:Delete",
	)
	web.Router("/api/dashboard/summary",
		&apicontrollers.DashboardAPIController{},
		"get:Summary",
	)
}
