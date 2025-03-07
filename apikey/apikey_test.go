package apikey_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

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
		c.JSON(200, apikey.Default(c).Data())
	})
	r.GET("/login", func(c *gin.Context) {
		v := apikey.Default(c)
		v.Set("id", "1")
		v.Set("name", "test")
		v.Set("group_id", "waefwaerfws")
		v.Save()
		c.String(200, v.Token())
	})
	err = r.Run(":8081")
	if err != nil {
		t.Fatal(err)
	}

}

func TestStoreDataStringSlice(t *testing.T) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", DB: 7})
	store, err := apikey.NewRedisStore(apikey.WithClient(client))
	if err != nil {
		t.Fatal(err)
	}
	apikey.Init(store)
	d := &apikey.DefaultData{}
	d.New()
	d.SetRoles([]string{"admin"})
	err = store.Save(ctx, d, 1*time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	dx := &apikey.DefaultData{}
	err = store.Get(ctx, dx)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(d.Roles())
}

func TestStoreDataMap(t *testing.T) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", DB: 7})
	store, err := apikey.NewRedisStore(apikey.WithClient(client))
	if err != nil {
		t.Fatal(err)
	}
	apikey.Init(store)
	type DataMap struct {
		apikey.DefaultData
		Values apikey.DataMap `json:"values" redis:"values"`
	}
	d := &DataMap{
		Values: map[string]any{
			"a": 1,
			"b": "2",
		},
	}
	d.New()
	d.SetRoles([]string{"admin"})
	err = store.Save(ctx, d, 1*time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func TestA(t *testing.T) {
	fmt.Println(strings.Split("admin", ","))
}
