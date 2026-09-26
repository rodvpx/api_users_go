package main

import (
	"api_users_go/internal/user"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("API iniciando...")

	http.HandleFunc("GET /users", user.GetUserHandler)
	http.HandleFunc("POST /users", user.CreateUserHandler)
	http.HandleFunc("GET /users/{id}", user.GetUserByIDHandler)
	http.HandleFunc("PATCH /users/{id}", user.UpdateUserHandler)
	http.HandleFunc("DELETE /users/{id}", user.DeleteUserHandler)

	fmt.Println("Servidor rodando em http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("erro ao iniciar servidor: ", err)
	}
}
