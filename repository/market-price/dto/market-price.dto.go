package dto

type MarketPriceDto struct {
    Ticker        string  `json:"ticker"`
    CurrentPrice  float64 `json:"currentPrice"`
    PriceGets  float64 `json:"priceGets"`
    CurrentRate float64 `json:"currentRate"`
}