package rest

import (
	authService "my-us-stock-backend/app/common/auth"
	"my-us-stock-backend/app/common/auth/logic"
	repoUser "my-us-stock-backend/app/repository/user"

	repoCrypto "my-us-stock-backend/app/repository/assets/crypto"
	repoFixedIncome "my-us-stock-backend/app/repository/assets/fixed-income"
	repoJapanFund "my-us-stock-backend/app/repository/assets/fund"
	repoStock "my-us-stock-backend/app/repository/assets/stock"
	repoMarketPrice "my-us-stock-backend/app/repository/market-price"
	repoMarketCrypto "my-us-stock-backend/app/repository/market-price/crypto"
	repoCurrency "my-us-stock-backend/app/repository/market-price/currency"
	repoFundPrice "my-us-stock-backend/app/repository/market-price/fund"
	repoTotalAsset "my-us-stock-backend/app/repository/total-assets"
	"my-us-stock-backend/app/rest/auth"
	totalAssets "my-us-stock-backend/app/rest/total-assets"

	"my-us-stock-backend/app/rest/admin"

	"my-us-stock-backend/app/rest/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupREST は REST API のルートとコントローラを設定します
func SetupREST(r *gin.Engine, db *gorm.DB) {
    // ユーザーリポジトリの初期化
    userRepo := repoUser.NewUserRepository(db)
    totalAssetRepo := repoTotalAsset.NewTotalAssetRepository(db)
    marketPriceRepo := repoMarketPrice.NewMarketPriceRepository(nil)
    usStockRepo := repoStock.NewUsStockRepository(db)
    currencyRepo := repoCurrency.NewCurrencyRepository(nil)
    japanFundRepo := repoJapanFund.NewJapanFundRepository(db)
    marketCryptoRepo := repoMarketCrypto.NewCryptoRepository(nil)
    fixedIncomeAssetRepo := repoFixedIncome.NewFixedIncomeRepository(db)
    cryptoRepo := repoCrypto.NewCryptoRepository(db)
    fundPriceRepo := repoFundPrice.NewFetchFundRepository(db)

    // 認証機能
    userLogic := logic.NewUserLogic()
    responseLogic := logic.NewResponseLogic()
    jwtLogic := logic.NewJWTLogic()
    authValidation := authService.NewAuthValidation()

    // RESTサービス、コントローラの初期化
    userRestService := user.NewUserService(userRepo)
    userController := user.NewUserController(userRestService)

    // 認証サービスの初期化
    authService := authService.NewAuthService(userRepo, userLogic, responseLogic, jwtLogic, authValidation)
    authController := auth.NewAuthController(authService)

    totalAssetService := totalAssets.NewTotalAssetService(totalAssetRepo, usStockRepo, marketPriceRepo, currencyRepo, japanFundRepo, cryptoRepo, fixedIncomeAssetRepo, marketCryptoRepo)
    totalAssetController := totalAssets.NewTotalAssetController(totalAssetService)

    adminService := admin.NewFundPriceService(fundPriceRepo)
    adminController := admin.NewFundPriceController(adminService)

    // RESTコントローラのルートを設定
    r.GET("/api/users/:id", userController.GetUser)
    r.POST("/api/users", userController.CreateUser)
    // 認証用
    r.POST("/api/v1/signin", authController.SignIn)
    r.POST("/api/v1/signup", authController.SignUp)
    r.POST("/api/v1/total-assets", totalAssetController.CreateTodayTotalAsset)
    r.POST("/api/v1/refresh", authController.RefreshAccessToken)
    // 管理画面用
    r.GET("/api/v1/admin/fund-prices", adminController.GetFundPrices)
    r.POST("/api/v1/admin/fund-prices", adminController.CreateFundPrice)
    r.PUT("/api/v1/admin/fund-prices", adminController.UpdateFundPrice)
}
