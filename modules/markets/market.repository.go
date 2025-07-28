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
func (m marketRepo) GetListProductByCategory(categoryID string) ([]productData, custom_errors.CustomError) {
	var products []productData
	results, err := connections.DbMySQL().Query("SELECT p.id, p.name, p.price, s.name, f.filename FROM products p  JOIN stores s ON p.store_id = s.id JOIN files f ON p.image_id = f.id WHERE p.category_id = ?", categoryID)
	if err != nil {
		return []productData{}, custom_errors.CustomError{
			Code:          http.StatusInternalServerError,
			MessageToSend: "Internal Server Error",
			Message:       err.Error(),
		}
	}
	for results.Next() {
		var data productData
		err := results.Scan(&data.ID, &data.Name, &data.Price, &data.StoreName, &data.ImagePath)
		if err != nil {
			return []productData{}, custom_errors.CustomError{
				Code:          http.StatusInternalServerError,
				MessageToSend: "Internal Server Error",
				Message:       err.Error(),
			}
		}
		products = append(products, data)
	}
	return products, custom_errors.CustomError{}
}
func factoryMarketRepository() marketRepo {
	if repo == (marketRepo{}) {
		repo = marketRepo{}
	}
	return repo
}
