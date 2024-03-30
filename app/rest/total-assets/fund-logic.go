package totalassets

import (
	"context"
	"my-us-stock-backend/app/database/model"
	"sync"
)

// 日本投資信託の評価総額を計算する
func calculateFundPriceTotal(ctx context.Context, ts *DefaultTotalAssetService, modelFunds []model.JapanFund) (float64, error) {
	var amountOfFund float64
	var wg sync.WaitGroup
	mu := sync.Mutex{}
	errors := make(chan error, len(modelFunds))

	for _, modelFund := range modelFunds {
		wg.Add(1)
		go func(mf model.JapanFund) {
			defer wg.Done()
			// 現在価格を取得
			fundPrice, err := ts.FundPriceRepo.FindFundPriceByCode(ctx, mf.Code)
			if err != nil {
				errors <- err
				return
			}
			mu.Lock()
			amountOfFund += mf.GetPriceTotal * fundPrice.Price / mf.GetPrice
			mu.Unlock()
		}(modelFund)
	}

	wg.Wait()
	close(errors)

	// エラーチェック
	for err := range errors {
		if err != nil {
			return 0, err
		}
	}

	return amountOfFund, nil
}