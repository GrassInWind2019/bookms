package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
)

const (
	saltStr = "bookms"
)

func EncryptPassword(inputPass string) (string, error) {
	hash := md5.New()
	hash.Write([]byte(inputPass))
	salt, _ := salt()
	return hex.EncodeToString(hash.Sum([]byte(salt))), nil
}

func VerifyPassword(storePass string, inputPass string) (bool, error) {
	encryptPass, err := EncryptPassword(inputPass)
	if err != nil {
		return false, err
	}

	if encryptPass != storePass {
		return false, errors.New("密码错误")
	}

	return true, nil
}

func salt() (string, error) {
	hash := md5.New()
	hash.Write([]byte(saltStr))
	return hex.EncodeToString(hash.Sum([]byte(saltStr))), nil
}