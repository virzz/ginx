package apikey

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

func generateRandomKey() string {
	k := make([]byte, 20)
	io.ReadFull(rand.Reader, k)
	return hex.EncodeToString(k)
}

type Data struct {
	Token  string         `json:"token"`
	Values map[string]any `json:"values"`
}

func (v *Data) New() { v.Token = generateRandomKey() }
