package stock

type CreateUsStockDto struct {
    Code   string  `json:"code"`
    GetPrice float64 `json:"getPrice"`
    Quantity float64 `json:"quantity"`
    UserId   uint  `json:"userId"`
    Sector   string  `json:"sector"`
    UsdJpy   float64 `json:"usdjpy"`
}