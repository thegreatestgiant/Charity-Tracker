package handlers

import (
	"fmt"
	"net/http"
	"os"
)

func StartServer(cfg *App) {
	port := os.Getenv("APP_PORT")

	http.HandleFunc("GET /health", cfg.PingDB)
	http.HandleFunc("POST /register", cfg.Register)
	http.HandleFunc("POST /login", cfg.Login)
	http.HandleFunc("POST /logout", Logout)
	// http.HandleFunc("GET /ledger", middleware.Authenticate(Ledger, cfg))

	fmt.Println("Starting Server")
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	fmt.Println("Stopping Server")
}
