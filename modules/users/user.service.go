package users

import (
	"digishop/configs"
	"digishop/utilities"
	custom_errors "digishop/utilities/errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/square/go-jose.v2/jwt"
)

type userService struct {
	repo iRepo
}

var service userService

func (u userService) RegisterUser(user RegisterUserRequest) (bool, custom_errors.CustomError) {
	strUUID, err := uuid.NewV7()
	if err != nil {
		log.Println(err)
		return true, custom_errors.CustomError{
			Code:          500,
			Message:       err.Error(),
			MessageToSend: "Internal Server Error",
		}
	}
	decryptedPassword, err := utilities.DecryptRSA(user.Password)
	if err != nil {
		log.Println(err)
		return true, custom_errors.CustomError{
			Code:          500,
			Message:       err.Error(),
			MessageToSend: "Internal Server Error",
		}
	}
	decryptedConfirmPassword, err := utilities.DecryptRSA(user.ConfirmPassword)
	if err != nil {
		log.Println(err)
		return true, custom_errors.CustomError{
			Code:          500,
			Message:       err.Error(),
			MessageToSend: "Internal Server Error",
		}
	}
	if string(decryptedConfirmPassword) != string(decryptedPassword) {

		return true, custom_errors.CustomError{
			Code:          400,
			Message:       "Password and Confirm Password must match!",
			MessageToSend: "Password and Confirm Password must match!",
		}
	}
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(decryptedPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return true, custom_errors.CustomError{
			Code:          500,
			Message:       err.Error(),
			MessageToSend: "Internal Server Error",
		}
	}
	user.Password = string(hashedPassword)
	user.ID = strUUID.String()
	switch user.StrUserType {
	case "Seller":
		user.UserType = 1
	case "Buyer":
		user.UserType = 0
	}
	return u.repo.RegisterUser(user)
}

