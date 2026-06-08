package controllers

// WishlistController wishlist SSR page handle করে
// Route: GET /wishlist (protected)
type WishlistController struct {
	BaseController
}

// Get wishlist page render করে
// Auth filter এর পরে এখানে আসা মানে user logged in
func (c *WishlistController) Get() {
	items := svc().WishlistService.GetWishlist(c.Username)

	c.Data["WishlistItems"] = items
	c.Data["PageTitle"] = "Travel Wishlist"
	c.TplName = "wishlist.tpl"
}