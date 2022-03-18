package serializer

import (
	"github.com/golang/protobuf/ptypes/any"
)

var Encoding = map[string]int64{
	"DEFAULT": 0,
}

type Serializer interface {
	Serialize(interface{}) (*any.Any, error)
	Unserialize(*any.Any) (interface{}, error)
}

type SerializableItem struct {
	EncodingType int64
	Payload      interface{}
	Encoder      Serializer
}

func NewDefaultSerializer(encoder Serializer, payload interface{}) *SerializableItem {
	si := &SerializableItem{
		EncodingType: Encoding["DEFAULT"],
		Payload:      payload,
		Encoder:      encoder,
	}
	return si
}
