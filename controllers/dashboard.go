package controllers

type DashboardController struct {
	BaseController
}

// DashboardController handles the user dashboard page.
func (c *DashboardController) Get() {
	username := c.Username

	summary := svc().DashboardService.GetSummary(username)
	destinations := svc().DashboardService.GetSavedDestinations(username)

	c.Data["Summary"] = summary
	c.Data["Destinations"] = destinations
	c.Data["PageTitle"] = "Travel Dashboard"
	c.TplName = "dashboard.tpl"
}
