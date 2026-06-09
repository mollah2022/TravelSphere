package controllers

// WishlistController handles the server-side rendered wishlist page.
// Route: GET /wishlist (protected).
type WishlistController struct {
	BaseController
}

// Get renders the wishlist page.
// If this method is reached after auth filter, the user is authenticated.
func (c *WishlistController) Get() {
	items := svc().WishlistService.GetWishlist(c.Username)

	c.Data["WishlistItems"] = items
	c.Data["PageTitle"] = "Travel Wishlist"
	c.TplName = "wishlist.tpl"
}
