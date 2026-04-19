package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type App struct {
	DB  *sql.DB
	JWT []byte
}

func validateRequest(w http.ResponseWriter, r *http.Request, method string, requiresBody bool) bool {
	if r.Method != method {
		http.Error(w, "Need "+method, http.StatusMethodNotAllowed)
		log.Printf("Wasn't %s", method)
		return false
	}
	if requiresBody && r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unsupported Content-Type", http.StatusNoContent)
		log.Println("Need json")
		return false
	}
	return true
}

func getUUID(w http.ResponseWriter, r *http.Request) uuid.UUID {
	user_id, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "UUID didn't come through", http.StatusBadRequest)
		log.Println("Couldn't get the cookies user_id...")
		return uuid.Nil
	}
	return user_id
}
