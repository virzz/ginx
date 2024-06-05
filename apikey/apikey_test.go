package apikey_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	apikey "github.com/virzz/ginx/apikey"
)

func TestToken(t *testing.T) {
	r := gin.Default()
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", DB: 7})
	store, err := apikey.NewRedisStore(apikey.WithClient(client))
	if err != nil {
		t.Fatal(err)
	}
	r.Use(apikey.Init(store))
	r.GET("/info", func(c *gin.Context) {
		c.JSON(200, apikey.Default(c).Values())
	})
	r.GET("/login", func(c *gin.Context) {
		v := apikey.Default(c)
		v.Set("id", "1")
		v.Set("name", "test")
		v.Save()
		c.String(200, v.Token())
	})
	err = r.Run(":8080")
	if err != nil {
		t.Fatal(err)
	}

}
