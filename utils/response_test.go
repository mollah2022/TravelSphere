package utils_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"TravelSphere/utils"

	"github.com/beego/beego/v2/server/web"
	ctxpkg "github.com/beego/beego/v2/server/web/context"
)

func newControllerRecorder() (*web.Controller, *httptest.ResponseRecorder) {
	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/json-test", nil)
	ctx := ctxpkg.NewContext()
	ctx.Reset(rw, req)

	return &web.Controller{Ctx: ctx, Data: make(map[interface{}]interface{})}, rw
}

func TestSendSuccess(t *testing.T) {
	c, rw := newControllerRecorder()
	utils.SendSuccess(c, map[string]string{"hello": "world"}, "ok", 201)

	if rw.Code != 201 {
		t.Fatalf("expected status 201, got %d", rw.Code)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(rw.Body.Bytes(), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if payload["success"] != true {
		t.Fatal("expected success true")
	}
}

func TestSendError(t *testing.T) {
	c, rw := newControllerRecorder()
	utils.SendError(c, "failed", 400)

	if rw.Code != 400 {
		t.Fatalf("expected status 400, got %d", rw.Code)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(rw.Body.Bytes(), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if payload["success"] != false {
		t.Fatal("expected success false")
	}
}
