package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB es el "Pool" global de conexiones que usaremos en todo el proyecto
var DB *pgxpool.Pool

func ConnectDB(databaseURL string) error { // conexion a la bd
	config, err := pgxpool.ParseConfig(databaseURL) // lee config env
	if err != nil {
		return fmt.Errorf("error al leer la configuración de la BD: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config) // crea pre-conexiones para carga rapida
	if err != nil {
		return fmt.Errorf("error al conectar con PostgreSQL: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil { // verifica que la conexión es válida
		return fmt.Errorf("la base de datos no responde: %v", err)
	}

	DB = pool
	fmt.Println("✅ Conexión exitosa a la base de datos")
	return nil
}
