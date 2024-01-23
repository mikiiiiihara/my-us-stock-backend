package totalasset

import (
	"context"
	"log"
	"math"
	"my-us-stock-backend/app/common/auth"
	"my-us-stock-backend/app/graphql/generated"
	"my-us-stock-backend/app/graphql/utils"
	repoCrypto "my-us-stock-backend/app/repository/assets/crypto"
	repoFixedIncome "my-us-stock-backend/app/repository/assets/fixed-income"
	repoJapanFund "my-us-stock-backend/app/repository/assets/fund"
	"my-us-stock-backend/app/repository/assets/stock"
	marketPrice "my-us-stock-backend/app/repository/market-price"
	repoMarketCrypto "my-us-stock-backend/app/repository/market-price/crypto"
	repoCurrency "my-us-stock-backend/app/repository/market-price/currency"
	repoTotalAsset "my-us-stock-backend/app/repository/total-assets"
	"sort"
	"time"
)

// TotalAssetService インターフェースの定義
type TotalAssetService interface {
    TotalAssets(ctx context.Context, day int) ([]*generated.TotalAsset, error)
	UpdateTotalAsset(ctx context.Context, input generated.UpdateTotalAssetInput) (*generated.TotalAsset, error)
}

// DefaultTotalAssetService 構造体の定義
type DefaultTotalAssetService struct {
    Auth auth.AuthService    // 認証サービスのインターフェース
	TotalAssetRepo repoTotalAsset.TotalAssetRepository
	StockRepo stock.UsStockRepository
	MarketPriceRepo marketPrice.MarketPriceRepository
	CurrencyRepo repoCurrency.CurrencyRepository
	JapanFundRepo repoJapanFund.JapanFundRepository
	CryptoRepo repoCrypto.CryptoRepository
	FixedIncomeRepo repoFixedIncome.FixedIncomeRepository
	MarketCryptoRepo repoMarketCrypto.CryptoRepository
}

// NewTotalAssetService は DefaultUserService の新しいインスタンスを作成します
func NewTotalAssetService(auth auth.AuthService, totalAssetRepo repoTotalAsset.TotalAssetRepository, stockRepo stock.UsStockRepository, marketPriceRepo marketPrice.MarketPriceRepository, currencyRepo repoCurrency.CurrencyRepository, japanFundRepo repoJapanFund.JapanFundRepository,	cryptoRepo repoCrypto.CryptoRepository,fixedIncomeRepo repoFixedIncome.FixedIncomeRepository, marketCryptoRepo repoMarketCrypto.CryptoRepository) TotalAssetService {
	return &DefaultTotalAssetService{auth, totalAssetRepo, stockRepo, marketPriceRepo, currencyRepo, japanFundRepo, cryptoRepo, fixedIncomeRepo, marketCryptoRepo}
}

// GetUserByID はユーザーをIDによって検索します
func (s *DefaultTotalAssetService) TotalAssets(ctx context.Context, day int) ([]*generated.TotalAsset, error) {
    // アクセストークンの検証
    userId, _ := s.Auth.FetchUserIdAccessToken(ctx)
    if userId == 0 {
        return nil, utils.UnauthenticatedError("Invalid user ID")
    }
    modelAssets, err := s.TotalAssetRepo.FetchTotalAssetListById(ctx, userId, day)
    if err != nil {
        return nil, utils.DefaultGraphQLError(err.Error())
    }
	// modelAssetsが空の場合は空の配列を返却する
	if len(modelAssets) == 0 {
		return []*generated.TotalAsset{}, nil
	}

	assets := make([]*generated.TotalAsset, len(modelAssets))
	for i, modelAsset := range modelAssets {

		assets[i] = &generated.TotalAsset{
			ID: utils.ConvertIdToString(modelAsset.ID),
			CashJpy: modelAsset.CashJpy,
			CashUsd: modelAsset.CashUsd,
			Stock: modelAsset.Stock,
			Fund: modelAsset.Fund,
			Crypto: modelAsset.Crypto,
			FixedIncomeAsset: modelAsset.FixedIncomeAsset,
			CreatedAt: modelAsset.CreatedAt.Format(time.RFC3339),
		}
	}

	// assetsをCreatedAtで昇順にソート
	sort.Slice(assets, func(i, j int) bool {
		return assets[i].CreatedAt < assets[j].CreatedAt
	})

    return assets, nil
}

