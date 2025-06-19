package users

import (
	"digishop/connections"
)

type userRepo struct{}

var repo userRepo

func (u userRepo) RegisterUser(user RegisterUserRequest) error {
	_, err := connections.DbMySQL().Query("CALL create_user(?,?,?,?,?,?,?,?)", user.ID, user.Username, user.FirstName, user.LastName, user.Email, user.Password, user.PhoneNumber, user.UserType)
	return err
}
