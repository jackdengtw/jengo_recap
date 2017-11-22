package util

import (
	"testing"
)

func TestAESDecode(t *testing.T) {
	plain := "hereicometotest"
	cip, err := AESEncode([]byte(KeyCoder), []byte(plain))
	if err != nil {
		t.Fatal("err when do aes encode because of :", err)
	}
	t.Logf("\n%s => encode %s \n", plain, string(cip))

	after, err := AESDecode([]byte(KeyCoder), cip)
	if err != nil {
		t.Fatal("err when do aes decode because of :", err)
	}
	t.Logf("\n%s => decode %s \n", string(cip), string(after))
	if plain != string(after) {
		t.Errorf("%s is not equal to %s after encode and decode", plain, after)
	}

}
