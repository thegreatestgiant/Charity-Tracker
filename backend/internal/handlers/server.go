package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/thegreatestgiant/Charity-Tracker/internal/middleware"
)

func StartServer(cfg *App) {
	check := func(jti uuid.UUID) bool {
		// jti, err := uuid.Parse(jtiStr)
		// if err != nil {
		// 	log.Printf("Couldn't get jti uuid: %v", err)
		// 	return false
		// }
		return cfg.blacklisted(jti)
	}
	port := os.Getenv("APP_PORT")

	http.HandleFunc("GET /health", cfg.PingDB)
	http.HandleFunc("POST /register", cfg.Register)
	http.HandleFunc("POST /login", cfg.Login)
	http.HandleFunc("POST /logout", cfg.Logout)
	http.HandleFunc("POST /refresh", middleware.AuthGuard(http.HandlerFunc(cfg.refresh), cfg.JWT, check))
	http.HandleFunc("POST /entries", middleware.AuthGuard(http.HandlerFunc(cfg.setEntry), cfg.JWT, check))
	http.HandleFunc("GET /entries", middleware.AuthGuard(http.HandlerFunc(cfg.getEntries), cfg.JWT, check))
	http.HandleFunc("GET /summary", middleware.AuthGuard(http.HandlerFunc(cfg.summary), cfg.JWT, check))
	http.HandleFunc("PATCH /users/settings", middleware.AuthGuard(http.HandlerFunc(cfg.updatePercent), cfg.JWT, check))
	http.HandleFunc("POST /users/change-password", middleware.AuthGuard(http.HandlerFunc(cfg.changePassword), cfg.JWT, check))

	fmt.Println("Starting Server")
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	fmt.Println("Stopping Server")
}
