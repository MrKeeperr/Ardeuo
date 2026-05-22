package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"ardeuo_backend/internal/handlers"
	"ardeuo_backend/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  No se encontró archivo .env. Asegúrate de tener uno en la raíz.")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("❌ ERROR: La variable DATABASE_URL no está definida en el .env")
	}

	err = repository.ConnectDB(dbURL)
	if err != nil {
		log.Fatal(err) // Si falla la BD, el servidor no debe arrancar
	}

	// Cerrar el pool de forma limpia cuando el servidor se apague en el futuro
	defer repository.DB.Close()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/api/users", handlers.GetUsers)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	fmt.Printf("🚀 El microservicio de Demeter está listo en %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
