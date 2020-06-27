package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
)

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey(priv []byte) (key *rsa.PrivateKey, err error) {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes

	if enc {
		fmt.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return
		}
	}
	key, err = x509.ParsePKCS1PrivateKey(b)
	return
}

// RsaDecrypt data
func RsaDecrypt(ciphertext []byte, instalKey *rsa.PrivateKey, factoryKey *rsa.PrivateKey) (data []byte, err error) {
	if data, err = rsa.DecryptPKCS1v15(rand.Reader, instalKey, ciphertext); err != nil {
		data, err = rsa.DecryptPKCS1v15(rand.Reader, factoryKey, ciphertext)
	}
	return
}

// Aes256Decrypt data
func Aes256Decrypt(data []byte, key string, iv string) (decryptedData []byte, err error) {

	decodedKey, err := hex.DecodeString(key)
	decodedIv, err := hex.DecodeString(iv)

	block, err := aes.NewCipher(decodedKey)
	if err != nil {
		panic(err)
	}

	decryptedData = make([]byte, len(data))
	mode := cipher.NewCBCDecrypter(block, decodedIv)
	mode.CryptBlocks(decryptedData, data)
	return
}
