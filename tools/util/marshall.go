package util

import "encoding/json"

func MarshalMessage(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func UnmarshalMessage(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
