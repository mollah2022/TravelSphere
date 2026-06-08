package controllers

import (
	"os"
	"strings"
)

// AuthController login/logout handle করে
type AuthController struct {
	BaseController
}

// ShowLogin GET /login — login page দেখায়
func (c *AuthController) ShowLogin() {
	// Already logged in হলে home এ redirect করো
	if c.IsLoggedIn {
		c.Redirect("/", 302)
		return
	}

	// কোথা থেকে redirect হয়ে এসেছে সেটা save করো
	// login এর পরে সেখানে ফিরে যাবে
	redirect := c.GetString("redirect")
	if redirect == "" {
		redirect = "/"
	}

	c.Data["RedirectURL"] = redirect
	c.Data["PageTitle"] = "Login"
	// Login page এ layout থাকবে কিন্তু header/footer minimal
	c.TplName = "login.tpl"
}

// DoLogin POST /login — login process করে
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

	// Username sanitize করো
	username = sanitizeUsername(username)

	// Session এ save করো
	c.SetSession("username", username)

	// Redirect করো
	c.Redirect(redirect, 302)
}

// DoLogout GET /logout — session destroy করে
func (c *AuthController) DoLogout() {
	c.DestroySession()
	c.Redirect("/", 302)
}

// sanitizeUsername username থেকে special characters remove করে
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

// getDemoCredentials .env থেকে demo credentials আনে
// (unused for now — no password required)
func getDemoCredentials() (string, string) {
	return os.Getenv("DEMO_USERNAME"), os.Getenv("DEMO_PASSWORD")
}