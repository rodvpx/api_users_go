package user

import (
	"context"
	"testing"
)

func TestUserService_Create_Validation(t *testing.T) {
	repo := NewInMemoryRepository()
	service := NewUserService(repo)

	// Tenta criar usuário sem nome para testar a regra de validação na Service
	input := CreateUserInput{
		Name:  "",
		Email: "email@valido.com",
	}

	_, err := service.Create(context.Background(), input)
	if err == nil {
		t.Error("esperava erro de validação ao criar usuário sem nome, mas recebeu nil")
	}
}
