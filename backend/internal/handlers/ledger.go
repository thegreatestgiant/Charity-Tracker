package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type EntryType string

const (
	Paycheck       EntryType = "paycheck"
	Donation       EntryType = "donation"
	RunningBalance EntryType = "running_balance"
)

func (e EntryType) IsValid() bool {
	switch e {
	case Paycheck, Donation, RunningBalance:
		return true
	}
	return false
}

type Ledger struct {
	UserID           uuid.UUID `json:"user_id"`
	TransactionID    int       `json:"transaction_id"`
	LedgerEntry      EntryType `json:"ledger_entry"` // Maps to your 'entry' enum
	Amount           float64   `json:"amount"`       // DECIMAL(18,2)
	CharityOwed      float64   `json:"charity_owed"` // Use pointers if these can be NULL
	CharityFulfilled float64   `json:"charity_fulfilled"`
	TransactionDate  time.Time `json:"transaction_date"`
}

func (cfg *App) setEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Need Post", http.StatusMethodNotAllowed)
		log.Println("Wasn't POST")
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unsupported Content-Type", http.StatusNoContent)
		log.Println("Need json")
		return
	}
	uuid, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "UUID didn't come through", http.StatusBadRequest)
		log.Println("Couldn't get the cookies user_id...")
		return
	}

	sqlInsert := "INSERT INTO Ledgers (user_id, ledger_entry, amount, charity_owed, charity_fulfilled) VALUES ($1, $2, $3, $4, $5)"
	entry := Ledger{}

	json.NewDecoder(r.Body).Decode(&entry)
	defer r.Body.Close()

	if !entry.LedgerEntry.IsValid() {
		http.Error(w, "Invalid ledger entry type", http.StatusBadRequest)
		log.Println("Ledger Entry was not of valid type")
		return
	}

	_, err := cfg.DB.Query(sqlInsert, uuid, entry.LedgerEntry, entry.Amount, entry.CharityOwed, entry.CharityFulfilled)
	if err != nil {
		http.Error(w, "Bad Query", http.StatusBadRequest)
		log.Printf("Couldn't add to db: %v", err)
		return
	}
}

func (cfg *App) getEntries(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Bad Method", http.StatusMethodNotAllowed)
		log.Println("Bad Method")
		return
	}
	entries := []Ledger{}
	query := "SELECT * FROM Ledgers WHERE user_id=$1"

	uuid, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "UUID didn't come through", http.StatusBadRequest)
		log.Println("Couldn't get the cookies user_id...")
		end(w, r, entries)
		return
	}

	rows, err := cfg.DB.Query(query, uuid)
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
	// Allows your frontend to actually read the data
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "The Ledger entries")
	json.NewEncoder(w).Encode(entries)
}
