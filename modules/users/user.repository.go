package users

import (
	"digishop/connections"
	custom_errors "digishop/utilities/errors"
	"net/http"
)

type userRepo struct{}

var repo userRepo

func (u userRepo) RegisterUser(user RegisterUserRequest) (bool, custom_errors.CustomError) {
	_, err := connections.DbMySQL().Query("CALL create_user(?,?,?,?,?,?,?,?)", user.ID, user.Username, user.FirstName, user.LastName, user.Email, user.Password, user.PhoneNumber, user.UserType)
	if err != nil {
		if err.Error() == "Error 1644 (45000): Email already registered." {
			customErr := custom_errors.CustomError{
				Code:          422,
				MessageToSend: "Email already registered. Please login",
				Message:       err.Error(),
			}
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
func (u userRepo) LoginUser(param LoginUserRequest) (custom_errors.CustomError, LoginUserRequest) {
	result, err := connections.DbMySQL().Query("CALL login(?,?)", param.Username, param.UserType)
	if err != nil {

		return custom_errors.CustomError{
			Code:          http.StatusBadRequest,
			Message:       err.Error(),
			MessageToSend: "Invalid username or password",
		}, LoginUserRequest{}

	}
	if result.Next() {
		var data LoginUserRequest
		err := result.Scan(&data.ID, &data.Username, &data.Password, &data.UserType)
		if err != nil {
			return custom_errors.CustomError{
				Code:          http.StatusBadRequest,
				Message:       err.Error(),
				MessageToSend: "Invalid username or password",
			}, LoginUserRequest{}
		}
		switch data.UserType {
		case 0:
			data.StrUserType = "Buyer"
		case 1:
			data.StrUserType = "Seller"
		}
		return custom_errors.CustomError{}, data
	}
	return custom_errors.CustomError{
		Code:          http.StatusInternalServerError,
		Message:       "Unhandled condition",
		MessageToSend: "Invalid username or password",
	}, LoginUserRequest{}
}
func factoryUserRepo() userRepo {
	if repo == (userRepo{}) {
		repo = userRepo{}
	}
	return repo
}
