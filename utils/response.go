package utils

import (
	"TravelSphere/models"

	"github.com/beego/beego/v2/server/web"
)

// SendSuccess controller থেকে success JSON response পাঠায়
func SendSuccess(c *web.Controller, data interface{}, message string, statusCode int) {
	c.Ctx.ResponseWriter.WriteHeader(statusCode)
	c.Data["json"] = models.NewSuccessResponse(data, message)
	c.ServeJSON()
}

// SendError controller থেকে error JSON response পাঠায়
func SendError(c *web.Controller, message string, statusCode int) {
	c.Ctx.ResponseWriter.WriteHeader(statusCode)
	c.Data["json"] = models.NewErrorResponse(message, statusCode)
	c.ServeJSON()
}