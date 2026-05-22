package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"ardeuo_backend/internal/handlers"
	appmiddleware "ardeuo_backend/internal/middleware"
	"ardeuo_backend/internal/repository"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
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

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Get("/api/users", handlers.GetUsers)

	r.Post("/api/login", handlers.LoginHandler)

	r.Route("/api/v1", func(router chi.Router) {
		router.Use(appmiddleware.RequireAuth) // Middleware para proteger las rutas dentro de /api/v1

		router.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"message": "¡Bienvenido a la bóveda secreta de Demeter!"}`))
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	fmt.Printf("🚀 El microservicio de Ardeuo está listo en %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
