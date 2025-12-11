package main

import (
	"fmt"
	"io"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "OK")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		var count int64
		db.Model(&Quote{}).Count(&count)
		io.WriteString(w, fmt.Sprintf("%d records\n", count))
	})

	http.ListenAndServe(":8080", nil)
}
