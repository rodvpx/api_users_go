package main

import (
	"api_users_go/internal/user"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("API iniciando...")

	repo := user.NewInMemoryRepository()

	handler := user.NewUserHandler(repo)

	http.HandleFunc("GET /users", handler.GetUserHandler)
	http.HandleFunc("POST /users", handler.CreateUserHandler)
	http.HandleFunc("GET /users/{id}", handler.GetUserByIDHandler)
	http.HandleFunc("PATCH /users/{id}", handler.UpdateUserHandler)
	http.HandleFunc("DELETE /users/{id}", handler.DeleteUserHandler)

	fmt.Println("Servidor rodando em http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("erro ao iniciar servidor: ", err)
	}
}
