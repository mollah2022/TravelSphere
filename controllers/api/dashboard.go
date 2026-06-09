package apicontrollers

import "TravelSphere/utils"

// DashboardAPIController JSON API for dashboard
type DashboardAPIController struct {
	CountriesAPIController
}

// Summary handles GET /api/dashboard/summary.
// Used for AJAX dashboard stats refresh.
// Returns: { total, planned, visited }.
func (c *DashboardAPIController) Summary() {
	if c.Ctx == nil || c.Ctx.Input == nil || c.Ctx.Input.CruSession == nil {
		utils.SendError(&c.Controller, "Unauthorized", 401)
		return
	}

	sess := c.GetSession("username")
	if sess == nil {
		utils.SendError(&c.Controller, "Unauthorized", 401)
		return
	}
	username := sess.(string)

	summary := svc().DashboardService.GetSummary(username)
	utils.SendSuccess(&c.Controller, summary, "", 200)
}
