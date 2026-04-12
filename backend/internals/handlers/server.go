package handlers

import (
	"fmt"
	"net/http"
	"os"
)

func StartServer(cfg *App) {
	port := os.Getenv("APP_PORT")

	http.HandleFunc("/health", cfg.PingDB)
	fmt.Println("Starting Server")
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	fmt.Println("Stopping Server")
}
