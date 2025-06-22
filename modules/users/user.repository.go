package users

import (
	"digishop/connections"
	custom_errors "digishop/utilities/errors"
)

type userRepo struct{}

var repo userRepo

func (u userRepo) RegisterUser(user RegisterUserRequest) (bool, custom_errors.CustomError) {
	_, err := connections.DbMySQL().Query("CALL create_user(?,?,?,?,?,?,?,?)", user.ID, user.Username, user.FirstName, user.LastName, user.Email, user.Password, user.PhoneNumber, user.UserType)
	if err != nil {
		if err.Error() == "Error 45000: Email already registered." {
			customErr := custom_errors.CustomError{
				Code:    422,
				Message: err.Error(),
			}
			customErr.Compile()
			return true, customErr
		} else {
			customErr := custom_errors.CustomError{
				Code:          500,
				Message:       err.Error(),
				MessageToSend: "Internal server error.",
			}
			return true, customErr
		}
	}
	return false, custom_errors.CustomError{}
}

func factoryUserRepo() userRepo {
	if repo == (userRepo{}) {
		repo = userRepo{}
	}
	return repo
}
