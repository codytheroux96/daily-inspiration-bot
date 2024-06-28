package main

import (
	"database/sql"
	"encoding/json"
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

var db *sql.DB

type Quote struct {
	ID     int
	Text   string `json:"q"`
	Author string `json:"a"`
	Posted bool   `json:"posted"`
}

func main() {
	log.Println("Starting up the application...")

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
	client := resty.New()
	resp, err := client.R().Get("https://zenquotes.io/api/quotes")
	if err != nil {
		return err
	}

	var quotes []Quote
	err = json.Unmarshal(resp.Body(), &quotes)
	if err != nil {
		return err
	}

	for _, quote := range quotes {
		_, err := db.Exec("INSERT INTO quotes (text, author, posted) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE text=text", quote.Text, quote.Author, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func getUnpostedQuote() (Quote, error) {
	var quote Quote
	query := "SELECT id, text, author, posted FROM quotes WHERE posted = FALSE ORDER BY RAND() LIMIT 1"
	row := db.QueryRow(query)
	err := row.Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Posted)
	if err != nil {
		return quote, err
	}
	return quote, nil
}

func markQuoteAsPosted(id int) {
	query := "UPDATE quotes SET posted = TRUE WHERE id = ?"
	_, err := db.Exec(query, id)
	if err != nil {
		log.Printf("Error marking quote as posted: %v", err)
	}
}
