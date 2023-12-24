package dto

type UpdateCryptoDto struct {
    ID       uint     `json:"id"`
    GetPrice *float64 `json:"getPrice,omitempty"`
    Quantity *float64     `json:"quantity,omitempty"`
}
