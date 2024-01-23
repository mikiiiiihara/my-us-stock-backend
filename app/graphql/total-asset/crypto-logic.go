package totalasset

import (
	"context"
	"my-us-stock-backend/app/database/model"
	"sync"
)

// 仮想通貨の評価総額を計算する
func calculateCryptoTotal(ctx context.Context, ts *DefaultTotalAssetService, modelCryptos []model.Crypto) (float64, error) {
	var amountOfCrypto float64
	var wg sync.WaitGroup
	mu := sync.Mutex{}
	errors := make(chan error, len(modelCryptos))

	for _, modelCrypto := range modelCryptos {
		wg.Add(1)
		go func(mc model.Crypto) {
			defer wg.Done()
			// 現在価格を取得
			cryptoPrice, err := ts.MarketCryptoRepo.FetchCryptoPrice(mc.Code)
			if err != nil {
				errors <- err
				return
			}
			mu.Lock()
			amountOfCrypto += mc.Quantity * cryptoPrice.Price
			mu.Unlock()
		}(modelCrypto)
	}

	wg.Wait()
	close(errors)

	// エラーチェック
	for err := range errors {
		if err != nil {
			return 0, err
		}
	}

	return amountOfCrypto, nil
}
