package utilities

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"

	"digishop/configs"
)

func EncryptRSA(message []byte) (string, error) {
	rsaPublicKey := configs.GetConfig().Service.PublicKeyRSA

	encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPublicKey, message, nil)
	if err != nil {
		return "", err
	}
	// Convert to base64 so we can transfer the data easily
	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}
func DecryptRSA(encryptedMessage string) ([]byte, error) {
	rsaPrivateKey := configs.GetConfig().Service.PrivateKeyRSA

	// Decode the base64 string
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedMessage)
	if err != nil {
		return nil, err
	}
	// Decrypt the message
	decryptedBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaPrivateKey, encryptedBytes, nil)
	if err != nil {
		return nil, err
	}
	return decryptedBytes, nil
}
