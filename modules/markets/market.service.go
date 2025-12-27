package markets

import (
	custom_errors "digishop/utilities/errors"
	"log"

	"github.com/google/uuid"
)

type marketService struct {
	repo iRepo
}

var service marketService

func (m marketService) GetAllCategorySrv() ([]category, custom_errors.CustomError) {

	return m.repo.GetAllCategory()
}
func (m marketService) GetListProductByCategorySrv(categoryID string) ([]productData, custom_errors.CustomError) {
	return m.repo.GetListProductByCategory(categoryID)
}
func (m marketService) GetProductDetailByIDSrv(productID string) (productDetail, custom_errors.CustomError) {
	return m.repo.GetProductDetailByID(productID)
}
func (m marketService) ExploreProductsSrv(search string) ([]productData, custom_errors.CustomError) {
	return m.repo.ExploreProducts(search)
}
func (m marketService) ManageCartSrv(userID string, productID string, quantity int) custom_errors.CustomError {
	strUUID, err := uuid.NewV7()
	if err != nil {
		log.Println(err)
		return custom_errors.CustomError{
			Code:          500,
			Message:       err.Error(),
			MessageToSend: "Internal Server Error",
		}
	}
	return m.repo.ManageCart(strUUID.String(), userID, productID, quantity)
}
func (m marketService) GetUserCartsSrv(userID string) ([]cartData, custom_errors.CustomError) {
	return m.repo.GetUserCarts(userID)
}
func (m marketService) GetUserNotificationsSrv(userID string) ([]notificationData, custom_errors.CustomError) {
	return m.repo.GetUserNotifications(userID)
}
func factoryMarketService(repo iRepo) marketService {
	if service == (marketService{}) {
		service = marketService{
			repo: repo,
		}
	}
	return service
}
