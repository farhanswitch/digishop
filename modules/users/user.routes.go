package users

import "github.com/go-chi/chi/v5"

func InitModule(router *chi.Mux) {
	repo := factoryUserRepo()
	controller := factoryUserController(repo)

	router.Get("/authenticate", controller.CheckAuthenticationCtrl)
	router.Post("/user/register", controller.RegisterUserCtrl)
	router.Post("/user/login", controller.LoginUserCtrl)
}
