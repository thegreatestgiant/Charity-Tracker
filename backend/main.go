package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/thegreatestgiant/Charity-Tracker/internals/db"
	"github.com/thegreatestgiant/Charity-Tracker/internals/handlers"
)

func main() {
	godotenv.Load(".env")

	db, err := db.OpenDB(os.Getenv("DB_DEV_URL"))
	if err != nil {
		log.Fatal(err)
	}

	cfg := &handlers.App{
		DB: db,
	}

	handlers.StartServer(cfg)
}
