package controllers

type WishlistController struct {
	BaseController
}

// WishlistController handles the wishlist page.
func (c *WishlistController) Get() {
	items := svc().WishlistService.GetWishlist(c.Username)

	c.Data["WishlistItems"] = items
	c.Data["PageTitle"] = "Travel Wishlist"
	c.TplName = "wishlist.tpl"
}
