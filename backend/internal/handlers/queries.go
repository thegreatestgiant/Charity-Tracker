package handlers

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
)

func (cfg *App) getDonationPercent(user_id uuid.UUID) float64 {
	query := "SELECT donation_percentage FROM users WHERE user_id=$1"
	percent := 10.0

	row := cfg.DB.QueryRow(query, user_id)
	if err := row.Scan(&percent); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Bad uuid: %v", user_id)
			return percent
		}
		log.Default().Printf("Couldn't get donation percent, using default: %v", err)
		return 10.0
	}

	log.Printf("Donation Percent: %.2f", percent)
	return percent
}

func (cfg *App) getEntry(user_id uuid.UUID) Ledger {
	query := "SELECT * FROM Ledgers WHERE user_id=$1 ORDER BY transaction_date DESC Limit 1"
	var entry Ledger

	row := cfg.DB.QueryRow(query, user_id)
	if err := row.Scan(&entry.UserID, &entry.TransactionID, &entry.LedgerEntry, &entry.Amount, &entry.CharityOwed, &entry.CharityFulfilled, &entry.TransactionDate); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Bad uuid: %v. Or the date was wrong: %v", user_id)
			return entry
		}
		log.Default().Printf("No such Entry", err)
		return entry
	}

	log.Printf("Entry: %v", entry)
	return entry
}

func (cfg *App) getAmountOwed(user_id uuid.UUID) float64 {
	query := "SELECT SUM(charity_owed) FROM Ledgers WHERE user_id=$1"
	owed := 0.0

	row := cfg.DB.QueryRow(query, user_id)
	if err := row.Scan(&owed); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Bad uuid: %v", user_id)
			return owed
		}
		log.Default().Printf("Couldn't sum owed: %v", err)
		return 10
	}

	log.Printf("Total Amount Owed: %.2f", owed)
	return owed
}

func (cfg *App) getAmountEarned(user_id uuid.UUID) float64 {
	query := "SELECT SUM(amount) FROM Ledgers WHERE user_id=$1 AND ledger_entry='paycheck'"
	earned := 0.0

	row := cfg.DB.QueryRow(query, user_id)
	if err := row.Scan(&earned); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Bad uuid: %v", user_id)
			return earned
		}
		log.Default().Printf("Couldn't sum fulfilled: %v", err)
		return 10
	}

	log.Printf("Total Amount Earned: %.2f", earned)
	return earned
}

func (cfg *App) getAmountDonated(user_id uuid.UUID) float64 {
	query := "SELECT SUM(amount) FROM Ledgers WHERE user_id=$1 AND ledger_entry='donation'"
	donated := 0.0

	row := cfg.DB.QueryRow(query, user_id)
	if err := row.Scan(&donated); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Bad uuid: %v", user_id)
			return donated
		}
		log.Default().Printf("Couldn't sum fulfilled: %v", err)
		return 10
	}

	log.Printf("Total Amount Donated: %.2f", donated)
	return donated
}

func (cfg *App) getAmountFulfilled(user_id uuid.UUID) float64 {
	fulfilled := (cfg.getAmountDonated(user_id) / cfg.getAmountOwed(user_id)) * 100

	log.Printf("Total Percent Fulfilled: %.2f%%", fulfilled)
	return fulfilled
}
