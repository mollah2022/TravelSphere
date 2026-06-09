package controllers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/beego/beego/v2/server/web"
	ctxpkg "github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/session"
)

func ensureGlobalSessions() {
	if web.GlobalSessions == nil {
		mgr, err := session.NewManager("memory", session.NewManagerConfig(session.CfgSetCookie(true)))
		if err != nil {
			panic(err)
		}
		web.GlobalSessions = mgr
		go web.GlobalSessions.GC()
	}
}

func newAuthController(url string, session map[interface{}]interface{}) (*AuthController, *httptest.ResponseRecorder) {
	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", url, nil)
	ctx := ctxpkg.NewContext()
	ctx.Reset(rw, req)
	ctx.Input.CruSession = &mockStore{vals: session}
	ensureGlobalSessions()

	return &AuthController{BaseController: BaseController{Controller: web.Controller{Ctx: ctx, Data: make(map[interface{}]interface{})}}}, rw
}

func TestDoLogin_Success(t *testing.T) {
	c, rw := newAuthController("/login?username=Alice&redirect=/dashboard", map[interface{}]interface{}{})
	c.DoLogin()

	if rw.Code != 302 {
		t.Fatalf("expected redirect 302, got %d", rw.Code)
	}
	if rw.Header().Get("Location") != "/dashboard" {
		t.Fatalf("expected redirect to /dashboard, got %q", rw.Header().Get("Location"))
	}
	if c.GetSession("username") != "Alice" {
		t.Fatalf("expected session username Alice, got %#v", c.GetSession("username"))
	}
}

func TestDoLogin_InvalidUsername(t *testing.T) {
	c, _ := newAuthController("/login?username=a&redirect=/", map[interface{}]interface{}{})
	c.DoLogin()

	if c.TplName != "login.tpl" {
		t.Fatalf("expected login.tpl, got %q", c.TplName)
	}
	if c.Data["Error"] == nil {
		t.Fatal("expected error message for short username")
	}
}

func TestDoLogout_RedirectsToHome(t *testing.T) {
	c, rw := newAuthController("/logout", map[interface{}]interface{}{"username": "alice"})
	c.DoLogout()

	if rw.Code != 302 {
		t.Fatalf("expected redirect 302, got %d", rw.Code)
	}
	if rw.Header().Get("Location") != "/" {
		t.Fatalf("expected redirect to /, got %q", rw.Header().Get("Location"))
	}
	if c.Ctx.Input.CruSession != nil {
		t.Fatal("expected session to be cleared after logout")
	}
}

func TestShowLogin_RedirectsWhenLoggedIn(t *testing.T) {
	c, rw := newAuthController("/login", map[interface{}]interface{}{})
	c.IsLoggedIn = true
	c.ShowLogin()

	if rw.Code != 302 {
		t.Fatalf("expected redirect 302, got %d", rw.Code)
	}
	if rw.Header().Get("Location") != "/" {
		t.Fatalf("expected redirect to /, got %q", rw.Header().Get("Location"))
	}
}

func TestShowLogin_SetsDefaultRedirect(t *testing.T) {
	c, _ := newAuthController("/login", map[interface{}]interface{}{})
	c.ShowLogin()

	if c.TplName != "login.tpl" {
		t.Fatalf("expected login.tpl, got %q", c.TplName)
	}
	if c.Data["RedirectURL"] != "/" {
		t.Fatalf("expected default redirect /, got %#v", c.Data["RedirectURL"])
	}
}

func TestDoLogin_InvalidLongUsername(t *testing.T) {
	username := strings.Repeat("a", 31)
	c, _ := newAuthController("/login?username="+username+"&redirect=/", map[interface{}]interface{}{})
	c.DoLogin()

	if c.TplName != "login.tpl" {
		t.Fatalf("expected login.tpl, got %q", c.TplName)
	}
	if c.Data["Error"] == nil {
		t.Fatal("expected error message for long username")
	}
}

func TestDoLogin_SanitizesUsername(t *testing.T) {
	c, _ := newAuthController("/login?username=Alice!@#&redirect=/", map[interface{}]interface{}{})
	c.DoLogin()

	if c.GetSession("username") != "Alice" {
		t.Fatalf("expected sanitized username Alice, got %#v", c.GetSession("username"))
	}
}
