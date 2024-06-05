package apikey

import (
	"github.com/bytedance/sonic"
)

type Serializer interface {
	Serialize(*Data) ([]byte, error)
	Deserialize([]byte, *Data) error
}

type SonicSerializer struct{}

func (s SonicSerializer) Serialize(v *Data) ([]byte, error) {
	return sonic.Marshal(v)
}

func (s SonicSerializer) Deserialize(d []byte, v *Data) error {
	return sonic.Unmarshal(d, v)
}
