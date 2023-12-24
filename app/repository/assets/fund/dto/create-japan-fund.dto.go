package dto

type CreateJapanFundDto struct {
    Name   string  `json:"name"`
    Code   string  `json:"code"`
    GetPriceTotal float64 `json:"getPriceTotal"`
    GetPrice float64 `json:"getPrice"`
    UserId   string  `json:"userId"`
}