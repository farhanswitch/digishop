package users

import (
	custom_errors "digishop/utilities/errors"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo iRepo
}

var service userService

func (u userService) RegisterUser(user RegisterUserRequest) (bool, custom_errors.CustomError) {
	strUUID, err := uuid.NewV7()
	if err != nil {
		return true, custom_errors.CustomError{
			Code:          500,
			Message:       err.Error(),
			MessageToSend: "Internal Server Error",
		}
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
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

func (u userService) LoginUser(param LoginUserRequest) (custom_errors.CustomError, LoginUserRequest) {
	switch param.StrUserType {
	case "Seller":
		param.UserType = 1
	case "Buyer":
		param.UserType = 0
	}
	errObj, loginData := u.repo.LoginUser(param)
	if errObj.Code > 10 {
		return errObj, loginData
	}
	err := bcrypt.CompareHashAndPassword([]byte(loginData.Password), []byte(param.Password))
	if err != nil {
		return custom_errors.CustomError{
			Code:          http.StatusBadRequest,
			Message:       err.Error(),
			MessageToSend: "Invalid username or password",
		}, LoginUserRequest{}
	}
	return errObj, loginData
}

func factoryUserService(repo iRepo) userService {
	if service == (userService{}) {
		service = userService{
			repo,
		}
	}
	return service
}
