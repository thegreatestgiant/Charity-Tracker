package handlers

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *App) refresh(w http.ResponseWriter, r *http.Request) {
	if !validateRequest(w, r, "POST", false) {
		return
	}

	user_id := getUUID(w, r)
	jti := getUUID(w, r)
	if jti == uuid.Nil || user_id == uuid.Nil {
		return
	}

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	hashed_refresh := cfg.getRefresh(user_id)

	err = bcrypt.CompareHashAndPassword([]byte(hashed_refresh), []byte(cookie.Value))
	if err != nil {
		log.Printf("Bad password: %v", err)
		return
	}

	cfg.denyList(jti)
	cfg.generateJWTWithCookies(w, user_id)
}
