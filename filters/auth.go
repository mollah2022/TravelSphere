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

	// Check for authentication cookie instead of session
	cookie, err := ctx.Request.Cookie("travelsphere_user")
	if err != nil || cookie.Value == "" {
		redirectURL := "/login?redirect=" + ctx.Input.URI()
		ctx.Redirect(302, redirectURL)
		return
	}
}
