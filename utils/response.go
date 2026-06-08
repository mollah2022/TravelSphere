package utils

import (
	"TravelSphere/models"

	"github.com/beego/beego/v2/server/web/context"
)

// JSONSuccess sends a standardized successful JSON response
func JSONSuccess(ctx *context.Context, data interface{}, message string) {
	ctx.Output.SetStatus(200)
	ctx.Output.JSON(models.NewSuccessResponse(data, message), false, false)
}

// JSONError sends a standardized error JSON response
func JSONError(ctx *context.Context, message string, code int) {
	ctx.Output.SetStatus(code)
	ctx.Output.JSON(models.NewErrorResponse(message, code), false, false)
}
