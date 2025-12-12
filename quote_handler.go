package main

import (
	"fmt"
	"io"
	"net/http"

	"gorm.io/gorm"
)

type quoteHandler struct {
	db *gorm.DB
}

func newQuoteHandler(db *gorm.DB) quoteHandler {
	q := quoteHandler{db: db}

	// Take out advisory lock in the database to prevent simultaneous migrations.

	// Auto-migrate the schema which should be idempotent
	db.AutoMigrate(&Quote{})

	return q
}

func (q quoteHandler) handleRequest(w http.ResponseWriter, req *http.Request) {
	var count int64
	q.db.Model(&Quote{}).Count(&count)
	io.WriteString(w, fmt.Sprintf("%d records\n", count))
}
