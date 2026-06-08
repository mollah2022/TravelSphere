package filters

import (
	"github.com/beego/beego/v2/server/web/context"
)

// AuthFilter checks whether the user is authenticated or not.
// If not authenticated, it redirects the user to the login page.
func AuthFilter(ctx *context.Context) {
	if ctx == nil || ctx.Input == nil {
		return
	}
	if ctx.Input.CruSession == nil {
		redirectURL := "/login?redirect=" + ctx.Input.URI()
		ctx.Redirect(302, redirectURL)
		return
	}

	sess := ctx.Input.Session("username")
	if sess == nil {
		redirectURL := "/login?redirect=" + ctx.Input.URI()
		ctx.Redirect(302, redirectURL)
	}
}
