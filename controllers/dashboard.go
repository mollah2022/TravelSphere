package controllers

// DashboardController dashboard SSR page handle করে
// Route: GET /dashboard (protected)
type DashboardController struct {
	BaseController
}

// Get dashboard page render করে
func (c *DashboardController) Get() {
	username := c.Username

	// Stats আনো
	summary := svc().DashboardService.GetSummary(username)

	// Saved destinations আনো
	destinations := svc().DashboardService.GetSavedDestinations(username)

	c.Data["Summary"] = summary
	c.Data["Destinations"] = destinations
	c.Data["PageTitle"] = "Travel Dashboard"
	c.TplName = "dashboard.tpl"
}