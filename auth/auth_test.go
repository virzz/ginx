package auth_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/virzz/ginx/auth"
	"github.com/virzz/ginx/auth/token"
	"github.com/virzz/vlog"
)

func doReq(s *gin.Engine, req *http.Request) (*http.Response, []byte, error) {
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	rsp := w.Result()
	buf, err := io.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	return rsp, buf, err
}

func newTestServer(c *auth.Config) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(auth.Init(c))
	return r
}

func TestTokenRedis(t *testing.T) {
	r := newTestServer(&auth.Config{
		Type: auth.AuthTypeToken,
		Host: "localhost",
		Port: 6379,
	})
	r.GET("/login", func(c *gin.Context) {
		v := token.Default(c)
		v.SetID("1")
		v.SetAccount("test")
		v.SetRoles([]string{"admin"})
		v.SetValues("aaaa", "aaaaa")
		v.SetValues("vvvv", "asdveasd")
		v.Save()
		c.String(200, v.Token())
	})
	r.GET("/info", func(c *gin.Context) {
		v := token.Default(c).Data()
		buf, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			c.String(500, "Marshal"+err.Error())
			return
		}
		t.Log(string(buf))
		err = json.Unmarshal(buf, v)
		if err != nil {
			c.String(500, "Unmarshal"+err.Error())
			return
		}
		c.JSON(200, v)
	}, token.AuthMW())
	r.GET("/admin", func(c *gin.Context) {
		v := token.Default(c).Data()
		c.JSON(200, v)
	}, token.AuthMW(), token.RoleMW("admin"))
	//构建返回值
	req, _ := http.NewRequest("GET", "/login", nil)
	_, body, err := doReq(r, req)
	if err != nil {
		t.Fatal(err)
	}
	token := string(body)
	t.Log(token)
	req, _ = http.NewRequest("GET", "/info", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	_, body, err = doReq(r, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}

func TestSessionRedis(t *testing.T) {
	r := newTestServer(&auth.Config{Type: auth.AuthTypeSession,
		Host: "localhost", Port: 6379, Secret: "agrdeagedsrgefa"})
	r.GET("/login", func(c *gin.Context) {
		sess := sessions.Default(c)
		sess.Set("id", "1234")
		sess.Set("account", "test")
		sess.Set("roles", []string{"admin"})
		err := sess.Save()
		if err != nil {
			c.String(500, err.Error())
			return
		}
		c.String(200, sess.ID())
	})
	r.GET("/info", func(c *gin.Context) {
		sess := sessions.Default(c)
		account, ok := sess.Get("account").(string)
		if !ok {
			t.Fatal("account not found")
		}
		roles := sess.Get("roles")
		t.Log(roles)
		c.JSON(200, gin.H{"id": sess.ID(), "account": account, "roles": roles})
	})
	//构建返回值
	req, _ := http.NewRequest("GET", "/login", nil)
	rsp, body, err := doReq(r, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
	req, _ = http.NewRequest("GET", "/info", nil)
	for _, c := range rsp.Cookies() {
		req.AddCookie(c)
	}
	_, body, err = doReq(r, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}

func TestSessionCookie(t *testing.T) {
	r := newTestServer(&auth.Config{
		Type:   auth.AuthTypeCookie,
		Secret: "agrdeagedsrgefa",
	})
	r.GET("/login", func(c *gin.Context) {
		v := auth.Default(c)
		v.SetID("1")
		v.SetAccount("test")
		v.SetRoles([]string{"admin"})
		v.SetValues("aaaa", "aaaaa")
		v.SetValues("vvvv", "asdveasd")
		v.Save()
		c.String(200, v.Token())
	})
	r.GET("/info", func(c *gin.Context) {
		v := auth.Default(c)
		fmt.Printf("%#v", v)
		c.JSON(200, gin.H{"id": v.ID(), "account": v.Account(), "roles": v.Roles()})
	})
	//构建返回值
	req, _ := http.NewRequest("GET", "/login", nil)
	rsp, body, err := doReq(r, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
	req, _ = http.NewRequest("GET", "/info", nil)
	for _, c := range rsp.Cookies() {
		req.AddCookie(c)
	}
	_, body, err = doReq(r, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}

func TestDefault(t *testing.T) {
	r := newTestServer(&auth.Config{Type: auth.AuthTypeToken, Host: "localhost", Port: 6379})
	r.GET("/login", func(c *gin.Context) {
		v := auth.Default(c)
		v.SetID("1")
		v.SetAccount("test")
		v.SetRoles([]string{"admin"})
		v.SetValues("aaaa", "aaaaa")
		v.SetValues("vvvv", "asdveasd")
		v.Save()
		c.String(200, v.Token())
	})
	r.GET("/info", func(c *gin.Context) {
		v := auth.Default(c)
		vlog.Info(v.Token())
		vlog.Info(v.ID())
		vlog.Info(v.Account())
		vlog.Info(strings.Join(v.Roles(), ","))
		buf, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			c.String(500, "Marshal"+err.Error())
			return
		}
		c.Data(200, "application/json", buf)
	}, token.AuthMW())
	r.GET("/admin", func(c *gin.Context) { c.String(200, "ok") }, token.AuthMW(), token.RoleMW("admin"))
	//构建返回值
	t.Log("Req /login")
	req, _ := http.NewRequest("GET", "/login", nil)
	_, body, err := doReq(r, req)
	if err != nil {
		t.Fatal(err)
	}
	token := string(body)
	t.Log(token)
	t.Log("Req /info")
	req, _ = http.NewRequest("GET", "/info", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	_, body, err = doReq(r, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
	t.Log("Req /admin")
	req.URL.Path = "/admin"
	_, body, err = doReq(r, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}
