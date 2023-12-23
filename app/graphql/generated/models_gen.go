// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package generated

type CreateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type MarketPrice struct {
	// ティッカーシンボル
	Ticker string `json:"ticker"`
	// 現在価格
	CurrentPrice float64 `json:"currentPrice"`
	// 変化額
	PriceGets float64 `json:"priceGets"`
	// 変化率
	CurrentRate float64 `json:"currentRate"`
}

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
