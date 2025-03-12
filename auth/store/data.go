package store

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"strings"

	"github.com/go-json-experiment/json"
)

func generateRandomKey() string {
	k := make([]byte, 20)
	io.ReadFull(rand.Reader, k)
	return hex.EncodeToString(k)
}

type (
	Data interface {
		Token() string
		ID() string
		Account() string
		Roles() []string
		SetToken(string) Data
		SetID(string) Data
		SetAccount(string) Data
		SetRoles([]string) Data

		New()
		Get(string) any
		Set(string, any) Data
	}

	DataStringSlice []string
	DataMap         map[string]any
)

func (d DataStringSlice) MarshalBinary() ([]byte, error) {
	return []byte(strings.Join(d, ",")), nil
}
func (d *DataStringSlice) UnmarshalBinary(buf []byte) error {
	*d = DataStringSlice(strings.Split(string(buf), ","))
	return nil
}
func (d *DataStringSlice) UnmarshalText(buf []byte) error {
	return d.UnmarshalBinary(buf)
}

func (d DataMap) MarshalBinary() ([]byte, error)    { return json.Marshal(d) }
func (d *DataMap) UnmarshalBinary(buf []byte) error { return json.Unmarshal(buf, d) }
func (d *DataMap) UnmarshalText(buf []byte) error   { return json.Unmarshal(buf, d) }
