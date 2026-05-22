package repository

import (
	"context"
	"time"

	"ardeuo_backend/internal/models"
)

const usersQueryTimeout = 5 * time.Second

func GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	queryCtx, cancel := context.WithTimeout(ctx, usersQueryTimeout)
	defer cancel()

	row := DB.QueryRow(
		queryCtx,
		`SELECT usuario_id, username, password_hash, rol_id, propietario_id, activo, fecha_creacion
		 FROM usuarios
		 WHERE username = $1`,
		username,
	)

	var user models.User
	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.RoleID,
		&user.PropietarioID,
		&user.Activo,
		&user.FechaCreacion,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	queryCtx, cancel := context.WithTimeout(ctx, usersQueryTimeout)
	defer cancel()

	row := DB.QueryRow(
		queryCtx,
		`SELECT usuario_id, username, password_hash, rol_id, propietario_id, activo, fecha_creacion
		 FROM usuarios
		 WHERE usuario_id = $1`,
		userID,
	)

	var user models.User
	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.RoleID,
		&user.PropietarioID,
		&user.Activo,
		&user.FechaCreacion,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUserPassword(ctx context.Context, userID int, passwordHash string) error {
	queryCtx, cancel := context.WithTimeout(ctx, usersQueryTimeout)
	defer cancel()

	commandTag, err := DB.Exec(
		queryCtx,
		`UPDATE usuarios SET password_hash = $1 WHERE usuario_id = $2`,
		passwordHash,
		userID,
	)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}

func UpdateUsername(ctx context.Context, userID int, username string) error {
	queryCtx, cancel := context.WithTimeout(ctx, usersQueryTimeout)
	defer cancel()

	commandTag, err := DB.Exec(
		queryCtx,
		`UPDATE usuarios SET username = $1 WHERE usuario_id = $2`,
		username,
		userID,
	)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}

func ListUsers(ctx context.Context) ([]models.User, error) {
	queryCtx, cancel := context.WithTimeout(ctx, usersQueryTimeout)
	defer cancel()

	rows, err := DB.Query(
		queryCtx,
		`SELECT usuario_id, username, rol_id, propietario_id, activo, fecha_creacion
		 FROM usuarios
		 ORDER BY usuario_id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.RoleID,
			&user.PropietarioID,
			&user.Activo,
			&user.FechaCreacion,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
