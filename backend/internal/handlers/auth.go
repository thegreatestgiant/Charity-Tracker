package handlers

import (
	"database/sql"
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

	sqlInsert := "INSERT INTO Users (email,username,password_hash,user_id) VALUES ($1,$2,$3,$4)"

	body := body{}
	json.NewDecoder(r.Body).Decode(&body)
	defer r.Body.Close()

	log.Printf("Here's the body: %v", body)

	if exists(body.Email, body.Username, cfg) {
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
	_, err = cfg.DB.Query(sqlInsert, body.Email, body.Username, hashedPassword, uuid.New())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Printf("Couldn't add to db: %v", err)
		return
	}

	fmt.Fprintln(w, "Updated password")
	log.Println("Updated Password")
}

func exists(email, username string, cfg *App) bool {
	query := "SELECT * FROM users WHERE email=$1 OR username=$2"
	// Will return nil if empty, and it doesn't exist
	err := cfg.DB.QueryRow(query, email, username).Scan()
	if err == sql.ErrNoRows {
		return false
	}
	return true
}

func Logout(w http.ResponseWriter, r *http.Request) {
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
	if r.Method != "POST" {
		http.Error(w, "Bad Method", http.StatusMethodNotAllowed)
		log.Println("Bad Method")
		return
	}

	body := body{}
	sqlQuery := "SELECT user_id,password_hash FROM users WHERE username=$1"

	json.NewDecoder(r.Body).Decode(&body)
	defer r.Body.Close()

	var uuid uuid.UUID
	var pass string
	row := cfg.DB.QueryRow(sqlQuery, body.Username)
	if err := row.Scan(&uuid, &pass); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Incorrect username %s or password", body.Username)
			return
		}
		log.Printf("Incorrect username %s or password", body.Username)
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(pass), []byte(body.Password))
	if err != nil {
		log.Printf("Bad password: %v", err)
		return
	}

	token, err := auth.MakeJWT(uuid, cfg.JWT, time.Hour*24)
	if err != nil {
		log.Println("Couldn't make token")
		return
	}
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    token,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	w.Write([]byte("Cookie has been set!"))
}
