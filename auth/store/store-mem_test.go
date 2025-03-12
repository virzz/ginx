package store_test

import (
	"context"
	"testing"

	"github.com/virzz/ginx/auth/store"
)

func TestMemStore(t *testing.T) {
	s, err := store.NewMemStore()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	d := &store.DefaultData{}
	d.New()
	d.SetRoles([]string{"admin"})
	err = s.Save(ctx, d, 1)
	if err != nil {
		t.Fatal(err)
	}
	token := d.Token()
	dx := &store.DefaultData{Token_: token}
	err = s.Get(ctx, dx)
	if err != nil {
		t.Fatal(err)
	}
}
