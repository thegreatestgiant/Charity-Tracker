package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/thegreatestgiant/Charity-Tracker/internals/auth"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *App) Register(username, email, password string) (response string) {
	sqlInsert := "INSERT INTO Users (email,username,password_hash) VALUES ($1,$2,$3)"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	_, err = cfg.DB.Query(sqlInsert, email, username, hashedPassword)
	if err != nil {
		log.Fatal(err)
	}
	return "Success"
}

func (cfg *App) Login(username, password string, w http.ResponseWriter) (success bool) {
	sqlQuery := "SELECT user_id FROM users WHERE username=$1 AND password=$2"

	var uuid uuid.UUID
	row := cfg.DB.QueryRow(sqlQuery, username, password)
	if err := row.Scan(&uuid); err != nil {
		if err == sql.ErrNoRows {
			log.Fatal("Incorrect username %s or password", username)
		}
		log.Fatal("Incorrect username %s or password", username)
	}

	token, err := auth.MakeJWT(uuid, cfg.JWT, time.Hour*24)
	if err != nil {
		log.Fatal(err)
	}
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    token,
		HttpOnly: true,
		Secure:   false,
	}

	http.SetCookie(w, cookie)

	w.Write([]byte("Cookie has been set!"))

	return true
}
