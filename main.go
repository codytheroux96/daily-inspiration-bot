package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	// "os/signal"
	// "syscall"

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
	log.Println("Starting up the application...")

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

	err = fetchAndStoreQuotes()
	if err != nil {
		log.Fatalf("Error fetching and storing quotes: %v", err)
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}
	defer dg.Close()

	log.Println("Bot is now running. Press CTRL+C to exit.")

	quote, err := getUnpostedQuote()
	if err != nil {
		log.Printf("error getting quote: %v", err)
	} else {
		_, err = dg.ChannelMessageSend(channelID, fmt.Sprintf("\"%s\" - %s", quote.Text, quote.Author))
		if err != nil {
			log.Printf("error sending message: %v", err)
		} else {
			markQuoteAsPosted(quote.ID)
		}
	}

	log.Println("Shutting down the application...")
}

func fetchAndStoreQuotes() error {
	return nil
}

func onReady() {

}

func dailyQuote() {

}

func getUnpostedQuote() (Quote, error) {
	var quote Quote

	return quote, nil
}

func markQuoteAsPosted(id int) {

}
