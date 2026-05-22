package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"regexp"
	"strings"

	"ardeuo_backend/internal/models"
	"ardeuo_backend/internal/repository"
	"ardeuo_backend/internal/security"

	"github.com/jackc/pgx/v5/pgconn"
)

type RegisterRequest struct {
	NombreCompleto string `json:"nombre_completo"`
	Email          string `json:"email"`
	Telefono       string `json:"telefono"`
	PaisID         int    `json:"pais_id"`
	DepartamentoID int    `json:"departamento_id"`
	MunicipioID    int    `json:"municipio_id"`
	CodigoPostal   string `json:"codigo_postal"`
	Direccion      string `json:"direccion"`
	Username       string `json:"username"`
	Password       string `json:"password"`
}

type RegisterResponse struct {
	UserID        int    `json:"user_id"`
	PropietarioID int    `json:"propietario_id"`
	Message       string `json:"message"`
}

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9._-]{3,30}$`)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Formato JSON inválido"})
		return
	}

	req.NombreCompleto = strings.TrimSpace(req.NombreCompleto)
	req.Email = strings.TrimSpace(req.Email)
	req.Telefono = strings.TrimSpace(req.Telefono)
	req.CodigoPostal = strings.TrimSpace(req.CodigoPostal)
	req.Direccion = strings.TrimSpace(req.Direccion)
	req.Username = strings.TrimSpace(req.Username)

	if req.NombreCompleto == "" || req.Email == "" || req.Telefono == "" || req.Username == "" || req.Password == "" || req.Direccion == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Campos obligatorios faltantes"})
		return
	}
	if req.PaisID <= 0 || req.DepartamentoID <= 0 || req.MunicipioID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ubicacion invalida"})
		return
	}
	if !isValidEmail(req.Email) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Email inválido"})
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

	passwordHash, err := security.HashPassword(req.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al encriptar contraseña"})
		return
	}

	propietario := models.Propietario{
		NombreCompleto: req.NombreCompleto,
		Email:          req.Email,
		Telefono:       req.Telefono,
	}
	direccion := models.DireccionInput{
		MunicipioID:  req.MunicipioID,
		CodigoPostal: req.CodigoPostal,
		Direccion:    req.Direccion,
	}

	userID, propietarioID, err := repository.CreatePropietarioWithUser(r.Context(), propietario, direccion, req.Username, passwordHash, 2)
	if err != nil {
		if isUniqueViolation(err) {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": "Username o email ya existe"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al registrar usuario"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(RegisterResponse{
		UserID:        userID,
		PropietarioID: propietarioID,
		Message:       "Registro exitoso",
	})
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func isValidEmail(value string) bool {
	parsed, err := mail.ParseAddress(value)
	return err == nil && parsed.Address == value
}

func isValidUsername(value string) bool {
	return usernameRegex.MatchString(value)
}
