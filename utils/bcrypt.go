package utils

import (
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func BcryptEncryptToHex(data string) (string, error) {
	encrypt, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(encrypt), nil
}

func BcryptCompareWithHex(from, matchTo string) error {
	fromByte, err := hex.DecodeString(from)
	if err != nil {
		return err
	}

	return bcrypt.CompareHashAndPassword(fromByte, []byte(matchTo))
}
