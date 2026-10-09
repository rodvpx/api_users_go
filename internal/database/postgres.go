package database

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

func RunMigrations(db *sql.DB) error {

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("Erro ao criar driver de migration: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)

	if err != nil {
		return fmt.Errorf("Erro ao instaciar golang-migrate: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("Erro ao executar migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		fmt.Println("Migrations: nenhuma alteração pendente (banco atualizado)")
	} else {
		fmt.Println("Migrations executadas com sucesso!")
	}

	return nil
}
