package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const quoteTemplatePath = "quote.html.tmpl"

type quoteHandler struct {
	db *gorm.DB
}

func newQuoteHandler(db *gorm.DB) quoteHandler {
	migrate(db)
	return quoteHandler{db: db}
}

func migrate(db *gorm.DB) {
	var lock bool

	// Take out advisory lock in the database to prevent simultaneous migrations.
	log.Println("acquiring advisory lock on database")
	db.Raw("SELECT pg_try_advisory_lock(0)").Scan(&lock)
	if !lock {
		log.Printf("failed to acquire lock, skipping migration\n")
		return
	}

	// Ensure we release the lock
	defer db.Raw("SELECT pg_advisory_unlock(0)")
	defer log.Println("removing advisory lock on database")

	// Auto-migrate the schema which should be idempotent
	log.Println("migrating database schema")
	if err := db.AutoMigrate(&Quote{}); err != nil {
		log.Printf("error automigrating schema: %v\n", err)
		return
	}

	// Check for existing data
	var count int64
	db.Model(&Quote{}).Count(&count)
	if count > 0 {
		log.Printf("already have %d records in quote database", count)
		return
	}

	// Ingestion of quote data if we have none.
	log.Println("no quote data found, ingesting new quote data")
	quotes, err := downloadAndParseQuotes()
	if err != nil {
		log.Printf("unable to retrieve quotes: %v\n", err)
		return
	}

	// Not very efficient for large amounts of data but sufficient for now
	for _, quote := range quotes {
		formattedQuote := fmt.Sprintf("%s -- %s", quote.QuoteText, quote.QuoteAuthor)
		db.Create(&Quote{Id: uuid.New().String(), Data: formattedQuote})
	}
	log.Printf("completed ingestion of %d quotes\n", len(quotes))
}

func renderQuote(text string) ([]byte, error) {
	data := struct{ Text string }{text}

	tmpl, err := template.ParseFiles(quoteTemplatePath)
	if err != nil {
		return []byte{}, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (q quoteHandler) handleRequest(w http.ResponseWriter, req *http.Request) {
	req.Body.Close()

	// Prevent caching
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Surrogate-Control", "no-store")

	// Retrieve one quote at random. Doesn't scale to large amounts of data.
	var quote Quote
	q.db.Order("RANDOM()").First(&quote)

	responseBody, err := renderQuote(quote.Data)
	if err != nil {
		log.Printf("error rendering template: %v", err)
		w.WriteHeader(500)
		return
	}

	w.Write(responseBody)
}
