package fixedincome

import (
	"context"
	"fmt"
	"my-us-stock-backend/app/database/model"

	"gorm.io/gorm"
)

// FixedIncomeRepository インターフェースの定義
type FixedIncomeRepository interface {
	FetchFixedIncomeAssetListById(ctx context.Context, userId uint) ([]model.FixedIncomeAsset, error)
    UpdateFixedIncomeAsset(ctx context.Context, dto UpdateFixedIncomeDto) (*model.FixedIncomeAsset, error)
	CreateFixedIncomeAsset(ctx context.Context, dto CreateFixedIncomeDto) (*model.FixedIncomeAsset, error)
	DeleteFixedIncomeAsset(ctx context.Context, id uint) error
}

// DefaultFixedIncomeRepository 構造体の定義
type DefaultFixedIncomeRepository struct {
    DB *gorm.DB
}

// 共通フィールドを選択するためのヘルパー関数です。
func selectBaseQuery(db *gorm.DB) *gorm.DB {
    return db.Select("id", "get_price_total", "dividend_rate", "code", "usd_jpy", "payment_month", "user_id")
}

// NewCryptoRepository は DefaultStrategyRepository の新しいインスタンスを作成します
func NewFixedIncomeRepository(db *gorm.DB) FixedIncomeRepository {
    return &DefaultFixedIncomeRepository{DB: db}
}

// 指定したuserIdのユーザーが保有する米国株式のリストを取得する
func (r *DefaultFixedIncomeRepository) FetchFixedIncomeAssetListById(ctx context.Context, userId uint) ([]model.FixedIncomeAsset, error) {
    var fixedIncomeAssets []model.FixedIncomeAsset
    err := selectBaseQuery(r.DB).Where("user_id = ?", userId).Find(&fixedIncomeAssets).Error
    if err != nil {
        return nil, err
    }
    return fixedIncomeAssets, nil
}

// 米国株式情報を更新します
func (r *DefaultFixedIncomeRepository) UpdateFixedIncomeAsset(ctx context.Context, dto UpdateFixedIncomeDto) (*model.FixedIncomeAsset, error) {
    // 更新用のマップを作成します
    newFixedIncomeAsset := map[string]interface{}{}

    if dto.GetPriceTotal != nil {
        newFixedIncomeAsset["get_price_total"] = dto.GetPriceTotal
    }
    if dto.DividendRate != nil {
        newFixedIncomeAsset["dividend_rate"] = dto.DividendRate
    }
	if dto.UsdJpy != nil {
        newFixedIncomeAsset["usdjpy"] = dto.UsdJpy
    }

    // 指定されたIDの株式情報を更新します
    if err := r.DB.Model(&model.FixedIncomeAsset{}).Where("id = ?", dto.ID).Updates(newFixedIncomeAsset).Error; err != nil {
        return nil, err
    }

    // 更新された情報を取得します
    var fixedIncomeAsset model.FixedIncomeAsset
    if err := r.DB.Where("id = ?", dto.ID).Find(&fixedIncomeAsset).Error; err != nil {
        return nil, err
    }

    return &fixedIncomeAsset, nil
}

// 米国株式情報を作成します
func (r *DefaultFixedIncomeRepository) CreateFixedIncomeAsset(ctx context.Context, dto CreateFixedIncomeDto) (*model.FixedIncomeAsset, error) {
    // 既に同じ銘柄が存在するかを確認
    var existingUsStock model.FixedIncomeAsset
    if err := selectBaseQuery(r.DB).Where("code = ?", dto.Code).First(&existingUsStock).Error; err == nil {
        return nil, fmt.Errorf("この銘柄は既に登録されています")
    }

    // 新しい米国株式情報を作成
    fixedIncomeAsset := &model.FixedIncomeAsset{
        Code:   dto.Code,
        GetPriceTotal: dto.GetPriceTotal,
        DividendRate: dto.DividendRate,
		UsdJpy: dto.UsdJpy,
        PaymentMonth: dto.PaymentMonth,
        UserId:   dto.UserId,
    }

    if err := r.DB.Create(&fixedIncomeAsset).Error; err != nil {
        return nil, err
    }

    return fixedIncomeAsset, nil
}

// 米国株式情報を削除します
func (r *DefaultFixedIncomeRepository) DeleteFixedIncomeAsset(ctx context.Context, id uint) error {
    // 指定されたIDの株式情報を検索して削除
    if err := r.DB.Where("id = ?", id).Delete(&model.FixedIncomeAsset{}).Error; err != nil {
        return err
    }
    return nil
}

