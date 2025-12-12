package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type quoteHandler struct {
	db *gorm.DB
}

func newQuoteHandler(db *gorm.DB) quoteHandler {
	migrate(db)
	return quoteHandler{db: db}
}

func migrate(db *gorm.DB) {
	lock := struct{ val bool }{}

	// Take out advisory lock in the database to prevent simultaneous migrations.
	log.Println("acquiring advisory lock on database")
	db.Raw("SELECT pg_try_advisory_lock(0)").Scan(&lock)
	if !lock.val {
		log.Printf("failed to acquire log, skipping migration")
		return
	}

	// Auto-migrate the schema which should be idempotent
	db.AutoMigrate(&Quote{})

	// Ingest data here
	time.Sleep(5 * time.Second)

	// Release lock
	db.Raw("SELECT pg_advisory_unlock(0)")
}

func (q quoteHandler) handleRequest(w http.ResponseWriter, req *http.Request) {
	var count int64
	q.db.Model(&Quote{}).Count(&count)
	io.WriteString(w, fmt.Sprintf("%d records\n", count))
}
