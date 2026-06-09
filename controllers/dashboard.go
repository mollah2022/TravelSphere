package controllers

// DashboardController handles the server-side rendered dashboard page.
// Route: GET /dashboard (protected).
type DashboardController struct {
	BaseController
}

// Get renders the dashboard page.
func (c *DashboardController) Get() {
	username := c.Username

	summary := svc().DashboardService.GetSummary(username)

	destinations := svc().DashboardService.GetSavedDestinations(username)

	c.Data["Summary"] = summary
	c.Data["Destinations"] = destinations
	c.Data["PageTitle"] = "Travel Dashboard"
	c.TplName = "dashboard.tpl"
}
