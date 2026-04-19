package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type refresh struct {
	Token string `json:"refresh_token"`
}

func (cfg *App) refresh(w http.ResponseWriter, r *http.Request) {
	if !validateRequest(w, r, "POST", true) {
		return
	}

	var refresh refresh
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

	json.NewDecoder(r.Body).Decode(&refresh)
	defer r.Body.Close()

	hashed_refresh := cfg.getRefresh(user_id, cookie.Expires)

	err = bcrypt.CompareHashAndPassword([]byte(hashed_refresh), []byte(cookie.Value))
	if err != nil {
		log.Printf("Bad password: %v", err)
		return
	}

	cfg.denyList(jti)
	cfg.generateTokensWithCookies(w, user_id)
}
