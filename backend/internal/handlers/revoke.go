package handlers

import (
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (cfg *App) revoke(w http.ResponseWriter, r *http.Request) {
	if !validateRequest(w, r, "POST", false) {
		return
	}

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	byte_hashed_refresh, err := bcrypt.GenerateFromPassword([]byte(cookie.Value), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Printf("Couldn't hash or smtg: %v", err)
		return
	}

	cfg.revokeRefresh(string(byte_hashed_refresh))
	log.Println("Revoked refresh token")
}
