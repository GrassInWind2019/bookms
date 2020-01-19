package utils

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
)

func DeepCopy(dst,src interface{}) error {
	if dst == nil {
		return errors.New("dst cannot be null")
	}
	if src == nil {
		return errors.New("src cannot be null")
	}
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, dst)
	return err
}

func DeepCopy2(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
