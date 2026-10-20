package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// Função auxiliar para instanciar a pilha completa em memória
func setupTestHandler() (*UserHandler, Repository) {
	repo := NewInMemoryRepository()
	service := NewUserService(repo)
	handler := NewUserHandler(service)
	return handler, repo
}

func TestGetUserHandler(t *testing.T) {
	handler, _ := setupTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()

	handler.GetUserHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("esperava status %d, mas recebeu %d", http.StatusOK, rec.Code)
	}

	expectedHeader := "application/json"
	if got := rec.Header().Get("Content-Type"); got != expectedHeader {
		t.Errorf("esperava header Content-Type %q, mas recebeu %q", expectedHeader, got)
	}

	var users []UserOutput
	if err := json.Unmarshal(rec.Body.Bytes(), &users); err != nil {
		t.Fatalf("erro ao deserializar JSON de resposta: %v", err)
	}

	if len(users) != 0 {
		t.Errorf("esperava lista vazia de usuários, mas recebeu %d itens", len(users))
	}
}

func TestCreateUserHandler(t *testing.T) {
	handler, _ := setupTestHandler()

	bodyJSON := `{"name":"Maria Silva","email":"maria@email.com"}`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	handler.CreateUserHandler(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("esperava status %d, mas recebeu %d", http.StatusCreated, rec.Code)
	}

	var createdUser UserOutput
	if err := json.Unmarshal(rec.Body.Bytes(), &createdUser); err != nil {
		t.Fatalf("erro ao deserializar JSON de resposta: %v", err)
	}

	if createdUser.ID == 0 {
		t.Error("esperava ID válido maior que 0")
	}

	if createdUser.Name != "Maria Silva" || createdUser.Email != "maria@email.com" {
		t.Errorf("dados retornados incorretos: %+v", createdUser)
	}
}

func TestCreateUserHandler_InvalidJSON(t *testing.T) {
	handler, _ := setupTestHandler()

	invalidJSON := `{"name": "Maria", "email":`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	handler.CreateUserHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("esperava status %d, mas recebeu %d", http.StatusBadRequest, rec.Code)
	}
}

func TestGetUserByIDHandler_NotFound(t *testing.T) {
	handler, _ := setupTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
	req.SetPathValue("id", "999")

	rec := httptest.NewRecorder()

	handler.GetUserByIDHandler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("esperava status %d, mas recebeu %d", http.StatusNotFound, rec.Code)
	}
}

func TestUpdateUserHandler(t *testing.T) {
	handler, repo := setupTestHandler()

	u, err := repo.Create(context.Background(), User{Name: "Carlos", Email: "carlos@email.com"})
	if err != nil {
		t.Fatalf("falha ao criar usuário base: %v", err)
	}

	patchJSON := `{"name":"Carlos Eduardo"}`
	idStr := strconv.Itoa(u.ID)

	req := httptest.NewRequest(http.MethodPatch, "/users/"+idStr, bytes.NewBufferString(patchJSON))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", idStr)

	rec := httptest.NewRecorder()

	handler.UpdateUserHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("esperava status %d, mas recebeu %d", http.StatusOK, rec.Code)
	}

	var updatedUser UserOutput
	if err := json.Unmarshal(rec.Body.Bytes(), &updatedUser); err != nil {
		t.Fatalf("erro ao deserializar JSON: %v", err)
	}

	if updatedUser.Name != "Carlos Eduardo" || updatedUser.Email != "carlos@email.com" {
		t.Errorf("dados pós-update incorretos: %+v", updatedUser)
	}
}

func TestDeleteUserHandler(t *testing.T) {
	handler, repo := setupTestHandler()

	u, err := repo.Create(context.Background(), User{Name: "To Delete", Email: "delete@email.com"})
	if err != nil {
		t.Fatalf("falha ao criar usuário base: %v", err)
	}

	idStr := strconv.Itoa(u.ID)
	req := httptest.NewRequest(http.MethodDelete, "/users/"+idStr, nil)
	req.SetPathValue("id", idStr)

	rec := httptest.NewRecorder()

	handler.DeleteUserHandler(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("esperava status %d, mas recebeu %d", http.StatusNoContent, rec.Code)
	}

	_, err = repo.GetByID(context.Background(), u.ID)
	if err == nil {
		t.Error("esperava erro ao buscar usuário deletado, mas ele ainda existe")
	}
}
