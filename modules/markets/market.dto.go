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

type iRepo interface {
	GetAllCategory() ([]category, custom_errors.CustomError)
	GetListProductByCategory(categoryID string) ([]productData, custom_errors.CustomError)
}
