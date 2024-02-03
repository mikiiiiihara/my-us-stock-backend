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
	UpdateFixedIncomeAsset(ctx context.Context, input generated.UpdateFixedIncomeAssetInput) (*generated.FixedIncomeAsset, error)
	DeleteFixedIncomeAsset(ctx context.Context, id string) (bool, error)
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
	// modelAssetsが空の場合は空の配列を返却する
	if len(modelAssets) == 0 {
		return []*generated.FixedIncomeAsset{}, nil
	}

	assets := make([]*generated.FixedIncomeAsset, len(modelAssets))
	for i, modelAsset := range modelAssets {
		// pq.Int64Array to []int conversion
		paymentMonths := make([]int, len(modelAsset.PaymentMonth))
		for i, month := range modelAsset.PaymentMonth {
			paymentMonths[i] = int(month)
		}

		assets[i] = &generated.FixedIncomeAsset{
			ID: utils.ConvertIdToString(modelAsset.ID),
			Code: modelAsset.Code,
			GetPriceTotal: modelAsset.GetPriceTotal,
			DividendRate: modelAsset.DividendRate,
			UsdJpy: modelAsset.UsdJpy,
			PaymentMonth: paymentMonths, // Updated field
		}
	}

    return assets, nil
}

func (s *DefaultAssetService) CreateFixedIncomeAsset(ctx context.Context, input generated.CreateFixedIncomeAssetInput) (*generated.FixedIncomeAsset, error) {
		// アクセストークンの検証
		userId, _ := s.Auth.FetchUserIdAccessToken(ctx)
		if userId == 0 {
			return nil, utils.UnauthenticatedError("Invalid user ID")
		}
		// pq.Int64Array to []int conversion (without using append)
		paymentMonths := make([]int64, 0, len(input.PaymentMonth))
		for _, month := range input.PaymentMonth {
			paymentMonths = append(paymentMonths, int64(month))
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

func (s *DefaultAssetService) UpdateFixedIncomeAsset(ctx context.Context, input generated.UpdateFixedIncomeAssetInput) (*generated.FixedIncomeAsset, error) {
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
	updateDto := FixedIncome.UpdateFixedIncomeDto{
		ID: updateId,
		GetPriceTotal: &input.GetPriceTotal,
		UsdJpy: input.UsdJpy,
	}
	modelAsset, err := s.Repo.UpdateFixedIncomeAsset(ctx, updateDto)
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

// 削除
func (s *DefaultAssetService) DeleteFixedIncomeAsset(ctx context.Context, id string) (bool, error) {
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
    var err = s.Repo.DeleteFixedIncomeAsset(ctx, deleteId)
	// 市場情報を追加して返却
    if err != nil {
     return false, utils.DefaultGraphQLError(err.Error())
    }
	return true, nil
}