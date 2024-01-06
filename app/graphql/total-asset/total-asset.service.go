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

	assets := make([]*generated.TotalAsset, 0)
	for _, modelAsset := range modelAssets {

		assets = append(assets, &generated.TotalAsset{
			ID: utils.ConvertIdToString(modelAsset.ID),
			CashJpy: modelAsset.CashJpy,
			CashUsd: modelAsset.CashUsd,
			Stock: modelAsset.Stock,
			Fund: modelAsset.Fund,
			Crypto: modelAsset.Crypto,
			FixedIncomeAsset: modelAsset.FixedIncomeAsset,
			CreatedAt: modelAsset.CreatedAt.Format(time.RFC3339),
		})
	}

	// assetsをCreatedAtで昇順にソート
	sort.Slice(assets, func(i, j int) bool {
		return assets[i].CreatedAt < assets[j].CreatedAt
	})

    return assets, nil
}