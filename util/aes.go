package util

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/golang/glog"
)

/* ref https://github.com/matt-wu/AES
use aes cfb mode, the key is for aes to encode and decode
*/
var KeyCoder = "qsefthuko;jengo6074@$userzdvhmlp"

var commonIV = []byte{0x1a, 0x0b, 0xc1, 0xd5, 0xec, 0x0f, 0x07, 0xf1, 0x2e, 0xbb, 0x8f, 0x04, 0xe6, 0xee, 0xfe, 0x9d}

/* ref https://github.com/matt-wu/AES
use aes cfb mode, two factors to decode: the key and commonIV
todo: enhance the key generation and storage
*/
func AESEncode(key []byte, plainText []byte) (cipherText []byte, err error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		glog.Error(err)
		return
	}
	cfb := cipher.NewCFBEncrypter(c, commonIV)
	cipherText = make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return
}

func AESDecode(key []byte, cipherText []byte) (plainText []byte, err error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		glog.Error(err)
		return
	}
	cfb := cipher.NewCFBDecrypter(c, commonIV)
	plainText = make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return
}
