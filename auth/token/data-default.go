package token

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

type DefaultData struct {
	Token_   string          `json:"token" redis:"token"`
	ID_      string          `json:"id" redis:"id"`
	Account_ string          `json:"account" redis:"account"`
	Roles_   DataStringSlice `json:"roles" redis:"roles"`
	Items_   DataMap         `json:"items" redis:"items"`
}

var _ Data = (*DefaultData)(nil)

func New() string {
	k := make([]byte, 20)
	io.ReadFull(rand.Reader, k)
	return hex.EncodeToString(k)
}

func (d *DefaultData) New() string {
	_d := &DefaultData{}
	_d.Token_ = New()
	*d = *_d
	return d.Token_
}
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

func (d *DefaultData) SetValues(k string, v any) Data {
	if d.Items_ == nil {
		d.Items_ = make(DataMap)
	}
	d.Items_[k] = v
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

func (d *DefaultData) Delete(key string) Data {
	switch key {
	case "id":
		d.ID_ = ""
	case "account":
		d.Account_ = ""
	case "roles":
		d.Roles_ = nil
	default:
		if d.Items_ != nil {
			if _, ok := d.Items_[key]; ok {
				delete(d.Items_, key)
			}
		}
	}
	return d
}
