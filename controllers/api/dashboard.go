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
	if c.Ctx == nil || c.Ctx.Request == nil {
		utils.SendError(&c.Controller, "Unauthorized", 401)
		return
	}

	cookie, err := c.Ctx.Request.Cookie("travelsphere_user")
	if err != nil || cookie.Value == "" {
		utils.SendError(&c.Controller, "Unauthorized", 401)
		return
	}
	username := cookie.Value

	summary := svc().DashboardService.GetSummary(username)
	utils.SendSuccess(&c.Controller, summary, "", 200)
}
