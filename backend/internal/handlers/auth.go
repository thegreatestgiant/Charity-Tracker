package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/thegreatestgiant/Charity-Tracker/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

type body struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (cfg *App) Register(w http.ResponseWriter, r *http.Request) {
	if !validateRequest(w, r, "POST", true) {
		return
	}

	var body body
	json.NewDecoder(r.Body).Decode(&body)
	defer r.Body.Close()

	log.Printf("Here's the body: %v", body)

	if cfg.userExists(body.Email, body.Username) {
		http.Error(w, "Already exists", http.StatusConflict)
		log.Println("Username or email already exist")
		return
	}

	if _, err := mail.ParseAddress(body.Email); err != nil {
		http.Error(w, "Invalid Email Format", http.StatusBadRequest)
		log.Printf("Invalid Email: %s", body.Email)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Printf("Couldn't hash or smtg: %v", err)
		return
	}

	err = cfg.setUser(body.Email, body.Username, string(hashedPassword))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Printf("Couldn't add to db: %v", err)
		return
	}

	fmt.Fprintln(w, "Updated password")
	log.Println("Updated Password")
}

func (cfg *App) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	claims, err := auth.Verifyer(cookie.Value, cfg.JWT)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Printf("Couldn't get claims", err)
		return
	}

	jti, err := uuid.Parse(claims.ID)
	if err != nil {
		log.Printf("Couldn't get jti uuid: %v", err)
		return
	}
	cfg.denyList(jti)
	cfg.revoke(w, r)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Path:     "/",
	})
	fmt.Fprintln(w, "Logging out ")
	log.Println("Logging out")
}

func (cfg *App) Login(w http.ResponseWriter, r *http.Request) {
	if !validateRequest(w, r, "POST", true) {
		return
	}

	var body body
	json.NewDecoder(r.Body).Decode(&body)
	defer r.Body.Close()

	user_id, hashedPass := cfg.getUser(body.Username)
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(body.Password))
	if err != nil {
		log.Printf("Bad password: %v", err)
		return
	}

	cfg.generateTokensWithCookies(w, user_id)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Set jwtCookie and RefreshCookie")
}
