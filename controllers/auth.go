package controllers

import (
	"net/http"
	"os"
	"strings"
)

// AuthController handles user authentication including login and logout operations.
type AuthController struct {
	BaseController
}

// ShowLogin handles GET /login and displays the login page.
func (c *AuthController) ShowLogin() {

	if c.IsLoggedIn {
		c.Redirect("/", 302)
		return
	}

	redirect := c.GetString("redirect")
	if redirect == "" {
		redirect = "/"
	}

	c.Data["RedirectURL"] = redirect
	c.Data["PageTitle"] = "Login"

	c.TplName = "login.tpl"
}

// DoLogin POST /login — login process
func (c *AuthController) DoLogin() {
	username := strings.TrimSpace(c.GetString("username"))
	redirect := c.GetString("redirect")
	if redirect == "" {
		redirect = "/"
	}

	// Username validation
	if username == "" {
		c.Data["Error"] = "Please enter a username."
		c.Data["RedirectURL"] = redirect
		c.Data["PageTitle"] = "Login"
		c.TplName = "login.tpl"
		return
	}

	if len(username) < 2 {
		c.Data["Error"] = "Username must be at least 2 characters."
		c.Data["RedirectURL"] = redirect
		c.Data["PageTitle"] = "Login"
		c.TplName = "login.tpl"
		return
	}

	if len(username) > 30 {
		c.Data["Error"] = "Username must be at most 30 characters."
		c.Data["RedirectURL"] = redirect
		c.Data["PageTitle"] = "Login"
		c.TplName = "login.tpl"
		return
	}

	username = sanitizeUsername(username)

	// Set cookie instead of session (more reliable)
	cookie := &http.Cookie{
		Name:     "travelsphere_user",
		Value:    username,
		MaxAge:   86400, // 24 hours
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(c.Ctx.ResponseWriter, cookie)

	c.Redirect(redirect, 302)
}

// DoLogout handles GET /logout and destroys the user session.
func (c *AuthController) DoLogout() {
	// Clear the user cookie by setting MaxAge to -1
	cookie := &http.Cookie{
		Name:     "travelsphere_user",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(c.Ctx.ResponseWriter, cookie)
	c.Redirect("/", 302)
}

// sanitizeUsername removes special characters from the username.
func sanitizeUsername(s string) string {
	var result strings.Builder
	for _, ch := range s {
		// শুধু letters, numbers, underscore, hyphen allow করো
		if (ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '_' || ch == '-' {
			result.WriteRune(ch)
		}
	}
	return result.String()
}

// getDemoCredentials loads demo credentials from .env.
// (currently unused — no password required)
func getDemoCredentials() (string, string) {
	return os.Getenv("DEMO_USERNAME"), os.Getenv("DEMO_PASSWORD")
}
