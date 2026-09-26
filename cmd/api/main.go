package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	fmt.Println("API iniciando...")

	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
				return

			}

			w.Header().Set("Content-Type", "application/json")

			u := User{
				ID:    1,
				Name:  "Marcio",
				Email: "marcio@mail.com",
			}

			err := json.NewEncoder(w).Encode(u)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		})

	fmt.Println("Servidor rodando em http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("erro ao iniciar servidor: ", err)
	}
}
