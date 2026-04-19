package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type setting struct {
	Percent float64 `json:"donation_percentage"`
}

func (cfg *App) setting(w http.ResponseWriter, r *http.Request) {
	if !validateRequest(w, r, "PATCH", true) {
		return
	}
	user_id := getUUID(w, r)
	if user_id == uuid.Nil {
		return
	}
	query := "UPDATE Users SET donation_percentage=$1 WHERE user_id=$2"
	setting := setting{}

	json.NewDecoder(r.Body).Decode(&setting)
	defer r.Body.Close()

	_, err := cfg.DB.Exec(query, setting.Percent, user_id)
	if err != nil {
		http.Error(w, "Not Updated", http.StatusInternalServerError)
		log.Printf("Couldn't update DB: %v", err)
	}

	w.Header().Set("Content-Type", "application/text")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Updated Donation Percent")
	cfg.getDonationPercent(user_id)
}
