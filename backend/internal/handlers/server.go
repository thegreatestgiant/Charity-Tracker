package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/thegreatestgiant/Charity-Tracker/internal/middleware"
)

func StartServer(cfg *App) {
	port := os.Getenv("APP_PORT")

	http.HandleFunc("GET /health", cfg.PingDB)
	http.HandleFunc("POST /register", cfg.Register)
	http.HandleFunc("POST /login", cfg.Login)
	http.HandleFunc("POST /logout", Logout)
	http.HandleFunc("POST /entries", middleware.AuthGuard(http.HandlerFunc(cfg.setEntry), cfg.JWT))
	http.HandleFunc("GET /entries", middleware.AuthGuard(http.HandlerFunc(cfg.getEntries), cfg.JWT))

	fmt.Println("Starting Server")
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	fmt.Println("Stopping Server")
}
