package store

import (
	"encoding"
	"strings"

	"github.com/go-json-experiment/json"
)

type (
	DataStringSlice []string
	DataMap         map[string]any
	Data            interface {
		Token() string
		ID() string
		Account() string
		Roles() []string
		Items() DataMap
		Get(string) any
		New()
		Set(string, any) Data
		SetToken(string) Data
		SetID(string) Data
		SetAccount(string) Data
		SetRoles([]string) Data
	}
)

var _ encoding.TextUnmarshaler = (*DataStringSlice)(nil)
var _ encoding.TextUnmarshaler = (*DataMap)(nil)

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
func (d *DataMap) UnmarshalText(buf []byte) error {
	_d := make(map[string]any)
	err := json.Unmarshal(buf, &_d)
	if err != nil {
		return err
	}
	*d = _d
	return nil
}
