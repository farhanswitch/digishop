package files

import (
	"digishop/middlewares"

	"github.com/go-chi/chi/v5"
)

func InitModule(router *chi.Mux) {
	controller := factoryFileController()
	router.With(middlewares.AuthMiddleware).Post("/file/product-photo/upload", controller.UploadFileCtrl)
	router.Get("/file/{filename}", controller.GetFileCtrl)

}
