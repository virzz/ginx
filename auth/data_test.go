package auth_test

import (
	"encoding/json"
	"testing"

	"github.com/virzz/ginx/auth"
)

var data = auth.DefaultData[string]{
	Token_:   "token",
	ID_:      "id",
	Account_: "account",
	Roles_:   []string{"admin"},
	Items_: map[string]any{
		"key": "value",
	},
}

func TestDataJSON(t *testing.T) {
	buf, err := json.Marshal(&data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(buf))
	v := &auth.DefaultData[string]{}
	err = json.Unmarshal(buf, v)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)
}

func TestDataSliceBinary(t *testing.T) {
	buf, err := data.Roles_.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(buf))
	v := &auth.DefaultData[string]{}
	err = v.Roles_.UnmarshalBinary(buf)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v.Roles_)
}

func TestDataMapBinary(t *testing.T) {
	buf, err := data.Items_.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(buf))
	v := &auth.DefaultData[string]{}
	err = v.Items_.UnmarshalText(buf)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v.Items_)
}
