package user

import "errors"

var (
	ErrUserNotFound       = errors.New("Usuário não encontrado")
	ErrEmailAlreadyExists = errors.New("Este e-mail já está em uso")
	ErrInvalidInput       = errors.New("Dados de entrada inválidos")
)
