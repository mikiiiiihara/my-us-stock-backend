package stock

type UpdateUsStockDto struct {
    ID       uint     `json:"id"`
    GetPrice *float64 `json:"getPrice,omitempty"`
    Quantity *float64     `json:"quantity,omitempty"`
    UsdJpy   *float64 `json:"usdjpy,omitempty"`
}
