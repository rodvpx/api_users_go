package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestGetUserHandler(t *testing.T) {

	repo := NewInMemoryRepository()
	handler := NewUserHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)

	rec := httptest.NewRecorder()

	handler.GetUserHandler(rec, req)

	// a) Valida o Status Code
	if rec.Code != http.StatusOK {
		t.Errorf("Espera status %d, mas recebeu %d", http.StatusOK, rec.Code)
	}

	// b) Valida o Content-Type nos Headers
	expectedHeader := "application/json"
	if got := rec.Header().Get("Content-Type"); got != expectedHeader {
		t.Errorf("Espera header Content-Type %q, mas recebeu %q", expectedHeader, got)
	}

	// c) Valida o corpo da resposta (JSON)
	var users []User
	if err := json.Unmarshal(rec.Body.Bytes(), &users); err != nil {
		t.Fatalf("erro ao deserializar JSON da resposta: %v", err)
	}

	if len(users) != 0 {
		t.Errorf("esperava lista vazia de usuários, mas recebeu %d itens", len(users))
	}

}

func TestCreateUserHandler(t *testing.T) {

	repo := NewInMemoryRepository()
	handler := NewUserHandler(repo)

	bodyJSON := `{"name":"Maria Silva","email":"maria@email.com"}`

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	handler.CreateUserHandler(rec, req)

	// a) Verifica o Status Code (201 Created)
	if rec.Code != http.StatusCreated {
		t.Errorf("esperava status %d, mas recebeu %d", http.StatusCreated, rec.Code)
	}

	// b) Deserializa o JSON de resposta para verificar o ID gerado
	var createdUser User
	if err := json.Unmarshal(rec.Body.Bytes(), &createdUser); err != nil {
		t.Fatalf("erro ao deserializar JSON de resposta: %v", err)
	}

	if createdUser.ID == 0 {
		t.Error("esperava que o usuário recebesse um ID válido maior que 0")
	}

	if createdUser.Name != "Maria Silva" || createdUser.Email != "maria@email.com" {
		t.Errorf("dados retornados incorretos: %+v", createdUser)
	}

}

func TestCreateUserHandler_InvalidJSON(t *testing.T) {
	repo := NewInMemoryRepository()
	handler := NewUserHandler(repo)

	// JSON malformatado (faltando aspa e chave de fechamento)
	invalidJSON := `{"name": "Maria", "email":`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	handler.CreateUserHandler(rec, req)

	// Valida se retornou 400 Bad Request
	if rec.Code != http.StatusBadRequest {
		t.Errorf("esperava status %d, mas recebeu %d", http.StatusBadRequest, rec.Code)
	}
}

func TestGetUserByIDHandler_NotFound(t *testing.T) {
	repo := NewInMemoryRepository()
	handler := NewUserHandler(repo)

	// Tenta buscar um ID que não existe (ex: ID 999)
	req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
	// No Go 1.22+, SetPathValue simula o valor extraído da rota pelo multiplexer
	req.SetPathValue("id", "999")

	rec := httptest.NewRecorder()

	handler.GetUserByIDHandler(rec, req)

	// Valida se retornou 404 Not Found
	if rec.Code != http.StatusNotFound {
		t.Errorf("esperava status %d, mas recebeu %d", http.StatusNotFound, rec.Code)
	}
}

func TestGetUserByIDHandler_InvalidID(t *testing.T) {
	repo := NewInMemoryRepository()
	handler := NewUserHandler(repo)

	// Passa um ID com texto em vez de número
	req := httptest.NewRequest(http.MethodGet, "/users/abc", nil)
	req.SetPathValue("id", "abc")

	rec := httptest.NewRecorder()

	handler.GetUserByIDHandler(rec, req)

	// Valida se retornou 400 Bad Request
	if rec.Code != http.StatusBadRequest {
		t.Errorf("esperava status %d, mas recebeu %d", http.StatusBadRequest, rec.Code)
	}
}

func TestUpdateUserHandler(t *testing.T) {
	repo := NewInMemoryRepository()
	handler := NewUserHandler(repo)

	// Inserimos um usuário inicial na memória para atualizar
	u, err := repo.Create(User{Name: "Carlos", Email: "carlos@email.com"})
	if err != nil {
		t.Fatalf("falha ao criar usuário base: %v", err)
	}

	// Payload simulando PATCH apenas com o nome
	patchJSON := `{"name":"Carlos Eduardo"}`
	idStr := strconv.Itoa(u.ID)

	req := httptest.NewRequest(http.MethodPatch, "/users/"+idStr, bytes.NewBufferString(patchJSON))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", idStr)

	rec := httptest.NewRecorder()

	handler.UpdateUserHandler(rec, req)

	if rec.Code != http.StatusOK && rec.Code != http.StatusAccepted {
		t.Errorf("esperava status de sucesso no update, mas recebeu %d", rec.Code)
	}

	var updatedUser User
	if err := json.Unmarshal(rec.Body.Bytes(), &updatedUser); err != nil {
		t.Fatalf("erro ao deserializar JSON: %v", err)
	}

	if updatedUser.Name != "Carlos Eduardo" || updatedUser.Email != "carlos@email.com" {
		t.Errorf("dados pós-update incorretos: %+v", updatedUser)
	}
}

func TestDeleteUserHandler(t *testing.T) {
	repo := NewInMemoryRepository()
	handler := NewUserHandler(repo)

	// Criamos um usuário para ser removido
	u, err := repo.Create(User{Name: "To Delete", Email: "delete@email.com"})
	if err != nil {
		t.Fatalf("falha ao criar usuário base: %v", err)
	}

	idStr := strconv.Itoa(u.ID)
	req := httptest.NewRequest(http.MethodDelete, "/users/"+idStr, nil)
	req.SetPathValue("id", idStr)

	rec := httptest.NewRecorder()

	handler.DeleteUserHandler(rec, req)

	// Deletion bem-sucedida deve retornar 204 No Content
	if rec.Code != http.StatusNoContent {
		t.Errorf("esperava status %d, mas recebeu %d", http.StatusNoContent, rec.Code)
	}

	// Garante que o usuário realmente sumiu do repositório
	_, err = repo.GetByID(u.ID)
	if err == nil {
		t.Error("esperava erro ao buscar usuário deletado, mas ele ainda existe")
	}
}
