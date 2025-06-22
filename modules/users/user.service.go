package users

import (
	custom_errors "digishop/utilities/errors"

	"github.com/google/uuid"
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

	user.ID = strUUID.String()
	if user.StrUserType == "Seller" {
		user.UserType = 1
	} else if user.StrUserType == "Buyer" {
		user.UserType = 0
	}
	return u.repo.RegisterUser(user)
}

func factoryUserService(repo iRepo) userService {
	if service == (userService{}) {
		service = userService{
			repo,
		}
	}
	return service
}
