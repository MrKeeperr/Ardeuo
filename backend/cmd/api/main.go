package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"ardeuo_backend/internal/handlers"
	appmiddleware "ardeuo_backend/internal/middleware"
	"ardeuo_backend/internal/repository"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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

	allowedOrigins := defaultAllowedOrigins()
	if origins := os.Getenv("CORS_ALLOWED_ORIGINS"); origins != "" {
		allowedOrigins = splitAndTrim(origins)
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	}))

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Get("/api/users", handlers.GetUsers)

	r.Post("/api/login", handlers.LoginHandler)
	r.Post("/api/register", handlers.RegisterHandler)

	r.Get("/api/locations/paises", handlers.ListPaises)
	r.Get("/api/locations/departamentos", handlers.ListDepartamentos)
	r.Get("/api/locations/municipios", handlers.ListMunicipios)

	r.Route("/api/v1", func(router chi.Router) {
		router.Use(appmiddleware.RequireAuth) // Middleware para proteger las rutas dentro de /api/v1

		router.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"message": "¡Bienvenido a la bóveda secreta de Demeter!"}`))
		})

		router.Put("/account/password", handlers.ChangePasswordHandler)
		router.Put("/account/username", handlers.ChangeUsernameHandler)

		router.Route("/admin", func(admin chi.Router) {
			admin.Use(appmiddleware.RequireRole(1))
			admin.Get("/tenants", handlers.AdminListTenants)
			admin.Patch("/tenants/{id}/status", handlers.AdminUpdateTenantStatus)
			admin.Get("/usage", handlers.AdminTenantUsage)
			admin.Get("/health", handlers.AdminHealth)
		})

		router.Route("/owner", func(owner chi.Router) {
			owner.Use(appmiddleware.RequireRole(2))

			owner.Get("/users", handlers.OwnerListUsers)
			owner.Post("/users", handlers.OwnerCreateUser)
			owner.Patch("/users/{id}", handlers.OwnerUpdateUser)

			owner.Get("/cultivos", handlers.OwnerListCultivos)
			owner.Post("/cultivos", handlers.OwnerCreateCultivo)
			owner.Patch("/cultivos/{id}", handlers.OwnerUpdateCultivo)
			owner.Post("/cultivos/{id}/close", handlers.OwnerCloseCultivo)
			owner.Put("/cultivos/{id}/umbrales", handlers.OwnerUpsertUmbral)

			owner.Get("/nodos", handlers.OwnerListNodos)
			owner.Post("/nodos", handlers.OwnerCreateNodo)
			owner.Patch("/nodos/{id}", handlers.OwnerUpdateNodo)

			owner.Get("/umbrales", handlers.OwnerListUmbrales)

			owner.Get("/audit/eventos", handlers.OwnerListAuditEvents)
			owner.Get("/audit/nodos/{id}", handlers.OwnerListNodoAudit)
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

func defaultAllowedOrigins() []string {
	return []string{
		"http://localhost:3000",
		"http://127.0.0.1:3000",
		"http://localhost:5173",
		"http://127.0.0.1:5173",
	}
}

func splitAndTrim(raw string) []string {
	parts := strings.Split(raw, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value != "" {
			values = append(values, value)
		}
	}
	return values
}
