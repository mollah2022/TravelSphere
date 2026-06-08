package apicontrollers

import (
	"TravelSphere/services"
	"TravelSphere/utils"

	"github.com/beego/beego/v2/server/web"
)

type DashboardAPIController struct {
	web.Controller
}

// It returns user-specific dashboard summary in JSON format.
func (c *DashboardAPIController) Summary() {
	var username string
	if c.Ctx.Input.CruSession != nil {
		if sess := c.Ctx.Input.Session("username"); sess != nil {
			username = sess.(string)
		}
	}

	if username == "" {
		utils.JSONError(c.Ctx, "Unauthorized", 401)
		return
	}

	summary := services.Container.DashboardService.GetSummary(username)
	utils.JSONSuccess(c.Ctx, summary, "")
}
