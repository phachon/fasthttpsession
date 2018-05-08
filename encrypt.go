package fasthttpsession

import (
	"encoding/json"
	"encoding/gob"
	"bytes"
	"encoding/base64"
)

// fasthttpsession encrypt tool
// - json
// - gob
// - base64

const (
	BASE64TABLE = "1234567890poiuytreqwasdfghjklmnbvcxzQWERTYUIOPLKJHGFDSAZXCVBNM-_"
)

func NewEncrypt() *encrypt {
	return &encrypt{}
}

type encrypt struct {

}

// json encode
func (s *encrypt) JsonEncode(data map[string]interface{}) ([]byte, error) {
	return json.Marshal(data)
}

// json decode
func (s *encrypt) JsonDecode(data []byte) (map[string]interface{}, error) {
	tempValue := make(map[string]interface{})
	err := json.Unmarshal(data, &tempValue)
	if err != nil {
		return tempValue, err
	}
	return tempValue, nil
}

// gob encode
func (s *encrypt) GobEncode(data map[string]interface{}) ([]byte, error) {
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
func (s *encrypt) GobDecode(data []byte) (map[string]interface{}, error) {

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

// base64 encode
func (s *encrypt) Base64Encode(data map[string]interface{}) ([]byte, error) {
	var coder = base64.NewEncoding(BASE64TABLE)
	b, err := s.GobEncode(data)
	if err != nil {
		return []byte{}, err
	}
	return []byte(coder.EncodeToString(b)), nil
}

// base64 decode
func (s *encrypt) Base64Decode(data []byte) (map[string]interface{}, error) {
	var coder = base64.NewEncoding(BASE64TABLE)
	b, err := coder.DecodeString(string(data))
	if err != nil {
		return nil, err
	}
	return s.GobDecode(b)
}