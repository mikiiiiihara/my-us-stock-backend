package fixedincome

import "time"

type CreateFixedIncomeDto struct {
    Code   string  `json:"code"`
    GetPriceTotal float64 `json:"getPriceTotal"`
    DividendRate float64 `json:"dividendRate"`
    UsdJpy   float64 `json:"usdjpy"`
    PaymentDate time.Time
    UserId   uint  `json:"userId"`
}