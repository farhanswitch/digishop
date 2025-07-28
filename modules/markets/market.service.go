package markets

import (
	custom_errors "digishop/utilities/errors"
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

func factoryMarketService(repo iRepo) marketService {
	if service == (marketService{}) {
		service = marketService{
			repo: repo,
		}
	}
	return service
}
