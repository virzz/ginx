package auth

import (
	"crypto/rand"
	"encoding"
	"encoding/hex"
	"io"

	"github.com/go-json-experiment/json"
)

type (
	IDType interface {
		string | int64 | uint64 | int | uint
	}
	DataStringSlice []string
	DataMap         map[string]any
	Data[T IDType]  interface {
		Token() string
		ID() T
		Account() string
		Roles() []string
		Items() DataMap
		Get(string) any
		New() string
		Set(string, any) Data[T]
		SetToken(string) Data[T]
		SetID(T) Data[T]
		SetAccount(string) Data[T]
		SetValues(string, any) Data[T]
		SetRoles([]string) Data[T]
		Delete(string) Data[T]
		Clear() Data[T]
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

type DefaultData[T IDType] struct {
	Token_   string          `json:"token" redis:"token"`
	ID_      T               `json:"id" redis:"id"`
	Account_ string          `json:"account" redis:"account"`
	Roles_   DataStringSlice `json:"roles" redis:"roles"`
	Items_   DataMap         `json:"items" redis:"items"`
}

var _ Data[int64] = (*DefaultData[int64])(nil)

func New() string {
	k := make([]byte, 20)
	io.ReadFull(rand.Reader, k)
	return hex.EncodeToString(k)
}

func (d *DefaultData[T]) New() string {
	_d := &DefaultData[T]{}
	_d.Token_ = New()
	*d = *_d
	return d.Token_
}
func (d *DefaultData[T]) ID() T           { return d.ID_ }
func (d *DefaultData[T]) Token() string   { return d.Token_ }
func (d *DefaultData[T]) Account() string { return d.Account_ }
func (d *DefaultData[T]) Roles() []string { return []string(d.Roles_) }
func (d *DefaultData[T]) Items() DataMap  { return d.Items_ }
func (d *DefaultData[T]) SetToken(v string) Data[T] {
	d.Token_ = v
	return d
}
func (d *DefaultData[T]) SetID(v T) Data[T] {
	d.ID_ = v
	return d
}
func (d *DefaultData[T]) SetAccount(v string) Data[T] {
	d.Account_ = v
	return d
}

func (d *DefaultData[T]) SetRoles(v []string) Data[T] {
	d.Roles_ = DataStringSlice(v)
	return d
}

func (d *DefaultData[T]) SetValues(k string, v any) Data[T] {
	if d.Items_ == nil {
		d.Items_ = make(DataMap)
	}
	d.Items_[k] = v
	return d
}

func (d *DefaultData[T]) Set(key string, val any) Data[T] {
	switch key {
	case "id":
		if v, ok := val.(T); ok {
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

func (d *DefaultData[T]) Get(key string) any {
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

func (d *DefaultData[T]) Delete(key string) Data[T] {
	switch key {
	case "id":
		_id := new(T)
		d.ID_ = *_id
	case "account":
		d.Account_ = ""
	case "roles":
		d.Roles_ = nil
	default:
		if d.Items_ != nil {
			delete(d.Items_, key)
		}
	}
	return d
}

func (d *DefaultData[T]) Clear() Data[T] {
	_d := new(DefaultData[T])
	*d = *_d
	return d
}
