package stores

import custom_errors "digishop/utilities/errors"

type storeData struct {
	ID      string `json:"id"`
	Name    string `json:"name" validate:"required,min=6,max=100"`
	Address string `json:"address" validate:"required,min=6,max=255"`
	UserID  string `json:"userID"`
}

type iRepo interface {
	RegisterStore(store storeData) (bool, custom_errors.CustomError)
	GetStoreByUserID(id string) (storeData, custom_errors.CustomError)
	UpdateStore(store storeData) (bool, custom_errors.CustomError)
}
