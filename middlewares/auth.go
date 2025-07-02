package middlewares

import (
	"digishop/configs"
	"digishop/utilities"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"gopkg.in/square/go-jose.v2/jwt"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Get header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"errors":"Unauthorized"}`))
			return
		}
		var token string
		if strings.HasPrefix(authHeader, "Bearer") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			log.Println("Invalid token")
		}
		decryptedClaims, err := utilities.JWEDecryptAES(token, []byte(configs.GetConfig().Service.EncryptKey))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"errors":"Unauthorized"}`))
			return
		}
		username, ok := decryptedClaims["username"].(string)
		if !ok {
			log.Println("Invalid token. There is no username in token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"errors":"Unauthorized"}`))
		}
		token1, err := utilities.RedisInstance().GetValue(fmt.Sprintf("TOKEN_%s_1", username))
		if err != nil {
			if err.Error() == "redis: nil" {
				token1 = ""
			} else {
				log.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"errors":"Unauthorized"}`))
				return
			}
		}
		token2, err := utilities.RedisInstance().GetValue(fmt.Sprintf("TOKEN_%s_2", username))
		if err != nil {
			if err.Error() == "redis: nil" {
				token2 = ""
			} else {
				log.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"errors":"Unauthorized"}`))
				return
			}
		}
		if token1 != token && token2 != token {
			log.Println("Token missmatch with token saved in Redis")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"errors":"Unauthorized"}`))
			return
		}
		floatExpireTime := decryptedClaims["Expiry"].(float64)
		intExpireTime := int64(floatExpireTime)
		formattedExpireTime := time.Unix(intExpireTime, 0)
		currentTime := time.Now()
		remainingTime := formattedExpireTime.Sub(currentTime).Seconds()
		// Check if the remaining time is less than or equal to refresh session and the token2 is not set
		if uint16(remainingTime) <= configs.GetConfig().Service.RefreshTime && token2 == "" {
			claims := map[string]interface{}{
				"username": decryptedClaims["username"],
				"id":       decryptedClaims["id"],
				"Issuer":   "Digishop",
				"Expiry":   jwt.NewNumericDate(time.Now().Add(time.Duration(configs.GetConfig().Service.SessionTime) * time.Second)),
				"IssuedAt": jwt.NewNumericDate(time.Now()),
			}
			strToken, err := utilities.JWEEncryptAES(claims, []byte(configs.GetConfig().Service.EncryptKey))
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"errors":"Unauthorized"}`))
				return
			}
			err = utilities.RedisInstance().SaveValue(fmt.Sprintf("TOKEN_%s_2", username), token1, time.Duration(configs.GetConfig().Service.RefreshTime)*time.Second)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"errors":"Unauthorized"}`))
				return

			}
			err = utilities.RedisInstance().SaveValue(fmt.Sprintf("TOKEN_%s_1", username), strToken, time.Duration(configs.GetConfig().Service.SessionTime)*time.Second)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"errors":"Unauthorized"}`))
				return
			}
			w.Header().Set("XRF-TOKEN", strToken)
		}
		w.Header().Set("XRF-TOKEN", token1)
		strDecryptedClaims, _ := json.Marshal(decryptedClaims)
		r.Header.Set("X-USER-DATA", string(strDecryptedClaims))
		next.ServeHTTP(w, r)
	})
}
