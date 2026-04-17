package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/thegreatestgiant/Charity-Tracker/internal/db"
	"github.com/thegreatestgiant/Charity-Tracker/internal/handlers"
)

func main() {
	godotenv.Load(".env")

	db, err := db.OpenDB(os.Getenv("DB_DEV_URL"))
	if err != nil {
		log.Fatal(err)
	}

	cfg := &handlers.App{
		DB:  db,
		JWT: []byte(os.Getenv("JWT_SECRET")),
	}

	handlers.StartServer(cfg)
}
