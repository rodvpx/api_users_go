package user

import (
	"database/sql"
	"errors"
	"fmt"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (p *postgresRepository) Create(u User) (User, error) {

	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`

	err := p.db.QueryRow(query, u.Name, u.Email).Scan(&u.ID)
	if err != nil {
		return User{}, fmt.Errorf("Erro ao inserir usuário: %w", err)
	}
	return u, nil
}

func (p *postgresRepository) GetAll() ([]User, error) {

	query := `SELECT id, name, email FROM users`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Erro ao buscar usuário: %w", err)
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User

		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, fmt.Errorf("Erro ao ler usuário %w", err)
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Erro na interação dos usuários: %w", err)
	}

	if users == nil {
		users = []User{}
	}

	return users, nil
}

func (p *postgresRepository) GetByID(id int) (User, error) {

	query := `SELECT id, name, email FROM users WHERE id = $1`
	var u User

	err := p.db.QueryRow(query, id).Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, errors.New("Usuário não econtrado")
		}
		return User{}, fmt.Errorf("Erro ao buscar usuário por ID: %w", err)
	}

	return u, nil
}

func (p *postgresRepository) Update(id int, u User) (User, error) {

	query := `UPDATE users SET name = $1, email = $2 WHERE id = $3`

	res, err := p.db.Exec(query, u.Name, u.Email, id)
	if err != nil {
		return User{}, fmt.Errorf("erro ao atualizar usuário: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return User{}, fmt.Errorf("Erro ao verificar linhas afetadas: %w", err)
	}
	if rows == 0 {
		return User{}, errors.New("Usuário não encontrado")
	}

	u.ID = id
	return u, nil
}

func (p *postgresRepository) Delete(id int) error {

	query := `DELETE FROM users WHERE id = $1`

	res, err := p.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Erro ao deletar usuário: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Erro ao verificar linhas afetadas: %w", err)
	}

	if rows == 0 {
		return errors.New("Usuário não encontrado")
	}

	return nil
}
