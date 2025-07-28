package markets

import custom_errors "digishop/utilities/errors"

type category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type productData struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	StoreName string  `json:"storeName"`
	ImagePath string  `json:"imagePath"`
}
type productDetail struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	StoreName    string  `json:"storeName"`
	ImagePath    string  `json:"imagePath"`
	Description  string  `json:"description"`
	CategoryName string  `json:"categoryName"`
	Amount       int     `json:"amount"`
}
type manageCartRequest struct {
	ProductID string `json:"productID" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required"`
	UserID    string `json:"userID" validate:"required"`
}
type cartData struct {
	ProductID        string  `json:"productID"`
	ProductName      string  `json:"productName"`
	ProductPrice     float64 `json:"productPrice"`
	ProductAmount    int     `json:"productAmount"`
	ProductImagePath string  `json:"productImagePath"`
	StoreName        string  `json:"storeName"`
	CartQuantity     int     `json:"cartQuantity"`
}
type iRepo interface {
	GetAllCategory() ([]category, custom_errors.CustomError)
	GetListProductByCategory(categoryID string) ([]productData, custom_errors.CustomError)
	GetProductDetailByID(productID string) (productDetail, custom_errors.CustomError)
	ExploreProducts(search string) ([]productData, custom_errors.CustomError)
	ManageCart(userID string, productID string, quantity int) custom_errors.CustomError
	GetUserCarts(userID string) ([]cartData, custom_errors.CustomError)
}
