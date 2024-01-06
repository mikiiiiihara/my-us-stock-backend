package totalassets

// 日本投資信託の評価総額を計算する
func calculateFundPriceTotal(code string, getPrice float64, getPriceTotal float64) float64 {
	return getPriceTotal*getFundMarketPrice(code)/getPrice
}

// TODO: 三菱UFJのAPI復旧したらrepository実装してそこから取得するようにする
func getFundMarketPrice(code string) float64 {
	switch code {
	case "SP500":
		return 24342.0
	case "全世界株":
		return 20926.0
	default:
		return 18584.0
	
	}
}