package crypto

import (
	"context"
	"my-us-stock-backend/app/common/auth"
	"my-us-stock-backend/app/database/model"
	"my-us-stock-backend/app/graphql/generated"
	"my-us-stock-backend/app/graphql/utils"
	"my-us-stock-backend/app/repository/assets/crypto"
	marketPrice "my-us-stock-backend/app/repository/market-price/crypto"
)

// CryptoService インターフェースの定義
type CryptoService interface {
    Cryptos(ctx context.Context) ([]*generated.Crypto, error)
	CreateCrypto(ctx context.Context, input generated.CreateCryptoInput) (*generated.Crypto, error)
    DeleteCrypto(ctx context.Context, id string) (bool, error)
}

// DefaultCryptoService 構造体の定義
type DefaultCryptoService struct {
    Repo crypto.CryptoRepository // インターフェースを利用
	MarketPriceRepo marketPrice.CryptoRepository
    Auth auth.AuthService        // 認証サービスのインターフェース
}

// NewCryptoService は DefaultUsStockService の新しいインスタンスを作成します
func NewCryptoService(stockRepo crypto.CryptoRepository, auth auth.AuthService, marketPriceRepo marketPrice.CryptoRepository) CryptoService {
    return &DefaultCryptoService{Repo: stockRepo, Auth: auth, MarketPriceRepo: marketPriceRepo}
}

// Cryptos はユーザーの米国株式情報リストを取得します
func (s *DefaultCryptoService) Cryptos(ctx context.Context) ([]*generated.Crypto, error) {
    // アクセストークンの検証（コメントアウトされている部分は必要に応じて実装してください）
    userId, _ := s.Auth.FetchUserIdAccessToken(ctx)
    if userId == 0 {
        return nil, utils.UnauthenticatedError("Invalid user ID")
    }

    modelCryptos, err := s.Repo.FetchCryptoListById(ctx, userId)
    if err != nil {
        return nil, utils.DefaultGraphQLError(err.Error())
    }

    // modelCryptosが空の場合は空の配列を返却する
    if len(modelCryptos) == 0 {
        return []*generated.Crypto{}, nil
    }

    type cryptoWithMarketPrice struct {
        crypto      *model.Crypto
        marketPrice   *marketPrice.Crypto
    }

    // ゴルーチンの実行結果を収集するためのチャネル
    results := make(chan cryptoWithMarketPrice, len(modelCryptos))
    errChan := make(chan error, len(modelCryptos))

    for _, modelCrypto := range modelCryptos {
        modelCryptoCopy := modelCrypto
        go func(ms *model.Crypto) {
            cryptoPrice, err := s.MarketPriceRepo.FetchCryptoPrice(ms.Code)
            if err != nil {
                errChan <- err
                return
            }

            results <- cryptoWithMarketPrice{
                crypto:      ms,
                marketPrice:   cryptoPrice,
            }
        }(&modelCryptoCopy)
    }

    var cryptos []*generated.Crypto
    for i := 0; i < len(modelCryptos); i++ {
        select {
        case result := <-results:
            cryptos = append(cryptos, &generated.Crypto{
                ID: utils.ConvertIdToString(result.crypto.ID),
                Code:         result.crypto.Code,
                GetPrice:     result.crypto.GetPrice,
                Quantity:     result.crypto.Quantity,
                CurrentPrice: result.marketPrice.Price, // 外部APIから取得
            })
        case err := <-errChan:
            return nil, utils.DefaultGraphQLError(err.Error())
        }
    }

    return cryptos, nil
}

// ユーザーの米国株式情報リストを新規作成します
func (s *DefaultCryptoService) CreateCrypto(ctx context.Context, input generated.CreateCryptoInput) (*generated.Crypto, error) {
    // アクセストークンの検証（コメントアウトされている部分は必要に応じて実装してください）
    userId, _ := s.Auth.FetchUserIdAccessToken(ctx)
    if userId == 0 {
        return nil, utils.UnauthenticatedError("Invalid user ID")
    }
	
	// 値入れ直し
	createDto := crypto.CreateCryptDto{
		Code: input.Code,
		GetPrice: input.GetPrice,
		Quantity: input.Quantity,
		UserId: userId,
	}

    modelStock, err := s.Repo.CreateCrypto(ctx, createDto)
    if err != nil {
        return nil, utils.DefaultGraphQLError(err.Error())
    }
    // 市場価格取得
    marketPrice, err := s.MarketPriceRepo.FetchCryptoPrice(createDto.Code)
    if err != nil {
        return nil, utils.DefaultGraphQLError(err.Error())
    }
	// 市場情報を追加して返却
	return &generated.Crypto{
        ID: utils.ConvertIdToString(modelStock.ID),
		Code: modelStock.Code,
		GetPrice: modelStock.GetPrice,
		Quantity:     modelStock.Quantity,
		CurrentPrice: marketPrice.Price,
	}, err
}

// 削除
func (s *DefaultCryptoService) DeleteCrypto(ctx context.Context, id string) (bool, error) {
    // アクセストークンの検証（コメントアウトされている部分は必要に応じて実装してください）
    userId, _ := s.Auth.FetchUserIdAccessToken(ctx)
    if userId == 0 {
        return false, utils.UnauthenticatedError("Invalid user ID")
    }

    // 削除対象id変換
    deleteId, convertError := utils.ConvertIdToUint(id)
    if convertError != nil || deleteId == 0 {
        return false, utils.DefaultGraphQLError("入力されたidが無効です")
       }
    var err = s.Repo.DeleteCrypto(ctx, deleteId)
	// 市場情報を追加して返却
    if err != nil {
     return false, utils.DefaultGraphQLError(err.Error())
    }
	return true, nil
}