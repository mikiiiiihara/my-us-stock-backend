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
	UpdateJapanFund(ctx context.Context, input generated.UpdateJapanFundInput) (*generated.JapanFund, error)
	DeleteJapanFund(ctx context.Context, id string) (bool, error)
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

	funds := make([]*generated.JapanFund, len(modelFunds))
	for i, modelFund := range modelFunds {
		funds[i] = &generated.JapanFund{
			ID: utils.ConvertIdToString(modelFund.ID),
			Code: modelFund.Code,
			Name: modelFund.Name,
			GetPrice: modelFund.GetPrice,
			GetPriceTotal: modelFund.GetPriceTotal,
			CurrentPrice: getFundMarketPrice(modelFund.Code), // TODO: 三菱UFJのAPI復旧したらrepository実装してそこから取得するようにする
		}
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
		CurrentPrice: getFundMarketPrice(modelFund.Code),
	}, nil
}

func (s *DefaultJapanFundService) UpdateJapanFund(ctx context.Context, input generated.UpdateJapanFundInput) (*generated.JapanFund, error) {
	// アクセストークンの検証
	userId, _ := s.Auth.FetchUserIdAccessToken(ctx)
	if userId == 0 {
		return nil, utils.UnauthenticatedError("Invalid user ID")
	}
	updateId, convertError := utils.ConvertIdToUint(input.ID)
	if convertError != nil || updateId == 0 {
        return nil, utils.DefaultGraphQLError("入力されたidが無効です")
       }

	// 更新用DTOの作成
	updateDto := JapanFund.UpdateJapanFundDto{
		ID: updateId,
        GetPrice: &input.GetPrice,
		GetPriceTotal: &input.GetPriceTotal,
	}
    modelFund, err := s.Repo.UpdateJapanFund(ctx, updateDto)
    if err != nil {
        return nil, utils.DefaultGraphQLError(err.Error())
    }
	return  &generated.JapanFund{
		ID: utils.ConvertIdToString(modelFund.ID),
		Code: modelFund.Code,
		Name: modelFund.Name,
		GetPriceTotal: modelFund.GetPriceTotal,
		GetPrice: modelFund.GetPrice,
		CurrentPrice: getFundMarketPrice(modelFund.Code),
	}, nil
}

// TODO: 三菱UFJのAPI復旧したらrepository実装してそこから取得するようにする
func getFundMarketPrice(code string) float64 {
	switch code {
	case "SP500":
		return 27091.0
	case "全世界株":
		return 23040.0
	default:
		return 18768.0
	}
}

// 削除
func (s *DefaultJapanFundService) DeleteJapanFund(ctx context.Context, id string) (bool, error) {
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
    var err = s.Repo.DeleteJapanFund(ctx, deleteId)
	// 市場情報を追加して返却
    if err != nil {
     return false, utils.DefaultGraphQLError(err.Error())
    }
	return true, nil
}