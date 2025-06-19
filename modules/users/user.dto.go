package users

type RegisterUserRequest struct {
	UserType    uint8  `json:"userType" validate:"required"`
	Username    string `json:"username" validate:"required, min=6, max=15"`
	FirstName   string `json:"firstName" validate:"required, min=3, max=50"`
	LastName    string `json:"lastName" validate:"min=0, max=50"`
	Password    string `json:"password" validate:"required"`
	Email       string `json:"email" validate:"required, email"`
	PhoneNumber string `json:"phoneNumber" validate:"required, min=10, max=15"`
	ID          string `json:"id"`
}
