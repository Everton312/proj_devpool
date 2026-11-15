package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Everton312/proj-devpool/internal/config"
	"github.com/Everton312/proj-devpool/internal/handler"
	"github.com/Everton312/proj-devpool/internal/middleware"
	"github.com/Everton312/proj-devpool/internal/repository"
	"github.com/Everton312/proj-devpool/internal/service"
)

func main() {

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}
	defer db.Close()

	fmt.Println("BANCO: OK")

	repos := repository.New(db)
	services := service.New(repos)
	userHandler := handler.NewUserHandler(services)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", userHandler.CreateUser)
	mux.HandleFunc("GET /users", userHandler.GetAllUsers)
	mux.HandleFunc("GET /users/{id}", userHandler.GetUserById)
	mux.HandleFunc("DELETE /users/{id}", userHandler.DeleteUserById)
	mux.HandleFunc("PUT /users/{id}", userHandler.Update)

	fmt.Println("RODANDO 8080")
	log.Fatal(http.ListenAndServe(":8080", middleware.EnableCORS(mux)))
}
