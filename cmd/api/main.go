package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("API iniciando...")

	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
				return

			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, `{"mensagem": "Olá, API!"}`)
		})

	fmt.Println("Servidor rodando em http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("erro ao iniciar servidor: ", err)
	}
}
