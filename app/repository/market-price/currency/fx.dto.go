package currency

// CurrencyPair は通貨ペアの情報を表します。
type CurrencyPair struct {
    High             string  `json:"high"`
    Open             string  `json:"open"`
    Bid              string `json:"bid"`
    CurrencyPairCode string  `json:"currencyPairCode"`
    Ask              string  `json:"ask"`
    Low              string  `json:"low"`
}

// Fx は通貨のレート情報を表します。
type Fx struct {
    Quotes []CurrencyPair `json:"quotes"`
}
