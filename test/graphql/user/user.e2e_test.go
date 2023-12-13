package user

import (
	"bytes"
	"context"
	"encoding/json"
	serviceCurrency "my-us-stock-backend/app/graphql/currency"
	"my-us-stock-backend/app/graphql/generated"
	serviceMarketPrice "my-us-stock-backend/app/graphql/market-price"
	serviceUser "my-us-stock-backend/app/graphql/user"
	repoCurrency "my-us-stock-backend/app/repository/currency"
	repoMarketPrice "my-us-stock-backend/app/repository/market-price"
	repoUser "my-us-stock-backend/app/repository/user"
	"my-us-stock-backend/app/repository/user/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type CustomQueryResolver struct {
	UserResolver     *serviceUser.Resolver
	CurrencyResolver *serviceCurrency.Resolver
	MarketPriceResolver *serviceMarketPrice.Resolver
}

func (r *CustomQueryResolver) Query() generated.QueryResolver {
    return r
}

func (r *CustomQueryResolver) Mutation() generated.MutationResolver {
    // ここでuserResolverを使用してMutationを実装する
    return r.UserResolver
}

func (r *CustomQueryResolver) User(ctx context.Context, id string) (*generated.User, error) {
    return r.UserResolver.User(ctx, id)
}

func (r *CustomQueryResolver) GetCurrentUsdJpy(ctx context.Context) (float64, error) {
    return r.CurrencyResolver.GetCurrentUsdJpy(ctx)
}

// GetMarketPricesメソッドの実装
func (r *CustomQueryResolver) GetMarketPrices(ctx context.Context, tickers []*string) ([]*generated.MarketPrice, error) {
    // 文字列スライスに変換
    tickerStrs := make([]string, len(tickers))
    for i, t := range tickers {
        tickerStrs[i] = *t
    }

    // サービスを呼び出して結果を取得
    return r.MarketPriceResolver.GetMarketPrices(ctx, tickerStrs)
}

func NewCustomQueryResolver(userResolver *serviceUser.Resolver, currencyResolver *serviceCurrency.Resolver, marketPriceResolver *serviceMarketPrice.Resolver) *CustomQueryResolver {
    return &CustomQueryResolver{
        UserResolver:     userResolver,
        CurrencyResolver: currencyResolver,
		MarketPriceResolver: marketPriceResolver,
    }
}

// setupTestDB はテスト用のデータベースをセットアップします
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	db.AutoMigrate(&model.User{})
	return db
}

// setupGraphQLServer はテスト用のGraphQLサーバーをセットアップします
func setupGraphQLServer(db *gorm.DB) *handler.Server {
    // リポジトリ、サービス、リゾルバの初期化
	currencyRepo := repoCurrency.NewCurrencyRepository(nil)
    currencyService := serviceCurrency.NewCurrencyService(currencyRepo)
	currencyResolver := serviceCurrency.NewResolver(currencyService)

    userRepo := repoUser.NewUserRepository(db)
    userService := serviceUser.NewUserService(userRepo)
    userResolver := serviceUser.NewResolver(userService)

	marketPriceRepo := repoMarketPrice.NewMarketPriceRepository(nil)
	marketPriceService := serviceMarketPrice.NewMarketPriceService(marketPriceRepo)
    marketPriceResolver := serviceMarketPrice.NewResolver(marketPriceService)

// CustomQueryResolverの初期化
resolver := NewCustomQueryResolver(userResolver, currencyResolver, marketPriceResolver)

return handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
}


// executeGraphQLRequest はGraphQLリクエストを実行し、レスポンスを返します
func executeGraphQLRequest(h http.Handler, query string) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(map[string]interface{}{
		"query": query,
	})
	req, _ := http.NewRequest("POST", "/graphql", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}

func TestUserE2E(t *testing.T) {
	db := setupTestDB(t)
	graphqlServer := setupGraphQLServer(db)

	// テスト用のユーザーを作成
	db.Create(&model.User{Name: "Test User", Email: "test@example.com"})

	// GraphQLリクエストの実行
	query := `{ user(id: "1") { id name email } }`
	w := executeGraphQLRequest(graphqlServer, query)

	// レスポンスの検証
	assert.Equal(t, http.StatusOK, w.Code)

	// レスポンスボディの解析
	var response struct {
		Data struct {
			User struct {
				ID    string `json:"id"`
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"user"`
		} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// レスポンスボディの内容の検証
	assert.Equal(t, "1", response.Data.User.ID)
	assert.Equal(t, "Test User", response.Data.User.Name)
	assert.Equal(t, "test@example.com", response.Data.User.Email)
}

func TestCreateUserE2E(t *testing.T) {
	db := setupTestDB(t)
	graphqlServer := setupGraphQLServer(db)

	// GraphQLミューテーションの実行
	mutation := `
		mutation {
			createUser(name: "Jane Doe", email: "jane@example.com") {
				id
				name
				email
			}
		}
	`
	w := executeGraphQLRequest(graphqlServer, mutation)

	// レスポンスの検証
	assert.Equal(t, http.StatusOK, w.Code)

	// レスポンスボディの解析
	var response struct {
		Data struct {
			CreateUser struct {
				ID    string `json:"id"`
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"createUser"`
		} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// レスポンスボディの内容の検証
	assert.NotEmpty(t, response.Data.CreateUser.ID)
	assert.Equal(t, "Jane Doe", response.Data.CreateUser.Name)
	assert.Equal(t, "jane@example.com", response.Data.CreateUser.Email)
}
