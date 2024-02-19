package fund

type UpdateFundPriceDto struct {
    ID       uint     `json:"id"`
    Price float64 `json:"price"`
}