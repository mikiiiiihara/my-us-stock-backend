package totalassets

import (
	"context"
	"fmt"
	"my-us-stock-backend/app/database/model"
	marketPrice "my-us-stock-backend/app/repository/market-price"
)

// 米国株式の評価総額を計算する
func calculateStockTotal(ctx context.Context, ts *DefaultTotalAssetService, modelStocks []model.UsStock) (float64, error) {
	// 株式コードのリストを作成
	// 初期化時に固定長指定→不要なメモリアロケーションが減少
	usStockCodes := make([]string, len(modelStocks))
	for i, modelStock := range modelStocks {
		usStockCodes[i] = modelStock.Code
	}

	// マーケットプライスを取得
	marketPrices, err := ts.MarketPriceRepo.FetchMarketPriceList(ctx, usStockCodes)
	if err != nil {
		return 0, err
	}
	for _, mp := range marketPrices {
		fmt.Println(mp.Ticker, mp.CurrentPrice)
	}

	// マーケットプライスデータをマップに変換
	priceMap := make(map[string]*marketPrice.MarketPriceDto)
	for _, mp := range marketPrices {
		priceMap[mp.Ticker] = &mp
	}
	fmt.Println(priceMap)

	// 現在のドル円を取得
	currentUsdJpy, err := ts.CurrencyRepo.FetchCurrentUsdJpy(ctx)
	if err != nil {
		return 0, err
	}

	// 株式の評価総額を計算
	var amountOfStock = 0.0
	for _, modelStock := range modelStocks {
		// マーケットプライスをマップから取得(O(1) の時間複雑度で取得)
		marketPrice, ok := priceMap[modelStock.Code]
		if !ok {
			// マーケットプライスが見つからない場合はエラーを返す
			return 0, fmt.Errorf("market price not found for stock code: %s", modelStock.Code)
		}

		stockValue := modelStock.Quantity * marketPrice.CurrentPrice * currentUsdJpy
		fmt.Printf("Stock Code: %s, Quantity: %f, Market Price: %f, USD/JPY: %f, Value: %f\n", modelStock.Code, modelStock.Quantity, marketPrice.CurrentPrice, currentUsdJpy, stockValue)
		amountOfStock += stockValue
	}

	return amountOfStock, nil
}
