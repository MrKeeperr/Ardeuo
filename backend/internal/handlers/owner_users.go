package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"ardeuo_backend/internal/models"
	"ardeuo_backend/internal/repository"
	"ardeuo_backend/internal/security"

	"github.com/go-chi/chi/v5"
)

type OwnerCreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	RoleID   int    `json:"rol_id"`
}

type OwnerUpdateUserRequest struct {
	RoleID int  `json:"rol_id"`
	Activo bool `json:"activo"`
}

func OwnerListUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	propietarioID, ok := getPropietarioIDFromRequest(w, r)
	if !ok {
		return
	}

	users, err := repository.ListUsersByPropietario(r.Context(), propietarioID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al consultar usuarios"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func OwnerCreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	propietarioID, ok := getPropietarioIDFromRequest(w, r)
	if !ok {
		return
	}

	var req OwnerCreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Formato JSON inválido"})
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Username y password son requeridos"})
		return
	}
	if !isValidUsername(req.Username) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Username inválido"})
		return
	}
	if len(req.Password) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "La contraseña debe tener al menos 8 caracteres"})
		return
	}
	if !isAllowedOwnerRole(req.RoleID) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Rol no permitido"})
		return
	}

	hash, err := security.HashPassword(req.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al encriptar contraseña"})
		return
	}

	userID, err := repository.CreateUserForPropietario(r.Context(), propietarioID, req.Username, hash, req.RoleID)
	if err != nil {
		if isUniqueViolation(err) {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": "Username ya existe"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al crear usuario"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"user_id": userID,
		"message": "Usuario creado",
	})
}

func OwnerUpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	propietarioID, ok := getPropietarioIDFromRequest(w, r)
	if !ok {
		return
	}

	userID, err := parseIDParam(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "usuario_id inválido"})
		return
	}

	var req OwnerUpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Formato JSON inválido"})
		return
	}
	if !isAllowedOwnerRole(req.RoleID) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Rol no permitido"})
		return
	}

	if err := repository.UpdateUserRoleStatus(r.Context(), propietarioID, userID, req.RoleID, req.Activo); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Usuario no encontrado"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al actualizar usuario"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"user_id": userID,
		"message": "Usuario actualizado",
	})
}

func isAllowedOwnerRole(roleID int) bool {
	return roleID == 2 || roleID == 3
}
