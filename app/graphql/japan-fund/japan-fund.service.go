package japanfund

import (
	"context"
	"my-us-stock-backend/app/common/auth"
	"my-us-stock-backend/app/graphql/generated"
	"my-us-stock-backend/app/graphql/utils"
	JapanFund "my-us-stock-backend/app/repository/assets/fund"
)

// JapanFundService インターフェースの定義
type JapanFundService interface {
    JapanFunds(ctx context.Context) ([]*generated.JapanFund, error)
	CreateJapanFund(ctx context.Context, input generated.CreateJapanFundInput) (*generated.JapanFund, error)
}

// DefaultJapanFundService 構造体の定義
type DefaultJapanFundService struct {
    Repo JapanFund.JapanFundRepository // インターフェースを利用
    Auth auth.AuthService    // 認証サービスのインターフェース
}

// NewJapanFundService は DefaultUserService の新しいインスタンスを作成します
func NewJapanFundService(repo JapanFund.JapanFundRepository, auth auth.AuthService) JapanFundService {
    return &DefaultJapanFundService{Repo: repo, Auth: auth}
}

// GetUserByID はユーザーをIDによって検索します
func (s *DefaultJapanFundService) JapanFunds(ctx context.Context) ([]*generated.JapanFund, error) {
    // アクセストークンの検証
    userId, _ := s.Auth.FetchUserIdAccessToken(ctx)
    if userId == 0 {
        return nil, utils.UnauthenticatedError("Invalid user ID")
    }
    modelFunds, err := s.Repo.FetchJapanFundListById(ctx, userId)
    if err != nil {
        return nil, utils.DefaultGraphQLError(err.Error())
    }
	// modelFundsが空の場合は空の配列を返却する
	if len(modelFunds) == 0 {
		return []*generated.JapanFund{}, nil
	}

	funds := make([]*generated.JapanFund, 0)
	for _, modelFund := range modelFunds {

		funds = append(funds, &generated.JapanFund{
			ID: utils.ConvertIdToString(modelFund.ID),
			Code: modelFund.Code,
			Name: modelFund.Name,
			GetPrice: modelFund.GetPrice,
			GetPriceTotal: modelFund.GetPriceTotal,
			CurrentPrice: getFundMarketPrice(modelFund.Code), // TODO: 三菱UFJのAPI復旧したらrepository実装してそこから取得するようにする
		})
	}

    return funds, nil
}

// CreateUser は新しいユーザーを作成します
func (s *DefaultJapanFundService) CreateJapanFund(ctx context.Context, input generated.CreateJapanFundInput) (*generated.JapanFund, error) {
		// アクセストークンの検証
		userId, _ := s.Auth.FetchUserIdAccessToken(ctx)
		if userId == 0 {
			return nil, utils.UnauthenticatedError("Invalid user ID")
		}
	// 更新用DTOの作成
    createDto := JapanFund.CreateJapanFundDto{
        Code: input.Code,
		Name: input.Name,
		GetPriceTotal: input.GetPriceTotal,
		GetPrice: input.GetPrice,
		UserId: userId,
    }
    modelFund, err := s.Repo.CreateJapanFund(ctx, createDto)
    if err != nil {
        return nil, utils.DefaultGraphQLError(err.Error())
    }
    return  &generated.JapanFund{
		ID: utils.ConvertIdToString(modelFund.ID),
		Code: modelFund.Code,
		Name: modelFund.Name,
		GetPriceTotal: modelFund.GetPriceTotal,
		GetPrice: modelFund.GetPrice,
		CurrentPrice: getFundMarketPrice(modelFund.Code),// TODO: 三菱UFJのAPI復旧したらrepository実装してそこから取得するようにする
	}, nil
}

// TODO: 三菱UFJのAPI復旧したらrepository実装してそこから取得するようにする
func getFundMarketPrice(code string) float64 {
	switch code {
	case "SP500":
		return 24281.0
	case "全世界株":
		return 21084.0
	default:
		return 18905.0
	
	}
}