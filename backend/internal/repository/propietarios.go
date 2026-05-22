package repository

import (
	"context"
	"time"

	"ardeuo_backend/internal/models"
)

const propietarioQueryTimeout = 5 * time.Second

func GetPropietarioByID(ctx context.Context, propietarioID int) (*models.Propietario, error) {
	queryCtx, cancel := context.WithTimeout(ctx, propietarioQueryTimeout)
	defer cancel()

	row := DB.QueryRow(
		queryCtx,
		`SELECT propietario_id, nombre_completo, email, telefono, direccion_id, activo
		 FROM propietarios
		 WHERE propietario_id = $1`,
		propietarioID,
	)

	var propietario models.Propietario
	if err := row.Scan(
		&propietario.ID,
		&propietario.NombreCompleto,
		&propietario.Email,
		&propietario.Telefono,
		&propietario.DireccionID,
		&propietario.Activo,
	); err != nil {
		return nil, err
	}

	return &propietario, nil
}
