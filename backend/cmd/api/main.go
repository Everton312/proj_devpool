package main

import (
	"log"
	"fmt"
	"github.com/Everton312/proj-devpool/internal/config"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}
	defer db.Close()

	fmt.Println("Conectado ao PostgreSQL com sucesso!")
}