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
	http.HandleFunc("GET /summary", middleware.AuthGuard(http.HandlerFunc(cfg.summary), cfg.JWT))
	http.HandleFunc("PATCH /users/settings", middleware.AuthGuard(http.HandlerFunc(cfg.updatePercent), cfg.JWT))
	http.HandleFunc("POST /users/change-password", middleware.AuthGuard(http.HandlerFunc(cfg.changePassword), cfg.JWT))

	fmt.Println("Starting Server")
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	fmt.Println("Stopping Server")
}
