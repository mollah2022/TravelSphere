package apicontrollers

import (
	"TravelSphere/models"
	"TravelSphere/services"
	"TravelSphere/utils"
	"encoding/json"

	"github.com/beego/beego/v2/server/web"
)

type WishlistAPIController struct {
	web.Controller
}

// username extracts the logged-in username from session.
// It returns empty string if session is not available.
func (c *WishlistAPIController) username() string {
	if c.Ctx.Input.CruSession == nil {
		return ""
	}
	sess := c.Ctx.Input.Session("username")
	if sess == nil {
		return ""
	}
	return sess.(string)
}

// List returns all wishlist items for the logged-in user in JSON format.
func (c *WishlistAPIController) List() {
	username := c.username()
	if username == "" {
		utils.JSONError(c.Ctx, "Unauthorized", 401)
		return
	}

	items := services.Container.WishlistService.GetWishlist(username)
	if items == nil {
		items = []*models.WishlistItem{}
	}
	utils.JSONSuccess(c.Ctx, items, "")
}

// Create adds a new item to the user's wishlist.
// It validates request body and returns the created item.
func (c *WishlistAPIController) Create() {
	username := c.username()
	if username == "" {
		utils.JSONError(c.Ctx, "Unauthorized", 401)
		return
	}

	var req models.CreateWishlistRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		utils.JSONError(c.Ctx, "Invalid request body", 400)
		return
	}

	item, err := services.Container.WishlistService.Create(username, &req)
	if err != nil {
		utils.JSONError(c.Ctx, err.Error(), 400)
		return
	}

	utils.JSONSuccess(c.Ctx, item, "Item added to wishlist")
}

// Update modifies an existing wishlist item for the logged-in user.
// It checks item ownership, validates input, and updates data.
func (c *WishlistAPIController) Update() {
	username := c.username()
	if username == "" {
		utils.JSONError(c.Ctx, "Unauthorized", 401)
		return
	}

	id := c.Ctx.Input.Param(":id")
	if id == "" {
		utils.JSONError(c.Ctx, "Missing item ID", 400)
		return
	}

	var req models.UpdateWishlistRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		utils.JSONError(c.Ctx, "Invalid request body", 400)
		return
	}

	item, err := services.Container.WishlistService.Update(username, id, &req)
	if err != nil {
		if err == models.ErrNotFound {
			utils.JSONError(c.Ctx, "Item not found", 404)
			return
		}
		if err == models.ErrUnauthorized {
			utils.JSONError(c.Ctx, "Unauthorized", 403)
			return
		}
		utils.JSONError(c.Ctx, err.Error(), 400)
		return
	}

	utils.JSONSuccess(c.Ctx, item, "Item updated")
}

// Delete removes a wishlist item for the logged-in user.
// It validates ownership and deletes the item if it exists.
func (c *WishlistAPIController) Delete() {
	username := c.username()
	if username == "" {
		utils.JSONError(c.Ctx, "Unauthorized", 401)
		return
	}

	id := c.Ctx.Input.Param(":id")
	if id == "" {
		utils.JSONError(c.Ctx, "Missing item ID", 400)
		return
	}

	err := services.Container.WishlistService.Delete(username, id)
	if err != nil {
		if err == models.ErrNotFound {
			utils.JSONError(c.Ctx, "Item not found", 404)
			return
		}
		if err == models.ErrUnauthorized {
			utils.JSONError(c.Ctx, "Unauthorized", 403)
			return
		}
		utils.JSONError(c.Ctx, err.Error(), 400)
		return
	}

	utils.JSONSuccess(c.Ctx, nil, "Item deleted")
}
