package stores

import custom_errors "digishop/utilities/errors"

type storeData struct {
	ID      string `json:"id"`
	Name    string `json:"name" validate:"required,min=6,max=100"`
	Address string `json:"address" validate:"required,min=6,max=255"`
	UserID  string `json:"userID"`
}

type productRequest struct {
	ID          string  `json:"id"`
	CategoryID  string  `json:"categoryID" validate:"required"`
	UserID      string  `json:"userID"`
	StoreID     string  `json:"storeID"`
	ImageID     string  `json:"imageID" validate:"required"`
	Name        string  `json:"name" validate:"required,min=6,max=100"`
	Description string  `json:"description" validate:"required,min=6"`
	Price       float64 `json:"price" validate:"required,min=1"`
	Amount      int     `json:"amount" validate:"required,min=1"`
}

type updateProductRequest struct {
	ID          string  `json:"id" validate:"required"`
	CategoryID  string  `json:"categoryID" validate:"required"`
	UserID      string  `json:"userID"`
	StoreID     string  `json:"storeID"`
	ImageID     string  `json:"imageID" validate:"required"`
	Name        string  `json:"name" validate:"required,min=6,max=100"`
	Description string  `json:"description" validate:"required,min=6"`
	Price       float64 `json:"price" validate:"required,min=1"`
	Amount      int     `json:"amount" validate:"required,min=1"`
}

type getListProductRequest struct {
	PaginationPage uint ` validate:"required,min=1"`
	PaginationRows uint ` validate:"required,min=1"`
	Offset         uint
	UserID         string
	StoreID        string
	SortField      string `validate:"required"`
	SortOrder      string `validate:"required,oneof=asc desc ASC DESC"`
	Search         string
}
type getListProductResponse struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	CategoryName string  `json:"categoryName"`
	Price        float64 `json:"price"`
	Amount       int     `json:"amount"`
	ImagePath    string  `json:"imagePath"`
	StoreName    string  `json:"storeName"`
}
type productDetails struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	CategoryID  string  `json:"categoryID"`
	Price       float64 `json:"price"`
	Amount      int     `json:"amount"`
	ImageID     string  `json:"imageID"`
	ImagePath   string  `json:"imagePath"`
}
type iRepo interface {
	RegisterStore(store storeData) (bool, custom_errors.CustomError)
	GetStoreByUserID(id string) (storeData, custom_errors.CustomError)
	GetProductImagePathByID(id string) (string, error)
	GetCategoryNameByID(id string) (string, error)
	CreateNewProduct(product productRequest) (bool, custom_errors.CustomError)
	UpdateStore(store storeData) (bool, custom_errors.CustomError)
	CheckIsValidUserProduct(userID string, productID string) error
	UpdateProducts(product updateProductRequest) (bool, custom_errors.CustomError)
	GetListProduct(param getListProductRequest) ([]getListProductResponse, custom_errors.CustomError)
	GetDetailProduct(id string) (productDetails, custom_errors.CustomError)
	DeleteStoreProduct(productID string, userID string) (bool, custom_errors.CustomError)
}
