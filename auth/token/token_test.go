package token_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/virzz/ginx/auth/store"
	token "github.com/virzz/ginx/auth/token"
)

func TestToken(t *testing.T) {
	r := gin.Default()
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", DB: 7})
	store, err := store.NewRedisStore(client)
	if err != nil {
		t.Fatal(err)
	}
	r.Use(token.Init(store))
	r.GET("/info", func(c *gin.Context) {
		c.JSON(200, token.Default(c).Data())
	})
	r.GET("/login", func(c *gin.Context) {
		v := token.Default(c)
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

func TestStoreDataStringSlice(t *testing.T) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", DB: 7})
	s, err := store.NewRedisStore(client)
	if err != nil {
		t.Fatal(err)
	}
	token.Init(s)
	d := &store.DefaultData{}
	d.New()
	d.SetRoles([]string{"admin"})
	err = s.Save(ctx, d, 1*time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	dx := &store.DefaultData{}
	err = s.Get(ctx, dx)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(d.Roles())
}

func TestStoreDataMap(t *testing.T) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", DB: 7})
	s, err := store.NewRedisStore(client)
	if err != nil {
		t.Fatal(err)
	}
	token.Init(s)
	type DataMap struct {
		store.DefaultData
		Values store.DataMap `json:"values" redis:"values"`
	}
	d := &DataMap{
		Values: map[string]any{
			"a": 1,
			"b": "2",
		},
	}
	d.New()
	d.SetRoles([]string{"admin"})
	err = s.Save(ctx, d, 1*time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func TestA(t *testing.T) {
	fmt.Println(strings.Split("admin", ","))
}
