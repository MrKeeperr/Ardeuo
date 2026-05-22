package repository

import (
	"context"
	"time"

	"ardeuo_backend/internal/models"
)

const locationsQueryTimeout = 5 * time.Second

func ListPaises(ctx context.Context) ([]models.Pais, error) {
	queryCtx, cancel := context.WithTimeout(ctx, locationsQueryTimeout)
	defer cancel()

	rows, err := DB.Query(
		queryCtx,
		`SELECT pais_id, nombre_pais
		 FROM paises
		 ORDER BY nombre_pais`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	paises := make([]models.Pais, 0)
	for rows.Next() {
		var pais models.Pais
		if err := rows.Scan(&pais.PaisID, &pais.Nombre); err != nil {
			return nil, err
		}
		paises = append(paises, pais)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return paises, nil
}

func ListDepartamentos(ctx context.Context, paisID *int) ([]models.Departamento, error) {
	queryCtx, cancel := context.WithTimeout(ctx, locationsQueryTimeout)
	defer cancel()

	query := `SELECT departamento_id, nombre_departamento, pais_id
		FROM departamentos`
	args := []any{}
	if paisID != nil {
		query += ` WHERE pais_id = $1`
		args = append(args, *paisID)
	}
	query += ` ORDER BY nombre_departamento`

	rows, err := DB.Query(queryCtx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	departamentos := make([]models.Departamento, 0)
	for rows.Next() {
		var departamento models.Departamento
		if err := rows.Scan(&departamento.DepartamentoID, &departamento.Nombre, &departamento.PaisID); err != nil {
			return nil, err
		}
		departamentos = append(departamentos, departamento)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return departamentos, nil
}

func ListMunicipios(ctx context.Context, departamentoID *int) ([]models.Municipio, error) {
	queryCtx, cancel := context.WithTimeout(ctx, locationsQueryTimeout)
	defer cancel()

	query := `SELECT municipio_id, nombre_municipio, departamento_id
		FROM municipios`
	args := []any{}
	if departamentoID != nil {
		query += ` WHERE departamento_id = $1`
		args = append(args, *departamentoID)
	}
	query += ` ORDER BY nombre_municipio`

	rows, err := DB.Query(queryCtx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	municipios := make([]models.Municipio, 0)
	for rows.Next() {
		var municipio models.Municipio
		if err := rows.Scan(&municipio.MunicipioID, &municipio.Nombre, &municipio.DepartamentoID); err != nil {
			return nil, err
		}
		municipios = append(municipios, municipio)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return municipios, nil
}
