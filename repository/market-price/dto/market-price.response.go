package dto

type MarketPriceResponse struct {
	Symbol               string  `json:"symbol"`               // ティッカー名
	Name                 string  `json:"name"`
	Price                float64 `json:"price"`                // 現実価格
	ChangesPercentage    float64 `json:"changesPercentage"`    // 騰落率
	Change               float64 `json:"change"`               // 変化額
	DayLow               float64 `json:"dayLow"`
	DayHigh              float64 `json:"dayHigh"`
	YearHigh             float64 `json:"yearHigh"`
	YearLow              float64 `json:"yearLow"`
	MarketCap            float64 `json:"marketCap"`
	PriceAvg50           float64 `json:"priceAvg50"`
	PriceAvg200          float64 `json:"priceAvg200"`
	Exchange             string  `json:"exchange"`
	Volume               float64 `json:"volume"`
	AvgVolume            float64 `json:"avgVolume"`
	Open                 float64 `json:"open"`
	PreviousClose        float64 `json:"previousClose"`
	EPS                  float64 `json:"eps"`
	PE                   float64 `json:"pe"`
	EarningsAnnouncement string  `json:"earningsAnnouncement"`
	SharesOutstanding    float64 `json:"sharesOutstanding"`
	Timestamp            int64   `json:"timestamp"`
}