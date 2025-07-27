package markets

import (
	"digishop/connections"
	custom_errors "digishop/utilities/errors"
	"net/http"
)

type marketRepo struct {
}

var repo marketRepo

func (m marketRepo) GetAllCategory() ([]category, custom_errors.CustomError) {
	var categories []category
	results, err := connections.DbMySQL().Query("SELECT id, name FROM categories")
	if err != nil {
		return []category{}, custom_errors.CustomError{
			Code:          http.StatusInternalServerError,
			MessageToSend: "Internal Server Error",
			Message:       err.Error(),
		}

	}
	for results.Next() {
		var data category
		err := results.Scan(&data.ID, &data.Name)
		if err != nil {
			return []category{}, custom_errors.CustomError{
				Code:          http.StatusInternalServerError,
				MessageToSend: "Internal Server Error",
				Message:       err.Error(),
			}
		}
		categories = append(categories, data)
	}
	return categories, custom_errors.CustomError{}
}

func factoryMarketRepository() marketRepo {
	if repo == (marketRepo{}) {
		repo = marketRepo{}
	}
	return repo
}
