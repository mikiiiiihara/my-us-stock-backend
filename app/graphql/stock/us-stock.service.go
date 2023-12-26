package stock

import (
	"context"
	"my-us-stock-backend/app/common/auth"
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
        return nil, err
    }

    // 空の配列を初期化
    usStocks := make([]*generated.UsStock, 0)
    for _, modelStock := range modelStocks {
        usStock := &generated.UsStock{
            ID: utils.ConvertIdToString(modelStock.ID),
            Code:         modelStock.Code,
            GetPrice:     modelStock.GetPrice,
            Dividend:     99.9, // TODO: 配当APIから取得
            Quantity:     modelStock.Quantity,
            Sector:       modelStock.Sector,
            UsdJpy:       modelStock.UsdJpy,
            CurrentPrice: 99.9, // TODO: 株価APIから取得
            PriceGets:    1.0,
            CurrentRate:  1.1,
        }
        usStocks = append(usStocks, usStock)
    }
    // エラーがなければ空の配列とnilを返す
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
    if err != nil {
        return nil, err
    }
	// 市場情報を追加して返却
	return &generated.UsStock{
        ID: utils.ConvertIdToString(modelStock.ID),
		Code: modelStock.Code,
		GetPrice: modelStock.GetPrice,
		Dividend: 99.9, // TODO: 配当APIから取得
		Quantity:     modelStock.Quantity,
		Sector:       modelStock.Sector,
		UsdJpy:       modelStock.UsdJpy,
		CurrentPrice: 99.9, // TODO: 株価APIから取得
		PriceGets:    1.0,
		CurrentRate:  1.1,
	}, err
}