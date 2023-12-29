package crypto

import (
	"context"
	"my-us-stock-backend/app/common/auth"
	"my-us-stock-backend/app/database/model"
	"my-us-stock-backend/app/graphql/generated"
	"my-us-stock-backend/app/graphql/utils"
	"my-us-stock-backend/app/repository/assets/stock"
	marketPrice "my-us-stock-backend/app/repository/market-price/crypto"
)

// CryptoService インターフェースの定義
type CryptoService interface {
    Cryptos(ctx context.Context) ([]*generated.UsStock, error)
	CreateCrypto(ctx context.Context, input generated.CreateUsStockInput) (*generated.UsStock, error)
}

// DefaultCryptoService 構造体の定義
type DefaultCryptoService struct {
    StockRepo stock.UsStockRepository // インターフェースを利用
	MarketPriceRepo marketPrice.CryptoRepository
    Auth auth.AuthService        // 認証サービスのインターフェース
}

// NewUsCryptoService は DefaultUsStockService の新しいインスタンスを作成します
func NewUsCryptoService(stockRepo stock.UsStockRepository, auth auth.AuthService, marketPriceRepo marketPrice.CryptoRepository) CryptoService {
    return &DefaultCryptoService{StockRepo: stockRepo, Auth: auth, MarketPriceRepo: marketPriceRepo}
}

// Cryptos はユーザーの米国株式情報リストを取得します
func (s *DefaultCryptoService) Cryptos(ctx context.Context) ([]*generated.UsStock, error) {
    // アクセストークンの検証（コメントアウトされている部分は必要に応じて実装してください）
    userId, _ := s.Auth.FetchUserIdAccessToken(ctx)
    if userId == 0 {
        return nil, utils.UnauthenticatedError("Invalid user ID")
    }

    modelCryptos, err := s.StockRepo.FetchUsStockListById(ctx, userId)
    if err != nil {
        return nil, err
    }

    type cryptoWithMarketPrice struct {
        crypto      *model.Crypto
        dividend   *marketPrice.Crypto
    }

    // ゴルーチンの実行結果を収集するためのチャネル
    results := make(chan cryptoWithMarketPrice, len(modelCryptos))
    errChan := make(chan error, len(modelCryptos))

    for _, modelCrypto := range modelCryptos {
        modelCryptoCopy := modelCrypto
        go func(ms *model.Crypto) {
            dividend, err := s.MarketPriceRepo.FetchCryptoPrice(ms.Code)
            if err != nil {
                errChan <- err
                return
            }

            results <- cryptoWithMarketPrice{
                crypto:      ms,
                dividend:   dividend,
            }
        }(&modelCryptoCopy)
    }

    var cryptos []*generated.UsStock
    for i := 0; i < len(modelCryptos); i++ {
        select {
        case result := <-results:
            cryptos = append(cryptos, &generated.UsStock{
                ID: utils.ConvertIdToString(result.crypto.ID),
                Code:         result.crypto.Code,
                GetPrice:     result.crypto.GetPrice,
                Dividend:     result.dividend.DividendTotal,
                Quantity:     result.crypto.Quantity,
                Sector:       result.crypto.Sector,
                UsdJpy:       result.crypto.UsdJpy,
                // CurrentPrice: result.marketPrice.CurrentPrice, 
                // PriceGets:    result.marketPrice.PriceGets,
                // CurrentRate:  result.marketPrice.CurrentRate,
            })
        case err := <-errChan:
            return nil, err
        }
    }

    return cryptos, nil
}

// UsStocks はユーザーの米国株式情報リストを新規作成します
func (s *DefaultCryptoService) CreateCrypto(ctx context.Context, input generated.CreateUsStockInput) (*generated.UsStock, error) {
    // アクセストークンの検証（コメントアウトされている部分は必要に応じて実装してください）
    userId, _ := s.Auth.FetchUserIdAccessToken(ctx)
    if userId == 0 {
        return nil, utils.UnauthenticatedError("Invalid user ID")
    }
	
	// 値入れ直し
	createDto := stock.CreateUsStockDto{
		Code: input.Code,
		GetPrice: input.GetPrice,
		Quantity: input.Quantity,
		UsdJpy: input.UsdJpy,
		Sector: input.Sector,
		UserId: userId,
	}

    modelStock, err := s.StockRepo.CreateUsStock(ctx, createDto)
    if err != nil {
        return nil, err
    }
    // 配当情報取得
    dividend, err := s.MarketPriceRepo.FetchDividend(ctx, createDto.Code)
    if err != nil {
        return nil, err
    }
    // 市場価格取得
    var codeInput = []string{createDto.Code}
    marketPrices, err := s.MarketPriceRepo.FetchMarketPriceList(ctx, codeInput)
    if err != nil {
        return nil, err
    }
	// 市場情報を追加して返却
	return &generated.UsStock{
        ID: utils.ConvertIdToString(modelStock.ID),
		Code: modelStock.Code,
		GetPrice: modelStock.GetPrice,
		Dividend: dividend.DividendTotal,
		Quantity:     modelStock.Quantity,
		Sector:       modelStock.Sector,
		UsdJpy:       modelStock.UsdJpy,
		CurrentPrice: marketPrices[0].CurrentPrice,
		PriceGets:    marketPrices[0].PriceGets,
		CurrentRate:  marketPrices[0].CurrentRate,
	}, err
}