package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"strconv"
)

const NUMBER_CHARS = "1234567890"

// return $KEY from environment variable, panic if $KEY not found
func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("(%s): missing env", key))
	}

	return value
}

func ConvertStringToUint(str string) (uint, error) {
	i, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(i), nil
}

func ConvertInt64ToUint32(num int64) (uint32, error) {
	rtn := uint32(num)
	if int64(rtn) != num {
		return 0, errors.New("overflow")
	}

	return rtn, nil
}

func GenerateRandomCode(length uint) (string, error) {
	buf := make([]byte, length)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	codeCharsLength := len(NUMBER_CHARS)
	for i := uint(0); i < length; i++ {
		buf[i] = NUMBER_CHARS[int(buf[i])%codeCharsLength]
	}

	return string(buf), nil
}
