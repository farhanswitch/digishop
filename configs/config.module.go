package configs

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type dbConfig struct {
	Host     string
	Port     int
	User     string
	Pass     string
	Database string
}
type redisConfig struct {
	Host string
	Port uint16
}
type serviceConfig struct {
	Port          uint16
	SessionTime   uint16
	RefreshTime   uint16
	EncryptKey    string
	PublicKeyRSA  *rsa.PublicKey
	PrivateKeyRSA *rsa.PrivateKey
}
type appConfig struct {
	Db      dbConfig
	Redis   redisConfig
	Service serviceConfig
}

func InitModule(path string) {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatalf("Cannot load environment file: %s\n", err.Error())
	}
}

var appConf appConfig

func GetConfig() appConfig {
	if appConf == (appConfig{}) {
		dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
		if err != nil {
			log.Fatalf("Invalid Database Port: %s\n", err.Error())
		}
		var dbConf dbConfig = dbConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     dbPort,
			User:     os.Getenv("DB_USER"),
			Pass:     os.Getenv("DB_PASS"),
			Database: os.Getenv("DB_NAME"),
		}
		appPort, err := strconv.Atoi(os.Getenv("APP_PORT"))
		if err != nil {
			log.Fatalf("Invalid Application Port: %s\n", err.Error())
		}
		sessionTime, err := strconv.Atoi(os.Getenv("SESSION_TIME"))
		if err != nil {
			log.Fatalf("Invalid Session Time: %s\n", err.Error())
		}
		refreshTime, err := strconv.Atoi(os.Getenv("REFRESH_TIME"))
		if err != nil {
			log.Fatalf("Invalid Refresh Time: %s\n", err.Error())
		}
		redisPort, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
		if err != nil {
			log.Fatalf("Invalid Redis Port: %s\n", err.Error())
		}
		var redisConf redisConfig = redisConfig{
			Host: os.Getenv("REDIS_HOST"),
			Port: uint16(redisPort),
		}
		strPublicKey := os.Getenv("RSA_PUBLIC")
		strPrivateKey := os.Getenv("RSA_PRIVATE")
		block, _ := pem.Decode([]byte(strPublicKey))
		if block == nil || block.Type != "PUBLIC KEY" {
			log.Fatalf("failed to decode PEM block containing public key. Error: %s\n", err.Error())
		}
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			log.Fatalf("failed to decode PEM block containing public key. Error: %s\n", err.Error())
		}
		rsaPub, ok := pub.(*rsa.PublicKey)
		if !ok {
			log.Fatalf("not an RSA public key.")
		}

		block, _ = pem.Decode([]byte(strPrivateKey))
		if block == nil || block.Type != "RSA PRIVATE KEY" {
			log.Fatalf("failed to decode PEM block containing private key. Error: %s\n", err.Error())
		}
		privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			log.Fatalf("failed to decode PEM block containing private key. Error: %s\n", err.Error())
		}
		var serviceConf serviceConfig = serviceConfig{
			Port:          uint16(appPort),
			SessionTime:   uint16(sessionTime),
			RefreshTime:   uint16(refreshTime),
			EncryptKey:    os.Getenv("ENCRYPT_KEY"),
			PublicKeyRSA:  rsaPub,
			PrivateKeyRSA: privateKey,
		}

		appConf = appConfig{
			Db:      dbConf,
			Redis:   redisConf,
			Service: serviceConf,
		}
		log.Println("Configuration loaded")
	}
	return appConf
}
