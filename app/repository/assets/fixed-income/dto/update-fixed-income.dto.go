package dto

type UpdateFixedIncomeDto struct {
    ID       uint     `json:"id"`
    GetPriceTotal *float64 `json:"getPriceTotal,omitempty"`
    DividendRate *float64     `json:"dividendRate,omitempty"`
    UsdJpy   *float64 `json:"usdjpy,omitempty"`
}
