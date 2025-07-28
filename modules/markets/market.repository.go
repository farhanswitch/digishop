package markets

import (
	"digishop/connections"
	"errors"

	custom_errors "digishop/utilities/errors"

	"net/http"
	"time"
)

func convertTimestamp(value interface{}) (time.Time, error) {
	var t time.Time

	switch v := value.(type) {
	case time.Time:
		// Sudah bertipe time.Time
		t = v
	case []byte:
		// Konversi []byte -> string -> time.Time
		timeStr := string(v)
		t, err := time.Parse("2006-01-02 15:04:05", timeStr)
		if err != nil {
			return time.Now(), err
		}
		return t, nil
	case string:
		// Kalau MySQL driver langsung kirim string
		t, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return time.Now(), err
		}
		return t, nil
	default:
		return time.Now(), errors.New("unsupported timestamp type")
	}
	return t, nil
}

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
func (m marketRepo) GetProductDetailByID(productID string) (productDetail, custom_errors.CustomError) {
	var product productDetail
	err := connections.DbMySQL().QueryRow("SELECT p.id, p.name, p.price, s.name, f.filename, p.description, c.name, p.amount FROM products p JOIN stores s ON p.store_id = s.id JOIN files f ON p.image_id = f.id JOIN categories c ON p.category_id = c.id WHERE p.id = ? LIMIT 1", productID).Scan(&product.ID, &product.Name, &product.Price, &product.StoreName, &product.ImagePath, &product.Description, &product.CategoryName, &product.Amount)
	if err != nil {
		return productDetail{}, custom_errors.CustomError{
			Code:          http.StatusInternalServerError,
			MessageToSend: "Internal Server Error",
			Message:       err.Error(),
		}
	}
	return product, custom_errors.CustomError{}
}
func (m marketRepo) ExploreProducts(search string) ([]productData, custom_errors.CustomError) {
	var products []productData
	results, err := connections.DbMySQL().Query("CALL explore_products(?)", search)
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
func (m marketRepo) ManageCart(userID string, productID string, quantity int) custom_errors.CustomError {
	_, err := connections.DbMySQL().Exec("CALL manage_cart(?, ?, ?)", userID, productID, quantity)
	if err != nil {
		return custom_errors.CustomError{
			Code:          http.StatusInternalServerError,
			MessageToSend: "Internal Server Error",
			Message:       err.Error(),
		}
	}
	return custom_errors.CustomError{}
}
func (m marketRepo) GetUserCarts(userID string) ([]cartData, custom_errors.CustomError) {
	var carts []cartData
	results, err := connections.DbMySQL().Query("CALL get_user_cart(?)", userID)
	if err != nil {
		return []cartData{}, custom_errors.CustomError{
			Code:          http.StatusInternalServerError,
			MessageToSend: "Internal Server Error",
			Message:       err.Error(),
		}
	}
	for results.Next() {
		var data cartData
		err := results.Scan(&data.ProductID, &data.ProductName, &data.ProductPrice, &data.ProductAmount, &data.ProductImagePath, &data.StoreName, &data.CartQuantity)
		if err != nil {
			return []cartData{}, custom_errors.CustomError{
				Code:          http.StatusInternalServerError,
				MessageToSend: "Internal Server Error",
				Message:       err.Error(),
			}
		}
		carts = append(carts, data)
	}
	return carts, custom_errors.CustomError{}
}
func (m marketRepo) GetUserNotifications(userID string) ([]notificationData, custom_errors.CustomError) {
	var notifications []notificationData
	results, err := connections.DbMySQL().Query("SELECT title, description, created_at FROM notifications WHERE user_id = ? ORDER BY created_at DESC", userID)
	if err != nil {
		return []notificationData{}, custom_errors.CustomError{
			Code:          http.StatusInternalServerError,
			MessageToSend: "Internal Server Error",
			Message:       err.Error(),
		}
	}
	for results.Next() {
		var valueCreated interface{}
		var data notificationData
		err := results.Scan(&data.Title, &data.Description, &valueCreated)
		if err != nil {
			return []notificationData{}, custom_errors.CustomError{
				Code:          http.StatusInternalServerError,
				MessageToSend: "Internal Server Error",
				Message:       err.Error(),
			}
		}
		data.CreatedAt, err = convertTimestamp(valueCreated)
		if err != nil {
			return []notificationData{}, custom_errors.CustomError{
				Code:          http.StatusInternalServerError,
				MessageToSend: "Internal Server Error",
				Message:       err.Error(),
			}
		}
		notifications = append(notifications, data)
	}
	return notifications, custom_errors.CustomError{}
}
func factoryMarketRepository() iRepo {
	if repo == (marketRepo{}) {
		repo = marketRepo{}
	}
	return repo
}
