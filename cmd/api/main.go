package main

import (
	"api_users_go/internal/database"
	"api_users_go/internal/middleware"
	"api_users_go/internal/user"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("API iniciando...")

	dsn := "postgres://postgres:postgres@localhost:5432/api_users_go_db?sslmode=disable"

	db, err := database.NewPostgresDB(dsn)
	if err != nil {
		fmt.Printf("Erro fatal ao conectar ao banco de dados: %v\n", err)
		return
	}

	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Erro fatal ao executar migrations: %v", err)
	}

	repo := user.NewPostgresRepository(db)
	userService := user.NewUserService(repo)
	handler := user.NewUserHandler(userService)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /users", handler.GetUserHandler)
	mux.HandleFunc("POST /users", handler.CreateUserHandler)
	mux.HandleFunc("GET /users/{id}", handler.GetUserByIDHandler)
	mux.HandleFunc("PATCH /users/{id}", handler.UpdateUserHandler)
	mux.HandleFunc("DELETE /users/{id}", handler.DeleteUserHandler)

	var handlerComMiddlewares http.Handler = mux
	handlerComMiddlewares = middleware.Logger(handlerComMiddlewares)
	handlerComMiddlewares = middleware.Recovery(handlerComMiddlewares)
	handlerComMiddlewares = middleware.CORS(handlerComMiddlewares)

	fmt.Println("Servidor rodando em http://localhost:8080")

	err = http.ListenAndServe(":8080", handlerComMiddlewares)
	if err != nil {
		fmt.Println("erro ao iniciar servidor: ", err)
	}
}
