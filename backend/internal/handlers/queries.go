package handlers

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

func (cfg *App) denyList(jti uuid.UUID) {
	query := "INSERT INTO denylist VALUES ($1, $2)"
	_, err := cfg.DB.Exec(query, jti, time.Now().Local().Add(time.Hour*24))
	if err != nil {
		log.Printf("Bad jti or time: %v", jti)
	}
}

func (cfg *App) blacklisted(jti uuid.UUID) bool {
	query := "SELECT jti FROM denylist WHERE jti=$1"
	row := cfg.DB.QueryRow(query, jti)
	if err := row.Scan(&uuid.UUID{}); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Not in denylist: %v", jti)
			return false
		}
		log.Default().Printf("Something went wrong: %v", err)
		return false
	}
	return true
}

func (cfg *App) getRefresh(user_id uuid.UUID, expires time.Time) string {
	query := "SELECT token FROM refresh_tokens WHERE user_id=$1 AND expires_at=$2"
	token := ""
	row := cfg.DB.QueryRow(query, user_id, expires)
	if err := row.Scan(&token); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No such token: %v", token)
			return ""
		}
		log.Default().Printf("Something went wrong: %v", err)
		return ""
	}
	return token
}

func (cfg *App) addRefresh(token string, user_id uuid.UUID, expires time.Time) {
	query := "INSERT INTO refresh_tokens (token, user_id, expires_at) VALUES ($1, $2, $3)"
	_, err := cfg.DB.Exec(query, token, user_id, expires)
	if err != nil {
		log.Println("Couldn't creat refresh token")
	}
}

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
	query := "SELECT ledger_entry, amount, charity_owed, charity_fulfilled FROM Ledgers WHERE user_id=$1 ORDER BY transaction_date DESC Limit 1"
	var entry Ledger

	row := cfg.DB.QueryRow(query, user_id)
	if err := row.Scan(&entry.LedgerEntry, &entry.Amount, &entry.CharityOwed, &entry.CharityFulfilled); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Bad uuid: %v. Or the date was wrong: %v", user_id)
			return entry
		}
		log.Default().Printf("No such Entry: %v", err)
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
