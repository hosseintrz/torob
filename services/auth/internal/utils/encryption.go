package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

var iv = []byte{65, 19, 37, 59, 90, 71, 42, 23, 78, 12, 06, 76, 55, 17, 47, 89}

func encodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Encrypt(key, text string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	plaintext := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	return encodeBase64(ciphertext), nil
}
