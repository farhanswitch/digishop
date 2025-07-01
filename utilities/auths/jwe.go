package auths

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

func ReadPublicKeyFromFile(filename string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}
	return rsaPub, nil
}
func ReadPrivateKeyFromFile(filename string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
func ParsePublicKey(strPublicKey string) (*rsa.PublicKey, error) {

	block, _ := pem.Decode([]byte(strPublicKey))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}
	return rsaPub, nil
}

func ParsePrivateKey(strPrivateKey string) (*rsa.PrivateKey, error) {

	block, _ := pem.Decode([]byte(strPrivateKey))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func WriteKeysToFile(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, filePathPublic string, filePathPrivate string) error {
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err := os.WriteFile(filePathPrivate, privPEM, 0600); err != nil {
		return err
	}

	pubASN1, _ := x509.MarshalPKIXPublicKey(publicKey)
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})
	if err := os.WriteFile(filePathPublic, pubPEM, 0644); err != nil {
		return err
	}
	return nil
}

func Encrypt(claims map[string]interface{}, publicKey *rsa.PublicKey) (string, error) {
	// Create encrypter
	encrypter, err := jose.NewEncrypter(
		jose.A256GCM,
		jose.Recipient{
			Algorithm: jose.RSA_OAEP_256,
			Key:       publicKey,
		},
		(&jose.EncrypterOptions{}).WithContentType("JWT"),
	)
	if err != nil {
		return "", err
	}
	// Encrypt JWT
	jweToken, err := jwt.Encrypted(encrypter).Claims(claims).CompactSerialize()
	if err != nil {
		panic(err)
	}
	return jweToken, nil
}
func Decrypt(jweToken string, privateKey *rsa.PrivateKey) (map[string]interface{}, error) {
	// Dekripsi
	parsed, err := jwt.ParseEncrypted(jweToken)
	if err != nil {
		return map[string]interface{}{}, err
	}

	var decryptedClaims map[string]interface{}
	err = parsed.Claims(privateKey, &decryptedClaims)
	if err != nil {
		return map[string]interface{}{}, err
	}
	return decryptedClaims, nil

}

// EncryptAES mengenkripsi claims menjadi JWE token menggunakan enkripsi simetris AES-256.
// 'secretKey' harus memiliki panjang 32 byte.
func EncryptAES(claims map[string]interface{}, secretKey []byte) (string, error) {
	// Validasi panjang kunci untuk AES-256
	if len(secretKey) != 32 {
		return "", errors.New("secret key must be 32 bytes long for AES-256")
	}

	// Membuat encrypter dengan algoritma AES-GCM Key Wrap dan A256GCM untuk konten
	encrypter, err := jose.NewEncrypter(
		jose.A256GCM, // Content encryption
		jose.Recipient{
			Algorithm: jose.A256GCMKW, // Key encryption
			Key:       secretKey,
		},
		(&jose.EncrypterOptions{}).WithContentType("JWT"),
	)
	if err != nil {
		return "", err
	}

	// Membuat dan mengenkripsi JWT dengan claims yang diberikan
	jweToken, err := jwt.Encrypted(encrypter).Claims(claims).CompactSerialize()
	if err != nil {
		return "", err
	}

	return jweToken, nil
}

// DecryptAES mendekripsi JWE token yang dienkripsi dengan AES-256.
// 'secretKey' harus sama dengan yang digunakan untuk enkripsi dan memiliki panjang 32 byte.
func DecryptAES(jweToken string, secretKey []byte) (map[string]interface{}, error) {
	// Validasi panjang kunci untuk AES-256
	if len(secretKey) != 32 {
		return nil, errors.New("secret key must be 32 bytes long for AES-256")
	}

	// Mengurai token JWE yang terenkripsi
	parsed, err := jwt.ParseEncrypted(jweToken)
	if err != nil {
		return nil, err
	}

	// Mendekripsi claims menggunakan secret key
	var decryptedClaims map[string]interface{}
	err = parsed.Claims(secretKey, &decryptedClaims)
	if err != nil {
		return nil, err
	}

	return decryptedClaims, nil
}
