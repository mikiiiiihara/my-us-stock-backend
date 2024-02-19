package fund

type CreateFundPriceDto struct {
    Name   string  `json:"name"`
    Code   string  `json:"code"`
    Price float64 `json:"price"`
}