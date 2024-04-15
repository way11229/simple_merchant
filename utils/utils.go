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

// convert string to uint, panic if convert error
func ConvertStringToUint(str string) uint {
	i, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		panic(fmt.Sprintf("(%s): convert to uint error", str))
	}

	return uint(i)
}
