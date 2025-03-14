package token_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	token "github.com/virzz/ginx/auth/token"
)

func TestTokenRedis(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	r.Use(token.Init(client))
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
	//构建返回值
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/login", nil)
	r.ServeHTTP(w1, req1)
	rsp1 := w1.Result()
	body, _ := io.ReadAll(rsp1.Body)
	token := string(body)
	t.Log(token)

	req2, _ := http.NewRequest("GET", "/info", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	rsp2 := w2.Result()
	body, _ = io.ReadAll(rsp2.Body)
	t.Log(string(body))
}
