package markets

import (
	custom_errors "digishop/utilities/errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// mockRepo is a mock of iRepo interface
type mockRepo struct {
	ctrl     *gomock.Controller
	recorder *mockRepoRecorder
}

// mockRepoRecorder is the mock recorder for mockRepo
type mockRepoRecorder struct {
	mock *mockRepo
}

// NewMockRepo creates a new mock instance
func NewMockRepo(ctrl *gomock.Controller) *mockRepo {
	mock := &mockRepo{ctrl: ctrl}
	mock.recorder = &mockRepoRecorder{mock}
	return mock
}

func (m *mockRepo) GetAllCategory() ([]category, custom_errors.CustomError) {
	return []category{
		{ID: "1", Name: "Category 1"},
		{ID: "2", Name: "Category 2"},
	}, custom_errors.CustomError{}
}

func (m *mockRepo) GetListProductByCategory(categoryID string) ([]productData, custom_errors.CustomError) {
	return []productData{
		{ID: "1", Name: "Product 1", Price: 100, StoreName: "Store 1", ImagePath: "image1.jpg"},
		{ID: "2", Name: "Product 2", Price: 200, StoreName: "Store 2", ImagePath: "image2.jpg"},
	}, custom_errors.CustomError{}
}

func (m *mockRepo) GetProductDetailByID(productID string) (productDetail, custom_errors.CustomError) {
	return productDetail{
		ID:           "1",
		Name:         "Product 1",
		Price:        100,
		StoreName:    "Store 1",
		ImagePath:    "image1.jpg",
		Description:  "Description 1",
		CategoryName: "Category 1",
		Amount:       10,
	}, custom_errors.CustomError{}
}

func (m *mockRepo) ExploreProducts(search string) ([]productData, custom_errors.CustomError) {
	return []productData{
		{ID: "1", Name: "Product 1", Price: 100, StoreName: "Store 1", ImagePath: "image1.jpg"},
		{ID: "2", Name: "Product 2", Price: 200, StoreName: "Store 2", ImagePath: "image2.jpg"},
	}, custom_errors.CustomError{}
}

func (m *mockRepo) ManageCart(userID string, productID string, quantity int) custom_errors.CustomError {
	return custom_errors.CustomError{}
}

func (m *mockRepo) GetUserCarts(userID string) ([]cartData, custom_errors.CustomError) {
	return []cartData{
		{ProductID: "1", ProductName: "Product 1", ProductPrice: 100, ProductAmount: 10, ProductImagePath: "image1.jpg", StoreName: "Store 1", CartQuantity: 2},
		{ProductID: "2", ProductName: "Product 2", ProductPrice: 200, ProductAmount: 20, ProductImagePath: "image2.jpg", StoreName: "Store 2", CartQuantity: 1},
	}, custom_errors.CustomError{}
}

func (m *mockRepo) GetUserNotifications(userID string) ([]notificationData, custom_errors.CustomError) {
	return []notificationData{
		{Title: "Title 1", Description: "Description 1", CreatedAt: time.Now()},
		{Title: "Title 2", Description: "Description 2", CreatedAt: time.Now()},
	}, custom_errors.CustomError{}
}

func TestGetAllCategorySrv(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepo(ctrl)
	service := marketService{repo: mockRepo}

	categories, err := service.GetAllCategorySrv()

	assert.Empty(t, err)
	assert.Len(t, categories, 2)
	assert.Equal(t, "Category 1", categories[0].Name)
	assert.Equal(t, "Category 2", categories[1].Name)
}

func TestGetListProductByCategorySrv(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepo(ctrl)
	service := marketService{repo: mockRepo}

	products, err := service.GetListProductByCategorySrv("1")

	assert.Empty(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "Product 1", products[0].Name)
	assert.Equal(t, "Product 2", products[1].Name)
}

func TestGetProductDetailByIDSrv(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepo(ctrl)
	service := marketService{repo: mockRepo}

	product, err := service.GetProductDetailByIDSrv("1")

	assert.Empty(t, err)
	assert.Equal(t, "Product 1", product.Name)
	assert.Equal(t, float64(100), product.Price)
	assert.Equal(t, "Store 1", product.StoreName)
}

func TestExploreProductsSrv(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepo(ctrl)
	service := marketService{repo: mockRepo}

	products, err := service.ExploreProductsSrv("test")

	assert.Empty(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "Product 1", products[0].Name)
	assert.Equal(t, "Product 2", products[1].Name)
}

func TestManageCartSrv(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepo(ctrl)
	service := marketService{repo: mockRepo}

	err := service.ManageCartSrv("user1", "product1", 2)

	assert.Empty(t, err)
}

func TestGetUserCartsSrv(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepo(ctrl)
	service := marketService{repo: mockRepo}

	carts, err := service.GetUserCartsSrv("user1")

	assert.Empty(t, err)
	assert.Len(t, carts, 2)
	assert.Equal(t, "Product 1", carts[0].ProductName)
	assert.Equal(t, "Product 2", carts[1].ProductName)
}

func TestGetUserNotificationsSrv(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepo(ctrl)
	service := marketService{repo: mockRepo}

	notifications, err := service.GetUserNotificationsSrv("user1")

	assert.Empty(t, err)
	assert.Len(t, notifications, 2)
	assert.Equal(t, "Title 1", notifications[0].Title)
	assert.Equal(t, "Title 2", notifications[1].Title)
}
