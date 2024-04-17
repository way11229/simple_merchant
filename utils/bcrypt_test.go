package utils

import "testing"

func TestBcrypt(t *testing.T) {
	pwd := "aaaaAAA1@@1"
	encodePwd, err := BcryptEncryptToHex(pwd)
	if err != nil {
		t.Fatalf("BcryptEncryptToHex error = %v", err)
	}

	if err := BcryptCompareWithHex(encodePwd, pwd); err != nil {
		t.Fatalf("BcryptCompareWithHex error = %v", err)
	}
}
