package apicontrollers

import (
	"TravelSphere/models"
	"TravelSphere/utils"
	"encoding/json"
)

// WishlistAPIController JSON API for wishlist CRUD
// সব route এ auth check করা হয়
type WishlistAPIController struct {
	CountriesAPIController // svc() inherit করার জন্য
}

// getUsername session থেকে username নেয়
// না থাকলে 401 response পাঠায় এবং false return করে
func (c *WishlistAPIController) getUsername() (string, bool) {
	if c.Ctx == nil || c.Ctx.Input == nil || c.Ctx.Input.CruSession == nil {
		utils.SendError(&c.Controller, "Unauthorized", 401)
		return "", false
	}

	sess := c.GetSession("username")
	if sess == nil {
		utils.SendError(&c.Controller, "Unauthorized", 401)
		return "", false
	}
	return sess.(string), true
}

// List GET /api/wishlist
// Logged in user এর সব wishlist items return করে
func (c *WishlistAPIController) List() {
	username, ok := c.getUsername()
	if !ok {
		return
	}

	items := svc().WishlistService.GetWishlist(username)
	utils.SendSuccess(&c.Controller, items, "", 200)
}

// Create POST /api/wishlist
// নতুন wishlist item তৈরি করে
// Body: { country_name, note, status }
func (c *WishlistAPIController) Create() {
	username, ok := c.getUsername()
	if !ok {
		return
	}

	// Request body parse করো
	var req models.CreateWishlistRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		utils.SendError(&c.Controller, "Invalid request body", 400)
		return
	}

	// XSS prevent করো
	req.CountryName = utils.SanitizeString(req.CountryName)
	req.Note = utils.SanitizeString(req.Note)

	// Service call করো
	item, err := svc().WishlistService.AddToWishlist(username, req)
	if err != nil {
		switch err {
		case models.ErrCountryNameRequired:
			utils.SendError(&c.Controller, err.Error(), 400)
		case models.ErrInvalidStatus:
			utils.SendError(&c.Controller, err.Error(), 400)
		default:
			utils.SendError(&c.Controller, "Failed to add to wishlist", 500)
		}
		return
	}

	utils.SendSuccess(&c.Controller, item, "Added to wishlist", 201)
}

// Update PUT /api/wishlist/:id
// Wishlist item এর note আর status update করে
// Body: { note, status }
func (c *WishlistAPIController) Update() {
	username, ok := c.getUsername()
	if !ok {
		return
	}

	id := c.Ctx.Input.Param(":id")
	if id == "" {
		utils.SendError(&c.Controller, "Invalid ID", 400)
		return
	}

	// Request body parse করো
	var req models.UpdateWishlistRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		utils.SendError(&c.Controller, "Invalid request body", 400)
		return
	}

	// XSS prevent করো
	req.Note = utils.SanitizeString(req.Note)

	// Service call করো
	item, err := svc().WishlistService.UpdateWishlistItem(username, id, req)
	if err != nil {
		switch err {
		case models.ErrUnauthorized:
			utils.SendError(&c.Controller, "Unauthorized", 403)
		case models.ErrNotFound:
			utils.SendError(&c.Controller, "Item not found", 404)
		case models.ErrInvalidStatus:
			utils.SendError(&c.Controller, err.Error(), 400)
		case models.ErrStatusRequired:
			utils.SendError(&c.Controller, err.Error(), 400)
		default:
			utils.SendError(&c.Controller, "Failed to update", 500)
		}
		return
	}

	utils.SendSuccess(&c.Controller, item, "Updated successfully", 200)
}

// Delete DELETE /api/wishlist/:id
// Wishlist item delete করে
func (c *WishlistAPIController) Delete() {
	username, ok := c.getUsername()
	if !ok {
		return
	}

	id := c.Ctx.Input.Param(":id")
	if id == "" {
		utils.SendError(&c.Controller, "Invalid ID", 400)
		return
	}

	// Service call করো
	err := svc().WishlistService.DeleteWishlistItem(username, id)
	if err != nil {
		switch err {
		case models.ErrUnauthorized:
			utils.SendError(&c.Controller, "Unauthorized", 403)
		case models.ErrNotFound:
			utils.SendError(&c.Controller, "Item not found", 404)
		default:
			utils.SendError(&c.Controller, "Failed to delete", 500)
		}
		return
	}

	utils.SendSuccess(&c.Controller, nil, "Deleted successfully", 200)
}
