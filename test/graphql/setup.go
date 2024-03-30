package graphql

import (
	authService "my-us-stock-backend/app/common/auth"
	"my-us-stock-backend/app/common/auth/logic"
	"my-us-stock-backend/app/graphql"
	"my-us-stock-backend/app/graphql/crypto"
	serviceCurrency "my-us-stock-backend/app/graphql/currency"
	serviceFixedIncomeAsset "my-us-stock-backend/app/graphql/fixed-income-asset"
	serviceJapanFund "my-us-stock-backend/app/graphql/japan-fund"
	serviceMarketPrice "my-us-stock-backend/app/graphql/market-price"
	serviceStock "my-us-stock-backend/app/graphql/stock"
	serviceTotalAsset "my-us-stock-backend/app/graphql/total-asset"
	serviceUser "my-us-stock-backend/app/graphql/user"
	repoCrypto "my-us-stock-backend/app/repository/assets/crypto"
	repoFixedIncome "my-us-stock-backend/app/repository/assets/fixed-income"
	repoJapanFund "my-us-stock-backend/app/repository/assets/fund"
	repoStock "my-us-stock-backend/app/repository/assets/stock"
	repoMarketPrice "my-us-stock-backend/app/repository/market-price"
	repoMarketCrypto "my-us-stock-backend/app/repository/market-price/crypto"
	repoCurrency "my-us-stock-backend/app/repository/market-price/currency"
	repoFundPrice "my-us-stock-backend/app/repository/market-price/fund"
	repoTotalAsset "my-us-stock-backend/app/repository/total-assets"
	repoUser "my-us-stock-backend/app/repository/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupOptions - GraphQLサーバーのセットアップオプション
type SetupOptions struct {
    MockHTTPClient    *http.Client
    CurrencyRepo      repoCurrency.CurrencyRepository
    UserRepo          repoUser.UserRepository
    MarketPriceRepo   repoMarketPrice.MarketPriceRepository
    MarketCryptoRepo   repoMarketCrypto.CryptoRepository
    UsStockRepo   repoStock.UsStockRepository
    CryptoRepo   repoCrypto.CryptoRepository
    FixedIncomeAssetRepo repoFixedIncome.FixedIncomeRepository
    JapanFundRepo repoJapanFund.JapanFundRepository
    TotalAssetRepo repoTotalAsset.TotalAssetRepository
    FundPriceRepo repoFundPrice.FundPriceRepository
}

// SetupGraphQLServer - GraphQLサーバーのセットアップ
func SetupGraphQLServer(db *gorm.DB, opts *SetupOptions) http.Handler {
    var currencyRepo repoCurrency.CurrencyRepository
    var userRepo repoUser.UserRepository
    var marketPriceRepo repoMarketPrice.MarketPriceRepository
    var marketCryptoRepo repoMarketCrypto.CryptoRepository
    var usStockRepo repoStock.UsStockRepository
    var cryptoRepo repoCrypto.CryptoRepository
    var fixedIncomeAssetRepo repoFixedIncome.FixedIncomeRepository
    var japanFundRepo repoJapanFund.JapanFundRepository
    var totalAssetRepo repoTotalAsset.TotalAssetRepository
    var fundPriceRepo repoFundPrice.FundPriceRepository

    // optsがnilでない場合にのみ、各リポジトリを設定
    if opts != nil {
        currencyRepo = opts.CurrencyRepo
        userRepo = opts.UserRepo
        marketPriceRepo = opts.MarketPriceRepo
        usStockRepo = opts.UsStockRepo
        cryptoRepo = opts.CryptoRepo
        fixedIncomeAssetRepo = opts.FixedIncomeAssetRepo
        japanFundRepo = opts.JapanFundRepo
        totalAssetRepo = opts.TotalAssetRepo
        fundPriceRepo = opts.FundPriceRepo
    }

    // 各リポジトリがまだnilの場合、デフォルトのリポジトリを使用
    if currencyRepo == nil {
        currencyRepo = repoCurrency.NewCurrencyRepository(nil)
    }
    if userRepo == nil {
        userRepo = repoUser.NewUserRepository(db)
    }
    if usStockRepo == nil {
        usStockRepo = repoStock.NewUsStockRepository(db)
    }

    if cryptoRepo == nil {
        cryptoRepo = repoCrypto.NewCryptoRepository(db)
    }

    if fixedIncomeAssetRepo == nil {
        fixedIncomeAssetRepo = repoFixedIncome.NewFixedIncomeRepository(db)
    }

    if totalAssetRepo == nil {
        totalAssetRepo = repoTotalAsset.NewTotalAssetRepository(db)
    }

    if japanFundRepo == nil {
        japanFundRepo = repoJapanFund.NewJapanFundRepository(db)
    }

    if fundPriceRepo == nil {
        fundPriceRepo = repoFundPrice.NewFetchFundRepository(db)
    }

    if marketPriceRepo == nil {
        // 注意: ここでは opts が nil の可能性があるため、opts.MockHTTPClient の前に nil チェックが必要
        var httpClient *http.Client
        if opts != nil {
            httpClient = opts.MockHTTPClient
        }
        marketPriceRepo = repoMarketPrice.NewMarketPriceRepository(httpClient)
    }

    if marketCryptoRepo == nil {
        // 注意: ここでは opts が nil の可能性があるため、opts.MockHTTPClient の前に nil チェックが必要
        var httpClient *http.Client
        if opts != nil {
            httpClient = opts.MockHTTPClient
        }
        marketCryptoRepo = repoMarketCrypto.NewCryptoRepository(httpClient)
    }

    // 認証機能
    userLogic := logic.NewUserLogic()
    responseLogic := logic.NewResponseLogic()
    jwtLogic := logic.NewJWTLogic()
    authValidation := authService.NewAuthValidation()

    // 認証サービスの初期化
    authService := authService.NewAuthService(userRepo, userLogic, responseLogic, jwtLogic, authValidation)

    // サービスとリゾルバの初期化
    currencyService := serviceCurrency.NewCurrencyService(currencyRepo)
    currencyResolver := serviceCurrency.NewResolver(currencyService)
    userService := serviceUser.NewUserService(userRepo, authService)
    userResolver := serviceUser.NewResolver(userService)
    marketPriceService := serviceMarketPrice.NewMarketPriceService(marketPriceRepo)
    marketPriceResolver := serviceMarketPrice.NewResolver(marketPriceService)

    usStockService := serviceStock.NewUsStockService(usStockRepo, authService,marketPriceRepo)
    usStockResolver := serviceStock.NewResolver(usStockService)

    cryptoService := crypto.NewCryptoService(cryptoRepo, authService, marketCryptoRepo)
    cryptoResolver := crypto.NewResolver(cryptoService)

    fixedIncomeAssetService := serviceFixedIncomeAsset.NewAssetService(fixedIncomeAssetRepo, authService)
    fixedIncomeAssetResolver := serviceFixedIncomeAsset.NewResolver(fixedIncomeAssetService)

    japanFundService := serviceJapanFund.NewJapanFundService(japanFundRepo, authService, fundPriceRepo)
    japanFundResolver := serviceJapanFund.NewResolver(japanFundService)

    totalAssetService := serviceTotalAsset.NewTotalAssetService(authService,totalAssetRepo, usStockRepo, marketPriceRepo, currencyRepo, japanFundRepo, cryptoRepo, fixedIncomeAssetRepo, marketCryptoRepo, fundPriceRepo)
    totalAssetResolver := serviceTotalAsset.NewResolver(totalAssetService)
    // Ginのルーターを初期化
    r := gin.Default()
    // GraphQLのエンドポイントを設定
    r.POST("/graphql", graphql.GinContextToGraphQLMiddleware(), graphql.Handler(userResolver, currencyResolver, marketPriceResolver,usStockResolver, cryptoResolver,fixedIncomeAssetResolver,japanFundResolver,totalAssetResolver))

    return r
}