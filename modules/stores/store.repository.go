package stores

import (
	"digishop/connections"
	custom_errors "digishop/utilities/errors"
	"net/http"
)

type storeRepo struct{}

var repo storeRepo

func (s storeRepo) RegisterStore(store storeData) (bool, custom_errors.CustomError) {
	_, err := connections.DbMySQL().Exec("INSERT INTO stores(id, name, address, user_id) VALUES(?,?,?,?)", store.ID, store.Name, store.Address, store.UserID)
	if err != nil {
		customErr := custom_errors.CustomError{
			Code:          http.StatusInternalServerError,
			MessageToSend: "Internal Server Error",
			Message:       err.Error(),
		}
		return true, customErr
	}
	return false, custom_errors.CustomError{}
}
func (s storeRepo) UpdateStore(store storeData) (bool, custom_errors.CustomError) {
	_, err := connections.DbMySQL().Exec("UPDATE stores SET name = ?, address = ? WHERE user_id = ?", store.Name, store.Address, store.UserID)
	if err != nil {
		customErr := custom_errors.CustomError{
			Code:          http.StatusInternalServerError,
			MessageToSend: "Internal Server Error",
			Message:       err.Error(),
		}
		return true, customErr
	}
	return false, custom_errors.CustomError{}
}
func (s storeRepo) GetStoreByUserID(userID string) (storeData, custom_errors.CustomError) {
	var store storeData
	err := connections.DbMySQL().QueryRow("SELECT id, name, address, user_id FROM stores WHERE user_id = ?", userID).Scan(&store.ID, &store.Name, &store.Address, &store.UserID)
	if err != nil && err.Error() != "sql: no rows in result set" {
		customErr := custom_errors.CustomError{
			Code:          http.StatusInternalServerError,
			MessageToSend: "Internal Server Error",
			Message:       err.Error(),
		}
		return storeData{}, customErr
	}
	return store, custom_errors.CustomError{}
}
func factoryStoreRepo() iRepo {
	if repo == (storeRepo{}) {
		repo = storeRepo{}
	}
	return repo
}
