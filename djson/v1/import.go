package djson

import (
	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func MarshalPretty(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", " ")
}

func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
