package stores

import (
	custom_errors "digishop/utilities/errors"
	"log"

	"github.com/google/uuid"
)

type storeService struct {
	repo iRepo
}

var service storeService

func (s storeService) RegisterStoreSrv(store storeData) (bool, custom_errors.CustomError) {
	currentStore, customErr := s.repo.GetStoreByUserID(store.UserID)
	if customErr != (custom_errors.CustomError{}) {
		return true, customErr
	}
	if currentStore != (storeData{}) {
		return true, custom_errors.CustomError{
			Code:          400,
			MessageToSend: "Store already exists",
			Message:       "This user already has a store",
		}
	}
	strUUID, err := uuid.NewV7()
	if err != nil {
		log.Println(err)
		return true, custom_errors.CustomError{
			Code:          500,
			Message:       err.Error(),
			MessageToSend: "Internal Server Error",
		}
	}
	store.ID = strUUID.String()
	return s.repo.RegisterStore(store)
}
func (s storeService) GetStoreByUserIDSrv(userID string) (storeData, custom_errors.CustomError) {
	return s.repo.GetStoreByUserID(userID)
}
func (s storeService) UpdateStoreSrv(store storeData) (bool, custom_errors.CustomError) {
	return s.repo.UpdateStore(store)
}
func (s storeService) CreateNewProductSrv(product productRequest) (bool, custom_errors.CustomError) {
	_, err := s.repo.GetCategoryNameByID(product.CategoryID)
	if err != nil {
		return true, custom_errors.CustomError{
			Code:          400,
			Message:       err.Error(),
			MessageToSend: "Invalid Category ID",
		}
	}
	store, customErr := s.repo.GetStoreByUserID(product.UserID)
	if customErr != (custom_errors.CustomError{}) {
		return true, custom_errors.CustomError{
			Code:          400,
			Message:       customErr.Message,
			MessageToSend: "You have no store",
		}
	}
	product.StoreID = store.ID
	_, err = s.repo.GetProductImagePathByID(product.ImageID)
	if err != nil {
		return true, custom_errors.CustomError{
			Code:          400,
			Message:       err.Error(),
			MessageToSend: "Invalid Image ID",
		}
	}
	strUUID, err := uuid.NewV7()
	if err != nil {
		log.Println(err)
		return true, custom_errors.CustomError{
			Code:          500,
			Message:       err.Error(),
			MessageToSend: "Internal Server Error",
		}
	}
	product.ID = strUUID.String()
	return s.repo.CreateNewProduct(product)
}
func factoryStoreService(repo iRepo) storeService {
	if service == (storeService{}) {
		service = storeService{
			repo: repo,
		}
	}
	return service
}
