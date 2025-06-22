package users

import custom_errors "digishop/utilities/errors"

type RegisterUserRequest struct {
	UserType    uint8  `json:"-"`
	StrUserType string `json:"userType" validate:"required,oneof='Seller' 'Buyer'"`
	Username    string `json:"username" validate:"required,min=6,max=15"`
	FirstName   string `json:"firstName" validate:"required,min=3,max=50"`
	LastName    string `json:"lastName" validate:"min=0,max=50"`
	Password    string `json:"password" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phoneNumber" validate:"required,min=10,max=15"`
	ID          string `json:"id"`
}

type LoginUserRequest struct {
	UserType    uint8  `json:"-"`
	StrUserType string `json:"userType" validate:"required,oneof='Seller' 'Buyer'"`
	Username    string `json:"username" validate:"required,min=6,max=15"`
	Password    string `json:"password" validate:"required"`
	ID          string `json:"id"`
}

type iRepo interface {
	RegisterUser(user RegisterUserRequest) (bool, custom_errors.CustomError)
	LoginUser(param LoginUserRequest) (custom_errors.CustomError, LoginUserRequest)
}
