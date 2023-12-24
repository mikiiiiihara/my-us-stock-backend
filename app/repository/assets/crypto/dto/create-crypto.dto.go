package dto

type CreateCryptDto struct {
    Code   string  `json:"code"`
    GetPrice float64 `json:"getPrice"`
    Quantity float64 `json:"quantity"`
    UserId   string  `json:"userId"`
}