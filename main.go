package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/go-resty/resty/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Quote struct {
	ID     int
	Text   string `json:"text"`
	Author string `json:"author"`
	Posted bool   `json:"posted"`
}

func main() {
	var db *sql.DB
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN is not found in the environment")
	}

	channelID := os.Getenv("CHANNELID")
	if channelID == "" {
		log.Fatal("CHANNELID is not found in the environment")
	}

	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Fatal("DSN is not found in the environment")
	}

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("error pinging the database: %v", err)
	}
}
