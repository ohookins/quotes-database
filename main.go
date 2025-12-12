package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s, %d bytes served\n", r.Method, r.RequestURI, r.Response.Status, r.Response.ContentLength)
	})
}

func main() {
	// Set up database from environment, or default to local database for testing.
	var dsn string
	var ok bool
	if dsn, ok = os.LookupEnv("DSN"); !ok {
		dsn = "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// basic healthcheck for apprunner, don't log requests as there will be a LOT
	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "OK")
	})

	qh := newQuoteHandler(db)

	http.Handle("/", logRequest(http.HandlerFunc(qh.handleRequest)))
	err = http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
