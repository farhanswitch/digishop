package markets

import custom_errors "digishop/utilities/errors"

type category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type iRepo interface {
	GetAllCategory() ([]category, custom_errors.CustomError)
}
