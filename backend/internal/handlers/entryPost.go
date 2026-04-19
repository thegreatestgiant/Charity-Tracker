package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type EntryType string

const (
	Paycheck EntryType = "paycheck"
	Donation EntryType = "donation"
)

func (e EntryType) IsValid() bool {
	switch e {
	case Paycheck, Donation:
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
	if !validateRequest(w, r, "POST", true) {
		return
	}

	user_id := getUUID(w, r)
	if user_id == uuid.Nil {
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

	percent := float64(cfg.getDonationPercent(user_id))

	owed := 0.0
	fulfilled := 0.0
	if entry.LedgerEntry == Paycheck {
		owed = entry.Amount * (percent / 100.0)
		owed = math.Round(owed*100) / 100
		log.Printf("Recieved Paycheck, amount owed: %.2f", owed)
	} else {
		fulfilled = (entry.Amount / cfg.getAmountOwed(user_id)) * 100
		fulfilled = math.Round(fulfilled*100) / 100
		log.Printf("Recieved Donation, fulfilled %.2f%%", fulfilled)
	}

	_, err := cfg.DB.Query(sqlInsert, user_id, entry.LedgerEntry, entry.Amount, owed, fulfilled)
	if err != nil {
		http.Error(w, "Bad Query", http.StatusBadRequest)
		log.Printf("Couldn't add to db: %v", err)
		return
	}

	entry = cfg.getEntry(user_id, entry.TransactionDate)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Recieved Entry")
	json.NewEncoder(w).Encode(entry)
}
