package dto

type UpdateJapanFundDto struct {
    ID       uint     `json:"id"`
    GetPriceTotal *float64 `json:"getPriceTotal,omitempty"`
    GetPrice *float64 `json:"getPrice,omitempty"`
}
