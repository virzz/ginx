package apikey

import (
	"crypto/rand"
	"encoding"
	"encoding/hex"
	"encoding/json"
	"io"
	"strings"
)

func generateRandomKey() string {
	k := make([]byte, 20)
	io.ReadFull(rand.Reader, k)
	return hex.EncodeToString(k)
}

type (
	DataStringSlice []string
	DataMap         map[string]any

	Data interface {
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

	DefaultData struct {
		Token_   string          `json:"token" redis:"token"`
		ID_      string          `json:"id" redis:"id"`
		Account_ string          `json:"account" redis:"account"`
		Roles_   DataStringSlice `json:"roles" redis:"roles"`
		Items_   DataMap         `json:"items" redis:"items"`
	}
)

var _ Data = (*DefaultData)(nil)
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

func (d *DefaultData) New()            { d.Token_ = generateRandomKey() }
func (d *DefaultData) Token() string   { return d.Token_ }
func (d *DefaultData) ID() string      { return d.ID_ }
func (d *DefaultData) Account() string { return d.Account_ }
func (d *DefaultData) Roles() []string { return []string(d.Roles_) }
func (d *DefaultData) Items() DataMap  { return d.Items_ }

func (d *DefaultData) SetToken(v string) Data {
	d.Token_ = v
	return d
}
func (d *DefaultData) SetID(v string) Data {
	d.ID_ = v
	return d
}
func (d *DefaultData) SetAccount(v string) Data {
	d.Account_ = v
	return d
}
func (d *DefaultData) SetRoles(v []string) Data {
	d.Roles_ = DataStringSlice(v)
	return d
}

func (d *DefaultData) Set(key string, val any) Data {
	switch key {
	case "id":
		if v, ok := val.(string); ok {
			d.ID_ = v
		}
	case "account":
		if v, ok := val.(string); ok {
			d.Account_ = v
		}
	case "roles":
		if v, ok := val.([]string); ok {
			d.Roles_ = DataStringSlice(v)
		}
	default:
		if d.Items_ == nil {
			d.Items_ = make(DataMap)
		}
		d.Items_[key] = val
	}
	return d
}

func (d *DefaultData) Get(key string) any {
	switch key {
	case "id":
		return d.ID_
	case "account":
		return d.Account_
	case "roles":
		return []string(d.Roles_)
	default:
		if d.Items_ != nil {
			if v, ok := d.Items_[key]; ok {
				return v
			}
		}
	}
	return nil
}
