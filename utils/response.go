package utils

import (
	"TravelSphere/models"

	"github.com/beego/beego/v2/server/web"
)

// SendSuccess sends a successful JSON response from the controller.
func SendSuccess(c *web.Controller, data interface{}, message string, statusCode int) {
	c.Ctx.ResponseWriter.WriteHeader(statusCode)
	c.Data["json"] = models.NewSuccessResponse(data, message)
	c.ServeJSON()
}

// SendError sends an error JSON response from the controller.
func SendError(c *web.Controller, message string, statusCode int) {
	c.Ctx.ResponseWriter.WriteHeader(statusCode)
	c.Data["json"] = models.NewErrorResponse(message, statusCode)
	c.ServeJSON()
}
