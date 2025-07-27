package markets

import (
	"github.com/go-chi/chi/v5"
)

func InitModule(router *chi.Mux) {
	repo := factoryMarketRepository()
	controller := factoryMarketController(repo)

	router.Get("/market/categories", controller.getAllCategoryCtrl)
}
