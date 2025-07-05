package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Estructura de la reseña que retorna Outscraper
type Review struct {
	ReviewText string `json:"review_text"`
	Rating     string `json:"rating"`
	Date       string `json:"date"`
	Verified   bool   `json:"verified_purchase"`
	Title      string `json:"review_title"`
}

func main() {
	apiKey := "YOUR_API_KEY"
	asin := "B08N5WRWNW" // reemplaza con el ASIN
	region := "com"
	url := fmt.Sprintf("https://api.app.outscraper.com/api/v1/amazon/reviews?query=%s&region=%s", asin, region)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-API-KEY", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var reviews []Review
	if err := json.Unmarshal(body, &reviews); err != nil {
		fmt.Println("❌ Error parsing JSON:", err)
		return
	}

	// Filtrar reseñas del último mes
	lastMonth := time.Now().AddDate(0, -1, 0)
	var recentReviews []Review

	for _, review := range reviews {
		parsedDate, err := time.Parse("Jan 2, 2006", review.Date)
		if err != nil {
			continue
		}
		if parsedDate.After(lastMonth) {
			recentReviews = append(recentReviews, review)
		}
	}

	// Exportar a archivo JSON
	jsonData, err := json.MarshalIndent(recentReviews, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("reviews.json", jsonData, 0644); err != nil {
		panic(err)
	}

	fmt.Printf("✅ %d reseñas guardadas en reviews.json\n", len(recentReviews))
}
