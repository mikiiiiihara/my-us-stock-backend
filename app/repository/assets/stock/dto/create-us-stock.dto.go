package dto

type CreateUsStockDto struct {
    Ticker   string  `json:"ticker"`
    GetPrice float64 `json:"getPrice"`
    Quantity float64 `json:"quantity"`
    UserId   string  `json:"userId"`
    Sector   string  `json:"sector"`
    UsdJpy   float64 `json:"usdjpy"`
}