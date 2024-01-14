package stock

import (
	"context"
	"my-us-stock-backend/app/common/auth"
	"my-us-stock-backend/app/database/model"
	"my-us-stock-backend/app/graphql/generated"
	"my-us-stock-backend/app/graphql/utils"
	"my-us-stock-backend/app/repository/assets/stock"
	marketPrice "my-us-stock-backend/app/repository/market-price"
)

// UsStockService インターフェースの定義
type UsStockService interface {
    UsStocks(ctx context.Context) ([]*generated.UsStock, error)
	CreateUsStock(ctx context.Context, input generated.CreateUsStockInput) (*generated.UsStock, error)
}

// DefaultUsStockService 構造体の定義
type DefaultUsStockService struct {
    StockRepo stock.UsStockRepository // インターフェースを利用
	MarketPriceRepo marketPrice.MarketPriceRepository
    Auth auth.AuthService        // 認証サービスのインターフェース
}

// NewUsStockService は DefaultUsStockService の新しいインスタンスを作成します
func NewUsStockService(stockRepo stock.UsStockRepository, auth auth.AuthService, marketPriceRepo marketPrice.MarketPriceRepository) UsStockService {
    return &DefaultUsStockService{StockRepo: stockRepo, Auth: auth, MarketPriceRepo: marketPriceRepo}
}

// UsStocks はユーザーの米国株式情報リストを取得します
func (s *DefaultUsStockService) UsStocks(ctx context.Context) ([]*generated.UsStock, error) {
    // アクセストークンの検証（コメントアウトされている部分は必要に応じて実装してください）
    userId, _ := s.Auth.FetchUserIdAccessToken(ctx)
    if userId == 0 {
        return nil, utils.UnauthenticatedError("Invalid user ID")
    }

    modelStocks, err := s.StockRepo.FetchUsStockListById(ctx, userId)
    if err != nil {
        return nil, utils.DefaultGraphQLError(err.Error())
    }
    // modelStocksが空の場合は空の配列を返却する
	if len(modelStocks) == 0 {
		return []*generated.UsStock{}, nil
	}
    // 米国株の市場価格情報取得
    // (本来はfor文内で呼びたいが、外部APIコール数削減のため一度に呼んでいる)
    usStockCodes := make([]string, len(modelStocks))
    for i, modelStock := range modelStocks {
        usStockCodes[i] = modelStock.Code
    }
    marketPrices, err := s.MarketPriceRepo.FetchMarketPriceList(ctx,usStockCodes)
    if err != nil {
        return nil, utils.DefaultGraphQLError(err.Error())
    }
    type stockWithMarketPrice struct {
        stock      *model.UsStock
        marketPrice *marketPrice.MarketPriceDto
        dividend   *marketPrice.DividendEntity
    }

    // ゴルーチンの実行結果を収集するためのチャネル
    results := make(chan stockWithMarketPrice, len(modelStocks))
    errChan := make(chan error, len(modelStocks))

    for _, modelStock := range modelStocks {
        modelStockCopy := modelStock
        go func(ms *model.UsStock) {
            dividend, err := s.MarketPriceRepo.FetchDividend(ctx, ms.Code)
            if err != nil {
                errChan <- utils.DefaultGraphQLError(err.Error())
                return
            }

            var marketPrice *marketPrice.MarketPriceDto
            for _, mp := range marketPrices {
                if mp.Ticker == ms.Code {
                    marketPrice = &mp
                    break
                }
            }

            results <- stockWithMarketPrice{
                stock:      ms,
                marketPrice: marketPrice,
                dividend:   dividend,
            }
        }(&modelStockCopy)
    }

    usStocks := make([]*generated.UsStock, len(modelStocks))
    for i := 0; i < len(modelStocks); i++ {
        select {
        case result := <-results:
    
            // marketPriceがnilかどうかをチェック
            var currentPrice, priceGets, currentRate float64
            if result.marketPrice != nil {
                currentPrice = result.marketPrice.CurrentPrice
                priceGets = result.marketPrice.PriceGets
                currentRate = result.marketPrice.CurrentRate
            } else {
                currentPrice = result.stock.GetPrice
                priceGets = 0.0
                currentRate = 0.0
            }
    
            usStocks[i] = &generated.UsStock{
                ID:           utils.ConvertIdToString(result.stock.ID),
                Code:         result.stock.Code,
                GetPrice:     result.stock.GetPrice,
                Dividend:     result.dividend.DividendTotal,
                Quantity:     result.stock.Quantity,
                Sector:       result.stock.Sector,
                UsdJpy:       result.stock.UsdJpy,
                CurrentPrice: currentPrice,
                PriceGets:    priceGets,
                CurrentRate:  currentRate,
            }
        case err := <-errChan:
            return nil, utils.DefaultGraphQLError(err.Error())
        }
    }
    

    return usStocks, nil
}

// UsStocks はユーザーの米国株式情報リストを新規作成します
func (s *DefaultUsStockService) CreateUsStock(ctx context.Context, input generated.CreateUsStockInput) (*generated.UsStock, error) {
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
    // すでに登録されていますの時はキャッチして、repo.updateを呼び出す
    if err != nil {
        return nil, utils.DefaultGraphQLError(err.Error())
    }
    // 配当情報取得
    dividend, err := s.MarketPriceRepo.FetchDividend(ctx, createDto.Code)
    if err != nil {
        return nil, utils.DefaultGraphQLError(err.Error())
    }
    // 市場価格取得
    var codeInput = []string{createDto.Code}
    marketPrices, err := s.MarketPriceRepo.FetchMarketPriceList(ctx, codeInput)
    if err != nil {
        return nil, utils.DefaultGraphQLError(err.Error())
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