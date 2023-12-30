package fixedincome

type CreateFixedIncomeDto struct {
    Code   string  `json:"code"`
    GetPriceTotal float64 `json:"getPriceTotal"`
    DividendRate float64 `json:"dividendRate"`
    UsdJpy   float64 `json:"usdjpy"`
    PaymentMonth []int64 `json:"paymentMonth"`
    UserId   uint  `json:"userId"`
}