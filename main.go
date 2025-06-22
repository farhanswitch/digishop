package main

import (
	"fmt"
	"log"
	"net/http"

	"digishop/configs"
	"digishop/connections"
	"digishop/modules/users"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	initModules()
	router := chi.NewRouter()

	initPlugins(router)
	internalModules(router)

	router.Get("/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"messagee":"hello world"}`))
	}))
	var port uint16 = configs.GetConfig().Service.Port
	log.Printf("Server running on port %d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}

func initModules() {
	configs.InitModule("./env/local.env")
	connections.DbMySQL()
}
func internalModules(router *chi.Mux) {
	users.InitModule(router)
}
func initPlugins(router *chi.Mux) {
	router.Use(middleware.Recoverer)
}
