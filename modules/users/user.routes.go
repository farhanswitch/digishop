package users

import "github.com/go-chi/chi/v5"

func InitModule(router *chi.Mux) {
	repo := factoryUserRepo()
	controller := factoryUserController(repo)

	router.Post("/user/register", controller.RegisterUserCtrl)
}
