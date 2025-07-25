package main

import (
	"fmt"
	"log"
	"net/http"

	"digishop/configs"
	"digishop/connections"
	"digishop/modules/files"
	"digishop/modules/stores"
	"digishop/modules/users"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
	connections.ConnectRedis()
}
func internalModules(router *chi.Mux) {
	users.InitModule(router)
	files.InitModule(router)
	stores.InitModule(router)
}
func initPlugins(router *chi.Mux) {
	router.Use(middleware.Recoverer)
	// Middleware CORS
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any major browsers
	}))
}
