package main

import (
	"fmt"
	"gs-app/backend/internal/handler"
	"gs-app/backend/internal/middleware"
	"gs-app/backend/internal/store"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

var defaultPackSizes = []int{250, 500, 1000, 2000, 5000}

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("PORT is not found in the env - Setting to 8000 by default")

		port = "8000"
	}

	logger := log.Default()

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	router.Use(middleware.LoggerMiddleware(logger))

	packStore := store.NewPacksStore(defaultPackSizes, logger)

	apiCfg := &handler.APIConfig{
		PackStore: packStore,
	}

	v1Router := chi.NewRouter()

	v1Router.Get("/health", apiCfg.CheckHealthHandler)
	v1Router.Get("/packs", apiCfg.GetPacksHandler)
	v1Router.Post("/calculate", apiCfg.CalculateHandler)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	logger.Printf("Server starting om port: %s", port)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
