package totalassets

import (
	"context"
	"math"

	repoCrypto "my-us-stock-backend/app/repository/assets/crypto"
	repoFixedIncome "my-us-stock-backend/app/repository/assets/fixed-income"
	repoJapanFund "my-us-stock-backend/app/repository/assets/fund"
	"my-us-stock-backend/app/repository/assets/stock"
	marketPrice "my-us-stock-backend/app/repository/market-price"
	repoMarketCrypto "my-us-stock-backend/app/repository/market-price/crypto"
	repoCurrency "my-us-stock-backend/app/repository/market-price/currency"
	repoTotalAsset "my-us-stock-backend/app/repository/total-assets"

	"github.com/gin-gonic/gin"
)

// TotalAssetService インターフェースの定義
type TotalAssetService interface {
	CreateTodayTotalAsset(ctx context.Context, c *gin.Context) (string, error)
}

// DefaultTotalAssetService 構造体の定義
type DefaultTotalAssetService struct {
	TotalAssetRepo repoTotalAsset.TotalAssetRepository
	StockRepo stock.UsStockRepository // インターフェースを利用
	MarketPriceRepo marketPrice.MarketPriceRepository
	CurrencyRepo repoCurrency.CurrencyRepository
	JapanFundRepo repoJapanFund.JapanFundRepository
	CryptoRepo repoCrypto.CryptoRepository
	FixedIncomeRepo repoFixedIncome.FixedIncomeRepository
	MarketCryptoRepo repoMarketCrypto.CryptoRepository
}

// DefaultTotalAssetService の新しいインスタンスを作成します
func NewTotalAssetService(totalAssetRepo repoTotalAsset.TotalAssetRepository, stockRepo stock.UsStockRepository, marketPriceRepo marketPrice.MarketPriceRepository, currencyRepo repoCurrency.CurrencyRepository, japanFundRepo repoJapanFund.JapanFundRepository,	cryptoRepo repoCrypto.CryptoRepository,fixedIncomeRepo repoFixedIncome.FixedIncomeRepository, marketCryptoRepo repoMarketCrypto.CryptoRepository) TotalAssetService {
	return &DefaultTotalAssetService{totalAssetRepo, stockRepo, marketPriceRepo, currencyRepo, japanFundRepo, cryptoRepo, fixedIncomeRepo, marketCryptoRepo}
}

// 資産新規登録処理
func (ts *DefaultTotalAssetService) CreateTodayTotalAsset(ctx context.Context, c *gin.Context) (string, error) {
    var requestParam CreateTotalAssetRequest
    if err := c.BindJSON(&requestParam); err != nil {
        // JSONパースエラーが発生した場合、400 Bad Requestを返す
        return "Bad Request", err
    }

	// 当日分の資産が登録されているか確認
	latestTotalAsset, err  := ts.TotalAssetRepo.FindTodayTotalAsset(ctx, uint(requestParam.UserId))
	if err == nil && latestTotalAsset != nil {
        return "すでに資産が登録されています。", nil
    }
	// 保有株式を取得
	var amountOfStock = 0.0
	modelStocks, err := ts.StockRepo.FetchUsStockListById(ctx, uint(requestParam.UserId))
	if err != nil {
        return "Internal Server Error", err
    }
    // modelStocksが空の場合は計算処理をスキップする
	if len(modelStocks) != 0 {
		// 米国株の市場価格情報取得
		// (本来はfor文内で呼びたいが、外部APIコール数削減のため一度に呼んでいる)
		usStockCodes := make([]string, 0)
		for _, modelStock := range modelStocks {
			usStockCodes = append(usStockCodes, modelStock.Code)
		}
		marketPrices, err := ts.MarketPriceRepo.FetchMarketPriceList(ctx,usStockCodes)
		if err != nil {
			return "Internal Server Error", err
		}
		// 現在のドル円を取得
		currentUsdJpy, err := ts.CurrencyRepo.FetchCurrentUsdJpy(ctx)
		if err != nil {
			return "Internal Server Error", err
		}
		// 株式の評価総額を計算
		for _, modelStock := range modelStocks {
			var marketPrice *marketPrice.MarketPriceDto
            for _, mp := range marketPrices {
                if mp.Ticker == modelStock.Code {
                    marketPrice = &mp
                    break
                }
            }
			// 株式評価総額に加算
			amountOfStock += modelStock.Quantity * marketPrice.CurrentPrice*currentUsdJpy
		}
	}
	// 日本投資信託の評価額を取得
	var amountOfFund = 0.0
	modelFunds, err := ts.JapanFundRepo.FetchJapanFundListById(ctx, uint(requestParam.UserId))
	if err != nil {
        return "Internal Server Error", err
    }
	// modelFundsが空の場合は計算処理をスキップする
	if len(modelFunds) != 0 {
		// 投資信託の評価総額を計算
		for _, modelFund := range modelFunds {
			amountOfFund += calculateFundPriceTotal(modelFund.Code,modelFund.GetPrice,modelFund.GetPriceTotal)
		}
	}

	// 仮想通貨の評価額を取得
	var amountOfCrypto = 0.0
	modelCryptos, err := ts.CryptoRepo.FetchCryptoListById(ctx, uint(requestParam.UserId))
	if err != nil {
        return "Internal Server Error", err
    }
	// 空の場合は計算処理をスキップする
	if len(modelCryptos) != 0 {
		// 仮想通貨の評価総額を計算
		for _, modelCrypto := range modelCryptos {
			// 現在価格を取得
			cryptoPrice, err := ts.MarketCryptoRepo.FetchCryptoPrice(modelCrypto.Code)
			if err != nil {
				return "Internal Server Error", err
			}
			amountOfCrypto += modelCrypto.Quantity*cryptoPrice.Price
		}
	}

	// 固定利回り資産の評価額を取得
	var amountOfFixedIncomeAsset= 0.0
	modelAssets, err := ts.FixedIncomeRepo.FetchFixedIncomeAssetListById(ctx, uint(requestParam.UserId))
	if err != nil {
        return "Internal Server Error", err
    }
	// 空の場合は計算処理をスキップする
	if len(modelAssets) != 0 {
		// 仮想通貨の評価総額を計算
		for _, modelAsset := range modelAssets {
			amountOfFixedIncomeAsset += modelAsset.GetPriceTotal
		}
	}
	// 登録処理を行うか、すでに資産が登録されているか確認
    totalAsset, err := ts.TotalAssetRepo.FetchTotalAssetListById(ctx, uint(requestParam.UserId),1)
    if err != nil {
        return  "Bad Request", err
    }
		// 登録内容準備
		createDto := repoTotalAsset.CreateTotalAssetDto{
			CashJpy: totalAsset[0].CashJpy,
			CashUsd: totalAsset[0].CashUsd,
			Stock: math.Round(amountOfStock),
			Fund: math.Round(amountOfFund),
			Crypto: math.Round(amountOfCrypto),
			FixedIncomeAsset: amountOfFixedIncomeAsset,
			UserId: uint(requestParam.UserId),
		}

	// 当日分の資産総額を新規登録
	_, err = ts.TotalAssetRepo.CreateTodayTotalAsset(ctx, createDto)
	if err != nil {
        return "Internal Server Error", err
    }

    return "OK", nil
}