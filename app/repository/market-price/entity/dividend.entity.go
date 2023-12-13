package entity

type DividendEntity struct {
    Ticker          string    `json:"ticker"`
    DividendTime    int       `json:"dividendTime"`
    DividendMonth   []int     `json:"dividendMonth"`
    DividendFixedMonth []int  `json:"dividendFixedMonth"`
    Dividend        float64   `json:"dividend"`
    DividendTotal   float64   `json:"dividendTotal"`
}