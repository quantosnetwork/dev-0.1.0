package factory

import (
	"encoding/json"
	"reflect"
)

// factory includes factory methods that are common to most factories interface

type Factory interface {
	ConvertTypeToMap(dataType interface{}) map[string]interface{}
	BuildFromBytes(data []byte, dataType interface{}) (interface{}, error)
	CreateEmpty(dataType interface{}) interface{}
}

type factory struct{}

func (f factory) ConvertTypeToMap(dataType interface{}) map[string]interface{} {
	b, _ := json.Marshal(dataType)

	err := json.Unmarshal(b, &dataType)
	if err != nil {
		return nil
	}
	return dataType.(map[string]interface{})

}

func (f factory) BuildFromBytes(data []byte, dataType interface{}) (interface{}, error) {
	_ = json.Unmarshal(data, &dataType)
	return dataType, nil
}

func (f factory) CreateEmpty(dataType interface{}) interface{} {
	typeOf := reflect.TypeOf(dataType)
	t := reflect.New(typeOf)
	return t
}

func GetFactory() Factory {
	return factory{}
}
