package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// QuoteData は株価データを表す構造体です。
type QuoteData struct {
    Symbol               string  `json:"symbol"`
    Name                 string  `json:"name"`
    Price                float64 `json:"price"`
    ChangesPercentage    float64 `json:"changesPercentage"`
    Change               float64 `json:"change"`
    DayLow               float64 `json:"dayLow"`
    DayHigh              float64 `json:"dayHigh"`
    YearHigh             float64 `json:"yearHigh"`
    YearLow              float64 `json:"yearLow"`
    MarketCap            int64   `json:"marketCap"`
    PriceAvg50           float64 `json:"priceAvg50"`
    PriceAvg200          float64 `json:"priceAvg200"`
    Exchange             string  `json:"exchange"`
    Volume               int     `json:"volume"`
    AvgVolume            int     `json:"avgVolume"`
    Open                 float64 `json:"open"`
    PreviousClose        float64 `json:"previousClose"`
    Eps                  float64 `json:"eps"`
    Pe                   float64 `json:"pe"`
    EarningsAnnouncement string  `json:"earningsAnnouncement"`
    SharesOutstanding    int64   `json:"sharesOutstanding"`
    Timestamp            int64   `json:"timestamp"`
}

// HistoricalData は履歴データを表す構造体です。
type HistoricalData struct {
	Date            string  `json:"date"`
	Label           string  `json:"label"`
	AdjDividend     float64 `json:"adjDividend"`
	Dividend        float64 `json:"dividend"`
	RecordDate      string  `json:"recordDate"`
	PaymentDate     string  `json:"paymentDate"`
	DeclarationDate string  `json:"declarationDate"`
}

// Response はレスポンスデータを表す構造体です。
type Response struct {
	Symbol      string           `json:"symbol"`
	Historical  []HistoricalData `json:"historical"`
}

func quoteOrderHandler(w http.ResponseWriter, r *http.Request) {
    // 固定のレスポンスデータを設定
    responseData := []QuoteData{
		{"AAPL","Apple Inc.",192.53, -0.5424,-1.05,191.725,194.39,199.62,124.17,2994380584000,186.3,179.2902,"NASDAQ", 41936044,53103488,193.9,193.58,6.12,31.46,"2024-01-31T10:59:00.000+0000",15552800000,1703883601},
	}

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(responseData)
}

func historicalPriceHandler(w http.ResponseWriter, r *http.Request) {
	// 固定のレスポンスデータを設定
	response := Response{
		Symbol: "AAPL",
		Historical: []HistoricalData{
			{"2023-11-10", "November 10, 23", 0.24, 0.24, "2023-11-13", "2023-11-16", "2023-11-02"},
			{"2023-08-11", "August 11, 23", 0.24, 0.24, "2023-08-14", "2023-08-17", "2023-08-03"},
			{"2023-05-12", "May 12, 23", 0.24, 0.24, "2023-05-15", "2023-05-18", "2023-05-04"},
			{"2023-02-10", "February 10, 23", 0.24, 0.24, "2023-02-13", "2023-02-16", "2023-02-02"},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/api/v3/quote-order/", quoteOrderHandler)
	http.HandleFunc("/api/v3/historical-price-full/stock_dividend/", historicalPriceHandler)

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
