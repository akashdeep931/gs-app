package main

import (
	"fmt"
	"gs-app/backend/handler"
	"gs-app/backend/middleware"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

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
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	router.Use(middleware.LoggerMiddleware(logger))

	v1Router := chi.NewRouter()

	v1Router.Get("/health", handler.CheckHealthHandler)

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
