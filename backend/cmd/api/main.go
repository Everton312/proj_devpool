package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/Everton312/proj-devpool/internal/config"
)

func main() {

	mux := http.NewServeMux()

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}
	defer db.Close()

	fmt.Println("BANCO: OK")

	http.HandleFunc("/", handleHello)

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleHello(w http.ResponseWriter, _ *http.Request) {
	wc, err := w.Write([]byte("Hello World\n"))
	if err != nil {
		slog.Error("error writing response", "err", err)
		return
	}

	fmt.Println("%d bytes written\n", wc)
}
