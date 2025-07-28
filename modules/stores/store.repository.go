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
func (s storeRepo) GetCategoryNameByID(id string) (string, error) {
	var category string
	err := connections.DbMySQL().QueryRow("SELECT name FROM categories WHERE id = ? LIMIT 1", id).Scan(&category)
	return category, err
}

func (s storeRepo) GetProductImagePathByID(id string) (string, error) {
	var path string
	err := connections.DbMySQL().QueryRow("SELECT path FROM files WHERE id = ? LIMIT 1", id).Scan(&path)
	return path, err
}
func (s storeRepo) CreateNewProduct(product productRequest) (bool, custom_errors.CustomError) {
	_, err := connections.DbMySQL().Exec("INSERT INTO products(id, category_id, store_id, image_id, name, description, price, amount) VALUES(?,?,?,?,?,?,?,?)", product.ID, product.CategoryID, product.StoreID, product.ImageID, product.Name, product.Description, product.Price, product.Amount)
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
func (s storeRepo) CheckIsValidUserProduct(userID string, productID string) error {
	var id int
	err := connections.DbMySQL().QueryRow("SELECT 1 FROM products JOIN stores ON products.store_id = stores.id WHERE stores.user_id = ? AND products.id = ? LIMIT 1", userID, productID).Scan(&id)
	return err
}
func (s storeRepo) UpdateProducts(product updateProductRequest) (bool, custom_errors.CustomError) {
	_, err := connections.DbMySQL().Exec("UPDATE products SET category_id = ?, store_id = ?, image_id = ?, name = ?, description = ?, price = ?, amount = ? WHERE id = ?", product.CategoryID, product.StoreID, product.ImageID, product.Name, product.Description, product.Price, product.Amount, product.ID)
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
func (s storeRepo) GetListProduct(param getListProductRequest) ([]getListProductResponse, custom_errors.CustomError) {
	var products []getListProductResponse
	results, err := connections.DbMySQL().Query("CALL get_list_products(?, ?, ?, ?, ?, ?)", param.StoreID, param.PaginationRows, param.Offset, param.SortField, param.SortOrder, param.Search)
	if err != nil {
		customErr := custom_errors.CustomError{
			Code:          http.StatusInternalServerError,
			MessageToSend: "Internal Server Error",
			Message:       err.Error(),
		}
		return []getListProductResponse{}, customErr
	}
	for results.Next() {
		var product getListProductResponse
		err = results.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Amount, &product.CategoryName, &product.ImagePath, &product.StoreName)
		if err != nil {
			return []getListProductResponse{}, custom_errors.CustomError{
				Code:          http.StatusInternalServerError,
				MessageToSend: "Internal Server Error",
				Message:       err.Error(),
			}
		}
		products = append(products, product)
	}
	return products, custom_errors.CustomError{}
}
func (s storeRepo) GetDetailProduct(id string) (productDetails, custom_errors.CustomError) {
	var product productDetails
	err := connections.DbMySQL().QueryRow("SELECT p.id, p.name, p.description, p.category_id, p.price, p.amount, p.image_id, f.filename FROM products p JOIN files f ON p.image_id = f.id  WHERE p.id = ? LIMIT 1", id).Scan(&product.ID, &product.Name, &product.Description, &product.CategoryID, &product.Price, &product.Amount, &product.ImageID, &product.ImagePath)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			customErr := custom_errors.CustomError{
				Code:          http.StatusNotFound,
				MessageToSend: "Product not found",
				Message:       err.Error(),
			}
			return productDetails{}, customErr
		}
		customErr := custom_errors.CustomError{
			Code:          http.StatusInternalServerError,
			MessageToSend: "Internal Server Error",
			Message:       err.Error(),
		}
		return productDetails{}, customErr
	}
	return product, custom_errors.CustomError{}
}
func (s storeRepo) DeleteStoreProduct(productID string, userID string) (bool, custom_errors.CustomError) {
	_, err := connections.DbMySQL().Exec("CALL delete_store_product(?,?)", userID, productID)
	if err != nil {
		if err.Error() == "Error 1644 (45001): Product not found." {
			customErr := custom_errors.CustomError{
				Code:          http.StatusNotFound,
				MessageToSend: "Product not found",
				Message:       err.Error(),
			}
			return true, customErr
		}
		customErr := custom_errors.CustomError{
			Code:          http.StatusInternalServerError,
			MessageToSend: "Internal Server Error",
			Message:       err.Error(),
		}
		return true, customErr
	}
	return false, custom_errors.CustomError{}
}

func factoryStoreRepo() iRepo {
	if repo == (storeRepo{}) {
		repo = storeRepo{}
	}
	return repo
}
