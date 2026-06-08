package controllers

import (
	"os"
	"strings"
)

// AuthController handles login and logout
type AuthController struct {
	BaseController
}

// ShowLogin shows the login page
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

// DoLogin handles login form submission
func (c *AuthController) DoLogin() {
	username := strings.TrimSpace(c.GetString("username"))
	redirect := c.GetString("redirect")
	if redirect == "" {
		redirect = "/"
	}

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

	c.SetSession("username", username)
	c.Redirect(redirect, 302)
}

// DoLogout logs out the user
func (c *AuthController) DoLogout() {
	c.DestroySession()
	c.Redirect("/", 302)
}

// sanitizeUsername removes unwanted characters from username
func sanitizeUsername(s string) string {
	var result strings.Builder
	for _, ch := range s {
		if (ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '_' || ch == '-' {
			result.WriteRune(ch)
		}
	}
	return result.String()
}

// getDemoCredentials gets demo username and password from environment
func getDemoCredentials() (string, string) {
	return os.Getenv("DEMO_USERNAME"), os.Getenv("DEMO_PASSWORD")
}
