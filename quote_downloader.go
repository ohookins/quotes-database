package main

import (
	"encoding/json"
	"net/http"
)

const quoteSourceURL = "https://raw.githubusercontent.com/JamesFT/Database-Quotes-JSON/refs/heads/master/quotes.json"

// Structure of above payload
type quoteJSON struct {
	QuoteText   string `json:"quoteText"`
	QuoteAuthor string `json:"quoteAuthor"`
}

// DownloadAndParseQuotes downloads the quotes.json file and parses it into a slice of QuoteJSON
func downloadAndParseQuotes() ([]quoteJSON, error) {
	resp, err := http.Get(quoteSourceURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var quotes []quoteJSON
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&quotes); err != nil {
		return nil, err
	}
	return quotes, nil
}
