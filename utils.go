package fasthttpsession

import (
	"encoding/json"
	"encoding/gob"
	"bytes"
)

// fasthttpsession utils

func NewUtils() *utils {
	return &utils{}
}

type utils struct {

}

// json encode
func (s *utils) JsonEncode(data map[string]interface{}) ([]byte, error) {
	return json.Marshal(data)
}

// json decode
func (s *utils) JsonDecode(data []byte) (map[string]interface{}, error) {
	tempValue := make(map[string]interface{})
	err := json.Unmarshal(data, &tempValue)
	if err != nil {
		return tempValue, err
	}
	return tempValue, nil
}

// gob encode
func (s *utils) GobEncode(data map[string]interface{}) ([]byte, error) {
	if len(data) == 0 {
		return []byte(""), nil
	}
	for _, v := range data {
		gob.Register(v)
	}
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return []byte(""), err
	}
	return buf.Bytes(), nil
}

// gob decode data to map
func (s *utils) GobDecode(data []byte) (map[string]interface{}, error) {

	if len(data) == 0 {
		return make(map[string]interface{}), nil
	}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	var out map[string]interface{}
	err := dec.Decode(&out)
	if err != nil {
		return make(map[string]interface{}), err
	}
	return out, nil
}
