package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"time"

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
		log.Printf("failed to acquire log, skipping migration")
		return
	}

	// Auto-migrate the schema which should be idempotent
	log.Println("migrating database schema")
	db.AutoMigrate(&Quote{})

	// Ingest data here
	log.Println("ingesting quote data")
	time.Sleep(5 * time.Second)

	// Release lock
	log.Println("removing advisory lock on database")
	db.Raw("SELECT pg_advisory_unlock(0)")
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

	// var count int64
	// q.db.Model(&Quote{}).Count(&count)
	// io.WriteString(w, fmt.Sprintf("%d records\n", count))

	// Placeholder
	quote := `I used to work in a fire hydrant factory.  You couldn't park anywhere near
the place.
		-- Steven Wright
`
	responseBody, err := renderQuote(quote)
	if err != nil {
		log.Printf("error rendering template: %v", err)
		w.WriteHeader(500)
		return
	}

	w.Write(responseBody)
}
