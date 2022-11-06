package encrypdecryp

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type Code struct {
	MySecret string
	Bytes    []byte
}

func (c *Code) Encrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(c.MySecret))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, c.Bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (c *Code) Decrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(c.MySecret))
	if err != nil {
		return "", err
	}
	cipherText, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}
	cfb := cipher.NewCFBDecrypter(block, c.Bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}
