package user

import (
	"testing"
)

func TestToUserOutput(t *testing.T) {
	u := User{
		ID:    1,
		Name:  "Ana Silva",
		Email: "ana@email.com",
	}

	output := ToUserOutput(u)

	if output.ID != u.ID || output.Name != u.Name || output.Email != u.Email {
		t.Errorf("ToUserOutput falhou. Esperava %+v, mas recebeu %+v", u, output)
	}
}

func TestToUserOutputList(t *testing.T) {
	t.Run("converte lista com elementos", func(t *testing.T) {
		users := []User{
			{ID: 1, Name: "User 1", Email: "u1@email.com"},
			{ID: 2, Name: "User 2", Email: "u2@email.com"},
		}

		outputs := ToUserOutputList(users)

		if len(outputs) != 2 {
			t.Fatalf("esperava 2 elementos, mas recebeu %d", len(outputs))
		}
		if outputs[0].Name != "User 1" || outputs[1].Name != "User 2" {
			t.Errorf("dados incorretos na conversão de lista: %+v", outputs)
		}
	})

	t.Run("converte lista vazia ou nula", func(t *testing.T) {
		outputs := ToUserOutputList([]User{})

		if len(outputs) != 0 {
			t.Errorf("esperava lista vazia de outputs, mas recebeu %d itens", len(outputs))
		}
	})
}
