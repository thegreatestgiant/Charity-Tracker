package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *App) getEntries(w http.ResponseWriter, r *http.Request) {
	if !validateRequest(w, r, "GET", false) {
		return
	}
	entries := []Ledger{}
	query := "SELECT * FROM Ledgers WHERE user_id=$1"

	user_id := getUUID(w, r)
	if user_id == uuid.Nil {
		return
	}

	rows, err := cfg.DB.Query(query, user_id)
	if err != nil {
		http.Error(w, "No entries", http.StatusNoContent)
		log.Default().Printf("Bad query: %v ", err)
		end(w, r, entries)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var entry Ledger
		if err = rows.Scan(&entry.UserID, &entry.TransactionID, &entry.LedgerEntry, &entry.Amount, &entry.CharityOwed, &entry.CharityFulfilled, &entry.TransactionDate); err != nil {
			log.Printf("Couldn't scan row: %v", err)
			end(w, r, entries)
			return
		}
		entries = append(entries, entry)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, "IDK", http.StatusInternalServerError)
		log.Printf("Not sure why, but here's an error: %v", err)
	}

	end(w, r, entries)
}

func end(w http.ResponseWriter, r *http.Request, entries []Ledger) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "The Ledger entries")
	json.NewEncoder(w).Encode(entries)
	log.Printf("Sent the ledgers: %v", entries)
}
