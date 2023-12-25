package marketprice

// Historical は配当情報の履歴を表します。
type Historical struct {
	Date           string  `json:"date"`           // 権利落月
	Label          string  `json:"label"`
	AdjDividend    float64 `json:"adjDividend"`
	Dividend       float64 `json:"dividend"`       // 配当支払額
	RecordDate     string  `json:"recordDate"`
	PaymentDate    string  `json:"paymentDate"`    // 配当支払月
	DeclarationDate string `json:"declarationDate"`
}

// DividendResponse はAPIからの配当情報レスポンスを表します。
type DividendResponse struct {
	Historical []Historical `json:"historical"`
	Symbol     string       `json:"symbol"`
}
