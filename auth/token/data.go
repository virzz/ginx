package token

import (
	"encoding"

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
		New() string
		Set(string, any) Data
		SetToken(string) Data
		SetID(string) Data
		SetAccount(string) Data
		SetValues(string, any) Data
		SetRoles([]string) Data
		Delete(string) Data
	}
)

var _ encoding.TextUnmarshaler = (*DataStringSlice)(nil)
var _ encoding.TextUnmarshaler = (*DataMap)(nil)

func (d DataStringSlice) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *DataStringSlice) UnmarshalBinary(buf []byte) error {
	v := []string{}
	err := json.Unmarshal(buf, &v)
	if err != nil {
		return err
	}
	*d = DataStringSlice(v)
	return nil
}

func (d *DataStringSlice) UnmarshalText(buf []byte) error {
	return d.UnmarshalBinary(buf)
}

func (d *DataStringSlice) UnmarshalJSON(buf []byte) error {
	return d.UnmarshalText(buf)
}

func (d DataMap) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *DataMap) UnmarshalBinary(buf []byte) error {
	return json.Unmarshal(buf, d)
}

func (d *DataMap) UnmarshalText(buf []byte) error {
	_d := make(map[string]any)
	if err := json.Unmarshal(buf, &_d); err != nil {
		return err
	}
	*d = _d
	return nil
}

func (d *DataMap) UnmarshalJSON(buf []byte) error {
	return d.UnmarshalText(buf)
}