func (s *DefaultTotalAssetService) UpdateTotalAsset(ctx context.Context, input generated.UpdateTotalAssetInput) (*generated.TotalAsset, error) {
    // アクセストークンの検証（コメントアウトされている部分は必要に応じて実装してください）
    userId, _ := s.Auth.FetchUserIdAccessToken(ctx)
    if userId == 0 {
        return nil, utils.UnauthenticatedError("Invalid user ID")
    }
	updateId, convertError := utils.ConvertIdToUint(input.ID)
	if convertError != nil || updateId == 0 {
        return nil, utils.DefaultGraphQLError("入力されたidが無効です")
       }
	   // 株式評価額再計算
	   var amountOfStock = 0.0
	   modelStocks, err := s.StockRepo.FetchUsStockListById(ctx, userId)
	   if err != nil {
		log.Fatalf("エラーが発生しました: %v", err)
		return nil, utils.DefaultGraphQLError("エラーが発生しました")
	   }
		// modelStocksが空の場合は計算処理をスキップする
		if len(modelStocks) != 0 {
			// 米国株の市場価格情報取得
			stockTotal,err := calculateStockTotal(ctx, s, modelStocks)
			if err != nil {
				log.Fatalf("エラーが発生しました: %v", err)
				return nil, utils.DefaultGraphQLError("エラーが発生しました")
			}
			// 資産総額に加算
			amountOfStock += stockTotal
		}
	   // 投資信託評価額再計算
	   var amountOfFund = 0.0
	   modelFunds, err := s.JapanFundRepo.FetchJapanFundListById(ctx, userId)
	   if err != nil {
			log.Fatalf("エラーが発生しました: %v", err)
			return nil, utils.DefaultGraphQLError("エラーが発生しました")
	   }
	   // modelFundsが空の場合は計算処理をスキップする
	   if len(modelFunds) != 0 {
		   // 投資信託の評価総額を計算
		   for _, modelFund := range modelFunds {
			   amountOfFund += calculateFundPriceTotal(modelFund.Code,modelFund.GetPrice,modelFund.GetPriceTotal)
		   }
	   }
	   // 暗号通貨評価額再計算
	   var amountOfCrypto = 0.0
	   modelCryptos, err := s.CryptoRepo.FetchCryptoListById(ctx, userId)
	   if err != nil {
		log.Fatalf("エラーが発生しました: %v", err)
		return nil, utils.DefaultGraphQLError("エラーが発生しました")
	   }
	   // 空の場合は計算処理をスキップする
	   if len(modelCryptos) != 0 {
		   // 仮想通貨の評価総額を計算
		   cryptoTotal,err := calculateCryptoTotal(ctx, s, modelCryptos)
		   if err != nil {
			log.Fatalf("エラーが発生しました: %v", err)
			return nil, utils.DefaultGraphQLError("エラーが発生しました")
		   }
		   // 資産総額に加算
		   amountOfCrypto += cryptoTotal
	   }
	   // 固定利回り資産評価額再計算
	   var amountOfFixedIncomeAsset= 0.0
	   modelAssets, err := s.FixedIncomeRepo.FetchFixedIncomeAssetListById(ctx, userId)
	   if err != nil {
		log.Fatalf("エラーが発生しました: %v", err)
		return nil, utils.DefaultGraphQLError("エラーが発生しました")
	   }
	   // 空の場合は計算処理をスキップする
	   if len(modelAssets) != 0 {
		   // 仮想通貨の評価総額を計算
		   for _, modelAsset := range modelAssets {
			   amountOfFixedIncomeAsset += modelAsset.GetPriceTotal
		   }
	   }
	   roundedAmountOfStock := math.Round(amountOfStock)
	   roundedAmountOfFund := math.Round(amountOfFund)
	   roundedAmountOfCrypto := math.Round(amountOfCrypto)
	   roundedAmountOfFixedIncomeAsset := math.Round(amountOfFixedIncomeAsset)

	   updateDto := repoTotalAsset.UpdateTotalAssetDto{
		   ID: updateId,
		   CashUsd: &input.CashUsd,
		   CashJpy: &input.CashJpy,
		   Stock: &roundedAmountOfStock,
		   Fund: &roundedAmountOfFund,
		   Crypto: &roundedAmountOfCrypto,
		   FixedIncomeAsset: &roundedAmountOfFixedIncomeAsset,
	   }
	   
	updatedAsset, err := s.TotalAssetRepo.UpdateTotalAsset(ctx, updateDto)
    if err != nil {
        return nil, utils.DefaultGraphQLError(err.Error())
    }
	return &generated.TotalAsset{
		ID: input.ID,
		CashJpy: updatedAsset.CashJpy,
		CashUsd: updatedAsset.CashUsd,
		Stock: updatedAsset.Stock,
		Fund: updatedAsset.Fund,
		Crypto: updatedAsset.Crypto,
		FixedIncomeAsset: updatedAsset.FixedIncomeAsset,
		CreatedAt: updatedAsset.CreatedAt.Format(time.RFC3339),
	}, nil
}