package fixedincomeasset

import (
	"context"
	"my-us-stock-backend/app/common/auth"
	"my-us-stock-backend/app/graphql/generated"
	"my-us-stock-backend/app/graphql/utils"
	FixedIncome "my-us-stock-backend/app/repository/assets/fixed-income"
)

// AssetService インターフェースの定義
type AssetService interface {
    FixedIncomeAssets(ctx context.Context) ([]*generated.FixedIncomeAsset, error)
    CreateFixedIncomeAsset(ctx context.Context, input generated.CreateFixedIncomeAssetInput) (*generated.FixedIncomeAsset, error)
}

// DefaultAssetService 構造体の定義
type DefaultAssetService struct {
    Repo FixedIncome.FixedIncomeRepository // インターフェースを利用
    Auth auth.AuthService    // 認証サービスのインターフェース
}

// NewAssetService は DefaultUserService の新しいインスタンスを作成します
func NewAssetService(repo FixedIncome.FixedIncomeRepository, auth auth.AuthService) AssetService {
    return &DefaultAssetService{Repo: repo, Auth: auth}
}

// GetUserByID はユーザーをIDによって検索します
func (s *DefaultAssetService) FixedIncomeAssets(ctx context.Context) ([]*generated.FixedIncomeAsset, error) {
    // アクセストークンの検証
    userId, _ := s.Auth.FetchUserIdAccessToken(ctx)
    if userId == 0 {
        return nil, utils.UnauthenticatedError("Invalid user ID")
    }
    modelAssets, err := s.Repo.FetchFixedIncomeAssetListById(ctx, userId)
    if err != nil {
        return nil, utils.DefaultGraphQLError(err.Error())
    }
	// modelCryptosが空の場合は空の配列を返却する
	if len(modelAssets) == 0 {
		return []*generated.FixedIncomeAsset{}, nil
	}

	assets := make([]*generated.FixedIncomeAsset, 0)
	for _, modelAsset := range modelAssets {
		// pq.Int64Array to []int conversion
		paymentMonths := make([]int, len(modelAsset.PaymentMonth))
		for i, month := range modelAsset.PaymentMonth {
			paymentMonths[i] = int(month)
		}

		assets = append(assets, &generated.FixedIncomeAsset{
			ID: utils.ConvertIdToString(modelAsset.ID),
			Code: modelAsset.Code,
			GetPriceTotal: modelAsset.GetPriceTotal,
			DividendRate: modelAsset.DividendRate,
			UsdJpy: modelAsset.UsdJpy,
			PaymentMonth: paymentMonths, // Updated field
		})
	}

    return assets, nil
}

// CreateUser は新しいユーザーを作成します
func (s *DefaultAssetService) CreateFixedIncomeAsset(ctx context.Context, input generated.CreateFixedIncomeAssetInput) (*generated.FixedIncomeAsset, error) {
		// アクセストークンの検証
		userId, _ := s.Auth.FetchUserIdAccessToken(ctx)
		if userId == 0 {
			return nil, utils.UnauthenticatedError("Invalid user ID")
		}
		// pq.Int64Array to []int conversion
		paymentMonths := make([]int64, len(input.PaymentMonth))
		for i, month := range input.PaymentMonth {
			paymentMonths[i] = int64(month)
		}
	// 更新用DTOの作成
    createDto := FixedIncome.CreateFixedIncomeDto{
        Code: input.Code,
		GetPriceTotal: input.GetPriceTotal,
		DividendRate: input.DividendRate,
		UsdJpy: input.UsdJpy,
		PaymentMonth: paymentMonths,
		UserId: userId,
    }
    modelAsset, err := s.Repo.CreateFixedIncomeAsset(ctx, createDto)
    if err != nil {
        return nil, utils.DefaultGraphQLError(err.Error())
    }
	// pq.Int64Array to []int conversion
	var newPaymentMonths = make([]int, len(modelAsset.PaymentMonth))
	for i, month := range modelAsset.PaymentMonth {
		newPaymentMonths[i] = int(month)
	}
    return  &generated.FixedIncomeAsset{
		ID: utils.ConvertIdToString(modelAsset.ID),
		Code: modelAsset.Code,
		GetPriceTotal: modelAsset.GetPriceTotal,
		DividendRate: modelAsset.DividendRate,
		UsdJpy: modelAsset.UsdJpy,
		PaymentMonth: newPaymentMonths, // Updated field
	}, nil
}
