package handlers

import (
	"encoding/json"
	"net/http"

	"ardeuo_backend/internal/security"
)

// Estructura de lo que esperamos recibir de React
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Estructura de lo que le devolveremos a React
type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

// LoginHandler procesa la solicitud de inicio de sesión
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Leer el JSON enviado por el usuario
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Formato JSON inválido"})
		return
	}

	// ---------------------------------------------------------
	// AQUI EN LA FASE B IREMOS A LA BASE DE DATOS:
	// dbUser, err := repository.GetUserByEmail(req.Email)
	// Si no existe -> Error 401
	// ---------------------------------------------------------

	// DATOS SIMULADOS COMO SI VINIERAN DE LA BD (Para probar la Fase A)
	dbHashedPassword, _ := security.HashPassword("123456") // En la vida real, este hash ya está guardado en PostgreSQL
	dbUserID := 5
	dbRoleID := 2 // Operario

	// 2. Comparar la contraseña enviada con el Hash de la Base de Datos
	if !security.CheckPasswordHash(req.Password, dbHashedPassword) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Credenciales incorrectas"})
		return
	}

	// 3. Generar el Token JWT
	tokenString, err := security.GenerateJWT(dbUserID, dbRoleID, req.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al generar el token"})
		return
	}

	// 4. Devolver éxito a React
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponse{
		Token:   tokenString,
		Message: "Login exitoso",
	})
}
