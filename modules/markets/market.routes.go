package markets

import (
	"digishop/middlewares"

	"github.com/go-chi/chi/v5"
)

func InitModule(router *chi.Mux) {
	repo := factoryMarketRepository()
	controller := factoryMarketController(repo)

	router.Get("/market/categories", controller.getAllCategoryCtrl)
	router.Get("/market/products-by-category", controller.getListProductByCategoryCtrl)
	router.Get("/market/product-detail/{id}", controller.getProductDetailByIDCtrl)
	router.Get("/market/explore-products", controller.exploreProductsCtrl)

	router.With(middlewares.AuthMiddleware).Post("/market/cart/submit", controller.manageCartCtrl)
	router.With(middlewares.AuthMiddleware).Get("/market/cart", controller.getUserCartCtrl)
}
