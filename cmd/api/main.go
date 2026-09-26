package main

import (
	"api_users_go/internal/user"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("API iniciando...")

	http.HandleFunc("/", user.GetUserHandler)

	fmt.Println("Servidor rodando em http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("erro ao iniciar servidor: ", err)
	}
}
