package crypto

type CreateCryptDto struct {
    Code   string  `json:"code"`
    GetPrice float64 `json:"getPrice"`
    Quantity float64 `json:"quantity"`
    UserId   uint  `json:"userId"`
}