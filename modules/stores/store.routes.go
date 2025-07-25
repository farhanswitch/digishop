package stores

import (
	"digishop/middlewares"

	"github.com/go-chi/chi/v5"
)

func InitModule(router *chi.Mux) {
	repo := factoryStoreRepo()
	controller := factoryStoreController(repo)

	router.With(middlewares.AuthMiddleware).Post("/store", controller.registerStoreCtrl)
}
