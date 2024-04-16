package utils

import (
	"fmt"
	"os"
	"strconv"
)

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
