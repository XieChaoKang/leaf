package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// AES CBC zeroPadding base64
func AESEncryptCBC(plaintext, k, iv string) (string, error) {
	data := []byte(plaintext)
	key := []byte(k)
	aesBlockEncrypter, err := aes.NewCipher(key)
	content := ZeroPadding(data, aesBlockEncrypter.BlockSize())
	encrypted := make([]byte, len(content))
	if err != nil {
		println(err.Error())
		return "", err
	}
	aesEncrypter := cipher.NewCBCEncrypter(aesBlockEncrypter, []byte(iv))
	aesEncrypter.CryptBlocks(encrypted, content)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// AES CBC zeroPadding base64
func AESDecryptCBC(chiper, k, iv string) (string, error) {
	arr, _ := base64.StdEncoding.DecodeString(chiper)
	key := []byte(k)
	decrypted := make([]byte, len(arr))
	aesBlockDecrypter, err := aes.NewCipher(key)
	if err != nil {
		println(err.Error())
		return "", err
	}
	aesDecrypter := cipher.NewCBCDecrypter(aesBlockDecrypter, []byte(iv))
	aesDecrypter.CryptBlocks(decrypted, arr)
	trimming := ZeroTrimming(decrypted)
	return string(trimming), nil
}

func ZeroPadding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(0)}, padding)
	return append(cipherText, padText...)
}

func ZeroTrimming(encrypt []byte) []byte {
	return bytes.TrimRight(encrypt, string([]byte{0}))
}
