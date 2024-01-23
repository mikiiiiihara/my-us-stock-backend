package totalasset

import (
	"context"
	"my-us-stock-backend/app/common/auth"
	"my-us-stock-backend/app/graphql/generated"
	"my-us-stock-backend/app/graphql/utils"
	Repo "my-us-stock-backend/app/repository/total-assets"
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
    Repo Repo.TotalAssetRepository // インターフェースを利用
    Auth auth.AuthService    // 認証サービスのインターフェース
}

// NewTotalAssetService は DefaultUserService の新しいインスタンスを作成します
func NewTotalAssetService(repo Repo.TotalAssetRepository, auth auth.AuthService) TotalAssetService {
    return &DefaultTotalAssetService{Repo: repo, Auth: auth}
}

// GetUserByID はユーザーをIDによって検索します
func (s *DefaultTotalAssetService) TotalAssets(ctx context.Context, day int) ([]*generated.TotalAsset, error) {
    // アクセストークンの検証
    userId, _ := s.Auth.FetchUserIdAccessToken(ctx)
    if userId == 0 {
        return nil, utils.UnauthenticatedError("Invalid user ID")
    }
    modelAssets, err := s.Repo.FetchTotalAssetListById(ctx, userId, day)
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
	updateDto := Repo.UpdateTotalAssetDto{
		ID: updateId,
		CashUsd: &input.CashUsd,
		CashJpy: &input.CashJpy,
	}
	updatedAsset, err := s.Repo.UpdateTotalAsset(ctx, updateDto)
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