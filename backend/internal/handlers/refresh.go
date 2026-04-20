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
	jti := getJti(w, r)
	if jti == uuid.Nil || user_id == uuid.Nil {
		return
	}

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	hashed_refresh := cfg.getRefresh(user_id)
	log.Printf("Hashed: %v ", hashed_refresh)
	log.Printf("cookie normal: %v ", cookie.Value)

	byte_hashed_refresh, err := bcrypt.GenerateFromPassword([]byte(cookie.Value), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Printf("Couldn't hash or smtg: %v", err)
		return
	}

	log.Printf("cookie encrypted: %v ", string(byte_hashed_refresh))

	err = bcrypt.CompareHashAndPassword([]byte(hashed_refresh), []byte(cookie.Value))
	if err != nil {
		log.Printf("Bad password: %v", err)
		return
	}

	cfg.denyList(jti)
	cfg.revoke(w, r)
	cfg.generateTokensWithCookies(w, user_id)
}
