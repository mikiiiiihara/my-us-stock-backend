package dto

type CreateUsStockDto struct {
    Code   string  `json:"code"`
    GetPrice float64 `json:"getPrice"`
    Quantity float64 `json:"quantity"`
    UserId   string  `json:"userId"`
    Sector   string  `json:"sector"`
    UsdJpy   float64 `json:"usdjpy"`
}