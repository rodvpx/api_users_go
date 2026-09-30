package database

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgresDB(dataSourceName string) (*sql.DB, error) {

	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("Erro ao abrir conexão com o banco: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Erro ao conectar com o banco de dados: %w", err)
	}

	fmt.Println("Conexão com o PostgreSQL estabelecida com sucesso!")
	return db, nil
}
