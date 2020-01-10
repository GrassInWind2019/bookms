package utils

import (
	"bytes"
	"encoding/gob"
)

func Decode(value string, r interface{}) error {
	buff := bytes.NewBuffer([]byte(value))
	dec := gob.NewDecoder(buff)
	return dec.Decode(r)
}

func Encode(value interface{}) (string, error) {
	buff := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buff)
	err := enc.Encode(value)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}
