package repository

import (
	"context"
	"time"

	"ardeuo_backend/internal/models"

	"github.com/jackc/pgx/v5"
)

const registrationQueryTimeout = 5 * time.Second

func CreatePropietarioWithUser(ctx context.Context, propietario models.Propietario, direccion models.DireccionInput, username string, passwordHash string, roleID int) (int, int, error) {
	queryCtx, cancel := context.WithTimeout(ctx, registrationQueryTimeout)
	defer cancel()

	tx, err := DB.BeginTx(queryCtx, pgx.TxOptions{})
	if err != nil {
		return 0, 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(queryCtx)
		}
	}()

	var direccionID int
	err = tx.QueryRow(
		queryCtx,
		`INSERT INTO direcciones (municipio_id, codigo_postal, direccion)
		 VALUES ($1, $2, $3)
		 RETURNING direccion_id`,
		direccion.MunicipioID,
		direccion.CodigoPostal,
		direccion.Direccion,
	).Scan(&direccionID)
	if err != nil {
		return 0, 0, err
	}

	var propietarioID int
	err = tx.QueryRow(
		queryCtx,
		`INSERT INTO propietarios (nombre_completo, email, telefono, direccion_id, activo)
		 VALUES ($1, $2, $3, $4, true)
		 RETURNING propietario_id`,
		propietario.NombreCompleto,
		propietario.Email,
		propietario.Telefono,
		direccionID,
	).Scan(&propietarioID)
	if err != nil {
		return 0, 0, err
	}

	var userID int
	err = tx.QueryRow(
		queryCtx,
		`INSERT INTO usuarios (username, password_hash, rol_id, propietario_id, activo)
		 VALUES ($1, $2, $3, $4, true)
		 RETURNING usuario_id`,
		username,
		passwordHash,
		roleID,
		propietarioID,
	).Scan(&userID)
	if err != nil {
		return 0, 0, err
	}

	if err = tx.Commit(queryCtx); err != nil {
		return 0, 0, err
	}

	return userID, propietarioID, nil
}
