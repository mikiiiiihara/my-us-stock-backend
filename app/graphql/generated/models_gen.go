// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package generated

type CreateCryptoInput struct {
	// ティッカーシンボル
	Code string `json:"code"`
	// 取得価格
	GetPrice float64 `json:"getPrice"`
	// 保有株数
	Quantity float64 `json:"quantity"`
}

type CreateFixedIncomeAssetInput struct {
	// 資産名称
	Code string `json:"code"`
	// 取得価格合計
	GetPriceTotal float64 `json:"getPriceTotal"`
	// １年当たり配当利回り
	DividendRate float64 `json:"dividendRate"`
	// 購入時為替
	UsdJpy *float64 `json:"usdJpy,omitempty"`
	// 配当支払い月
	PaymentMonth []int `json:"paymentMonth"`
}

type CreateJapanFundInput struct {
	// ティッカーシンボル
	Code string `json:"code"`
	// 銘柄名
	Name string `json:"name"`
	// 取得価格
	GetPrice float64 `json:"getPrice"`
	// 取得価格総額
	GetPriceTotal float64 `json:"getPriceTotal"`
}

type CreateUsStockInput struct {
	// ティッカーシンボル
	Code string `json:"code"`
	// 取得価格
	GetPrice float64 `json:"getPrice"`
	// 保有株数
	Quantity float64 `json:"quantity"`
	// セクター
	Sector string `json:"sector"`
	// 購入時為替
	UsdJpy float64 `json:"usdJpy"`
}

type CreateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Crypto struct {
	ID string `json:"id"`
	// ティッカーシンボル
	Code string `json:"code"`
	// 取得価格
	GetPrice float64 `json:"getPrice"`
	// 保有株数
	Quantity float64 `json:"quantity"`
	// 現在価格
	CurrentPrice float64 `json:"currentPrice"`
}

type FixedIncomeAsset struct {
	ID string `json:"id"`
	// 資産名称
	Code string `json:"code"`
	// 取得価格合計
	GetPriceTotal float64 `json:"getPriceTotal"`
	// １年当たり配当利回り
	DividendRate float64 `json:"dividendRate"`
	// 購入時為替
	UsdJpy *float64 `json:"usdJpy,omitempty"`
	// 配当支払い月
	PaymentMonth []int `json:"paymentMonth"`
}

type JapanFund struct {
	ID string `json:"id"`
	// ティッカーシンボル
	Code string `json:"code"`
	// 銘柄名
	Name string `json:"name"`
	// 取得価格
	GetPrice float64 `json:"getPrice"`
	// 取得価格総額
	GetPriceTotal float64 `json:"getPriceTotal"`
	// 現在価格
	CurrentPrice float64 `json:"currentPrice"`
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