func (u userService) LoginUser(param LoginUserRequest) (custom_errors.CustomError, map[string]any) {
	switch param.StrUserType {
	case "Seller":
		param.UserType = 1
	case "Buyer":
		param.UserType = 0
	}
	errObj, loginData := u.repo.LoginUser(param)
	if errObj.Code > 10 {
		return errObj, map[string]any{}
	}

	decryptedPassword, err := utilities.DecryptRSA(param.Password)
	if err != nil {
		log.Println(err)
		return custom_errors.CustomError{
			Code:          http.StatusBadRequest,
			Message:       err.Error(),
			MessageToSend: "Invalid username or password",
		}, map[string]any{}
	}
	err = bcrypt.CompareHashAndPassword([]byte(loginData.Password), decryptedPassword)
	if err != nil {
		log.Println(err, "999")
		return custom_errors.CustomError{
			Code:          http.StatusBadRequest,
			Message:       err.Error(),
			MessageToSend: "Invalid username or password",
		}, map[string]any{}
	}
	claims := map[string]interface{}{
		"username": loginData.Username,
		"id":       loginData.ID,
		"Issuer":   "Digishop",
		"Expiry":   jwt.NewNumericDate(time.Now().Add(time.Duration(configs.GetConfig().Service.SessionTime) * time.Second)),
		"IssuedAt": jwt.NewNumericDate(time.Now()),
	}
	strToken, err := utilities.JWEEncryptAES(claims, []byte(configs.GetConfig().Service.EncryptKey))
	if err != nil {
		log.Println(err)
		return custom_errors.CustomError{
			Code:          http.StatusBadRequest,
			Message:       err.Error(),
			MessageToSend: "Invalid username or password",
		}, map[string]any{}
	}
	err = utilities.RedisInstance().DeleteValue(fmt.Sprintf("TOKEN_%s_1", loginData.Username))
	if err != nil {
		log.Println(err)
		return custom_errors.CustomError{
			Code:          http.StatusBadRequest,
			Message:       err.Error(),
			MessageToSend: "Invalid username or password",
		}, map[string]any{}
	}
	err = utilities.RedisInstance().DeleteValue(fmt.Sprintf("TOKEN_%s_2", loginData.Username))
	if err != nil {
		log.Println(err)
		return custom_errors.CustomError{
			Code:          http.StatusBadRequest,
			Message:       err.Error(),
			MessageToSend: "Invalid username or password",
		}, map[string]any{}
	}
	err = utilities.RedisInstance().SaveValue(fmt.Sprintf("TOKEN_%s_1", loginData.Username), strToken, time.Duration(configs.GetConfig().Service.SessionTime)*time.Second)

	claims["token"] = strToken
	return errObj, claims
}
func (u userService) CheckAuthentication(token string) (string, custom_errors.CustomError) {
	decryptedClaims, err := utilities.JWEDecryptAES(token, []byte(configs.GetConfig().Service.EncryptKey))
	if err != nil {
		log.Println(err)
		return "", custom_errors.CustomError{
			Code:          http.StatusUnauthorized,
			Message:       err.Error(),
			MessageToSend: "Unauthencticated",
		}
	}
	username, ok := decryptedClaims["username"].(string)
	if !ok {
		log.Println(err)
		return "", custom_errors.CustomError{
			Code:          http.StatusUnauthorized,
			Message:       "Token has no property username",
			MessageToSend: "Unauthencticated",
		}
	}
	token1, err := utilities.RedisInstance().GetValue(fmt.Sprintf("TOKEN_%s_1", username))
	if err != nil {
		if err.Error() == "redis: nil" {
			token1 = ""
		} else {
			log.Println(err)

			return "", custom_errors.CustomError{
				Code:          http.StatusUnauthorized,
				Message:       err.Error(),
				MessageToSend: "Unauthencticated",
			}
		}
	}
	token2, err := utilities.RedisInstance().GetValue(fmt.Sprintf("TOKEN_%s_2", username))
	if err != nil {
		if err.Error() == "redis: nil" {
			token2 = ""
		} else {
			log.Println(err)

			return "", custom_errors.CustomError{
				Code:          http.StatusUnauthorized,
				Message:       err.Error(),
				MessageToSend: "Unauthencticated",
			}
		}
	}
	if token1 != token && token2 != token {
		return "", custom_errors.CustomError{
			Code:          http.StatusUnauthorized,
			Message:       "Token missmatch with token saved in Redis",
			MessageToSend: "Unauthencticated",
		}
	}
	floatExpireTime := decryptedClaims["Expiry"].(float64)
	intExpireTime := int64(floatExpireTime)

	formattedExpiryTime := time.Unix(intExpireTime, 0)
	currentTime := time.Now()
	remainingTime := formattedExpiryTime.Sub(currentTime).Seconds()
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
			log.Println("Failed to refresh token")
			log.Println(err)
			return token1, custom_errors.CustomError{}
		}
		err = utilities.RedisInstance().SaveValue(fmt.Sprintf("TOKEN_%s_2", username), token1, time.Duration(configs.GetConfig().Service.RefreshTime)*time.Second)
		if err != nil {
			log.Println(err)
			return "", custom_errors.CustomError{
				Code:          http.StatusUnauthorized,
				Message:       err.Error(),
				MessageToSend: "Unauthencticated",
			}
		}
		err = utilities.RedisInstance().SaveValue(fmt.Sprintf("TOKEN_%s_1", username), strToken, time.Duration(configs.GetConfig().Service.SessionTime)*time.Second)
		if err != nil {
			log.Println(err)
			return "", custom_errors.CustomError{
				Code:          http.StatusUnauthorized,
				Message:       err.Error(),
				MessageToSend: "Unauthencticated",
			}
		}
		token1 = strToken
	}
	return token1, custom_errors.CustomError{}
}

func factoryUserService(repo iRepo) userService {
	if service == (userService{}) {
		service = userService{
			repo,
		}
	}
	return service
}
