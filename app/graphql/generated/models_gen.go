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

type UsStock struct {
	ID string `json:"id"`
	// ティッカーシンボル
	Code string `json:"code"`
	// 取得価格
	GetPrice float64 `json:"getPrice"`
	// １年当たり配当
	Dividend float64 `json:"dividend"`
	// 保有株数
	Quantity float64 `json:"quantity"`
	// セクター
	Sector string `json:"sector"`
	// 購入時為替
	UsdJpy float64 `json:"usdJpy"`
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
