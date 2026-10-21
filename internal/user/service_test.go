package user

import (
	"context"
	"errors"
	"testing"
)

func setupTestService() (Service, Repository) {
	repo := NewInMemoryRepository()
	service := NewUserService(repo)
	return service, repo
}

func TestUserService_Create(t *testing.T) {
	t.Run("criação com sucesso", func(t *testing.T) {
		service, _ := setupTestService()

		input := CreateUserInput{
			Name:  "João Souza",
			Email: "joao@email.com",
		}

		output, err := service.Create(context.Background(), input)
		if err != nil {
			t.Fatalf("não esperava erro ao criar usuário: %v", err)
		}

		if output.ID == 0 {
			t.Error("esperava ID gerado maior que 0")
		}
		if output.Name != input.Name || output.Email != input.Email {
			t.Errorf("dados de saída incorretos: %+v", output)
		}
	})

	t.Run("falha na validação de nome vazio", func(t *testing.T) {
		service, _ := setupTestService()

		input := CreateUserInput{
			Name:  "",
			Email: "semnome@email.com",
		}

		_, err := service.Create(context.Background(), input)
		if err == nil {
			t.Error("esperava erro ao criar usuário sem nome, mas recebeu nil")
		}
		if !errors.Is(err, ErrInvalidInput) {
			t.Errorf("esperava erro contendo ErrInvalidInput, mas recebeu: %v", err)
		}
	})
}

func TestUserService_GetAll(t *testing.T) {
	service, repo := setupTestService()

	// Insere usuários prévios
	_, _ = repo.Create(context.Background(), User{Name: "User 1", Email: "u1@email.com"})
	_, _ = repo.Create(context.Background(), User{Name: "User 2", Email: "u2@email.com"})

	outputs, err := service.GetAll(context.Background())
	if err != nil {
		t.Fatalf("não esperava erro no GetAll: %v", err)
	}

	if len(outputs) != 2 {
		t.Errorf("esperava 2 usuários, mas recebeu %d", len(outputs))
	}
}

func TestUserService_GetByID(t *testing.T) {
	service, repo := setupTestService()

	created, _ := repo.Create(context.Background(), User{Name: "Carlos", Email: "carlos@email.com"})

	t.Run("sucesso ao buscar ID existente", func(t *testing.T) {
		output, err := service.GetByID(context.Background(), created.ID)
		if err != nil {
			t.Fatalf("não esperava erro ao buscar ID existente: %v", err)
		}
		if output.ID != created.ID || output.Name != "Carlos" {
			t.Errorf("dados incorretos retornados: %+v", output)
		}
	})

	t.Run("erro de usuário não encontrado", func(t *testing.T) {
		_, err := service.GetByID(context.Background(), 999)
		if err == nil {
			t.Error("esperava erro ao buscar ID 999, mas recebeu nil")
		}
		if !errors.Is(err, ErrUserNotFound) {
			t.Errorf("esperava erro ErrUserNotFound, recebeu %v", err)
		}
	})
}

func TestUserService_Update(t *testing.T) {
	service, repo := setupTestService()

	created, _ := repo.Create(context.Background(), User{Name: "Original", Email: "orig@email.com"})

	t.Run("sucesso ao atualizar nome parcialmente", func(t *testing.T) {
		novoNome := "Nome Atualizado"
		input := UpdateUserInput{
			Name: &novoNome,
		}

		output, err := service.Update(context.Background(), created.ID, input)
		if err != nil {
			t.Fatalf("falha ao atualizar: %v", err)
		}

		if output.Name != "Nome Atualizado" || output.Email != "orig@email.com" {
			t.Errorf("dados pós-update incorretos: %+v", output)
		}
	})

	t.Run("falha ao atualizar usuário inexistente", func(t *testing.T) {
		novoNome := "Novo"
		input := UpdateUserInput{Name: &novoNome}

		_, err := service.Update(context.Background(), 999, input)
		if err == nil {
			t.Error("esperava erro ao atualizar ID 999, mas recebeu nil")
		}
	})
}

func TestUserService_Delete(t *testing.T) {
	service, repo := setupTestService()

	created, _ := repo.Create(context.Background(), User{Name: "Para Deletar", Email: "del@email.com"})

	t.Run("sucesso ao deletar", func(t *testing.T) {
		err := service.Delete(context.Background(), created.ID)
		if err != nil {
			t.Fatalf("não esperava erro ao deletar: %v", err)
		}

		_, err = service.GetByID(context.Background(), created.ID)
		if err == nil {
			t.Error("esperava erro ao buscar usuário que foi deletado")
		}
	})

	t.Run("falha ao deletar ID inexistente", func(t *testing.T) {
		err := service.Delete(context.Background(), 999)
		if err == nil {
			t.Error("esperava erro ao deletar ID 999, mas recebeu nil")
		}
	})
}
