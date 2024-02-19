package japanfund

import (
	"context"
	"my-us-stock-backend/app/common/auth"
	"my-us-stock-backend/app/database/model"
	"my-us-stock-backend/app/graphql/generated"
	"my-us-stock-backend/app/graphql/utils"
	JapanFund "my-us-stock-backend/app/repository/assets/fund"
	marketPrice "my-us-stock-backend/app/repository/market-price/fund"
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
	MarketPriceRepo marketPrice.FundPriceRepository
}

// NewJapanFundService は DefaultUserService の新しいインスタンスを作成します
func NewJapanFundService(repo JapanFund.JapanFundRepository, auth auth.AuthService, marketPriceRepo marketPrice.FundPriceRepository) JapanFundService {
    return &DefaultJapanFundService{Repo: repo, Auth: auth, MarketPriceRepo: marketPriceRepo}
}

func (s *DefaultJapanFundService) JapanFunds(ctx context.Context) ([]*generated.JapanFund, error) {
    userId, _ := s.Auth.FetchUserIdAccessToken(ctx)
    if userId == 0 {
        return nil, utils.UnauthenticatedError("Invalid user ID")
    }
    modelFunds, err := s.Repo.FetchJapanFundListById(ctx, userId)
    if err != nil {
        return nil, utils.DefaultGraphQLError(err.Error())
    }
    if len(modelFunds) == 0 {
        return []*generated.JapanFund{}, nil
    }

    funds := make([]*generated.JapanFund, len(modelFunds))
    errChan := make(chan error, 1)
    resultChan := make(chan struct {
        Index int
        FundPrice *model.FundPrice
    }, len(modelFunds))

    for i, modelFund := range modelFunds {
        go func(i int, code string) {
            fundPrice, err := s.MarketPriceRepo.FindFundPriceByCode(ctx, code)
            if err != nil {
                errChan <- err
                return
            }
            resultChan <- struct {
                Index int
                FundPrice *model.FundPrice
            }{i, fundPrice}
        }(i, modelFund.Code)
    }

    for i := 0; i < len(modelFunds); i++ {
        select {
        case result := <-resultChan:
            funds[result.Index] = &generated.JapanFund{
                ID: utils.ConvertIdToString(modelFunds[result.Index].ID),
                Code: modelFunds[result.Index].Code,
                Name: modelFunds[result.Index].Name,
                GetPrice: modelFunds[result.Index].GetPrice,
                GetPriceTotal: modelFunds[result.Index].GetPriceTotal,
                CurrentPrice: result.FundPrice.Price,
            }
        case err := <-errChan:
            return nil, utils.DefaultGraphQLError(err.Error())
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
	fundPrice, err := s.MarketPriceRepo.FindFundPriceByCode(ctx, createDto.Code)
	if err != nil {
        return nil, utils.DefaultGraphQLError(err.Error())
    }
    return  &generated.JapanFund{
		ID: utils.ConvertIdToString(modelFund.ID),
		Code: modelFund.Code,
		Name: modelFund.Name,
		GetPriceTotal: modelFund.GetPriceTotal,
		GetPrice: modelFund.GetPrice,
		CurrentPrice: fundPrice.Price,
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
	fundPrice, err := s.MarketPriceRepo.FindFundPriceByCode(ctx, modelFund.Code)
	if err != nil {
        return nil, utils.DefaultGraphQLError(err.Error())
    }
	return  &generated.JapanFund{
		ID: utils.ConvertIdToString(modelFund.ID),
		Code: modelFund.Code,
		Name: modelFund.Name,
		GetPriceTotal: modelFund.GetPriceTotal,
		GetPrice: modelFund.GetPrice,
		CurrentPrice: fundPrice.Price,
	}, nil
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