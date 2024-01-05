package totalassets

import (
	"context"
	"errors"
	"my-us-stock-backend/app/database/model"
	"time"

	"gorm.io/gorm"
)

// TotalAssetRepository インターフェースの定義
type TotalAssetRepository interface {
	FetchTotalAssetListById(ctx context.Context, userId uint) ([]model.TotalAsset, error)
    UpdateTotalAsset(ctx context.Context, dto UpdateTotalAssetDto) (*model.TotalAsset, error)
	CreateTodayTotalAsset(ctx context.Context, dto CreateTotalAssetDto) (*model.TotalAsset, error)
}

// DefaultTotalAssetRepository 構造体の定義
type DefaultTotalAssetRepository struct {
    DB *gorm.DB
}

// NewTotalAssetRepository は DefaultTotalAssetRepository の新しいインスタンスを作成します
func NewTotalAssetRepository(db *gorm.DB) TotalAssetRepository {
    return &DefaultTotalAssetRepository{DB: db}
}

// 指定したuserIdのユーザーが保有する米国株式のリストを取得する
func (r *DefaultTotalAssetRepository) FetchTotalAssetListById(ctx context.Context, userId uint) ([]model.TotalAsset, error) {
    var assets []model.TotalAsset
    err := r.DB.Where("user_id = ?", userId).Find(&assets).Error
    if err != nil {
        return nil, err
    }
    return assets, nil
}

// 米国株式情報を更新します
func (r *DefaultTotalAssetRepository) UpdateTotalAsset(ctx context.Context, dto UpdateTotalAssetDto) (*model.TotalAsset, error) {
    var existingAsset model.TotalAsset
    // 対象となるレコードが存在するかどうかをチェック
    if err := r.DB.First(&existingAsset, dto.ID).Error; err != nil {
        return nil, err // レコードが存在しない場合、エラーを返す
    }
    // 更新用のマップを作成します
    newAsset := map[string]interface{}{}

    if dto.CashJpy != nil {
        newAsset["cash_jpy"] = dto.CashJpy
    }
    if dto.CashUsd != nil {
        newAsset["cash_usd"] = dto.CashUsd
    }
    if dto.Stock != nil {
        newAsset["stock"] = dto.Stock
    }
    if dto.Crypto != nil {
        newAsset["crypto"] = dto.Crypto
    }
    if dto.FixedIncomeAsset != nil {
        newAsset["fixed_income_asset"] = dto.FixedIncomeAsset
    }


    // 指定されたIDの株式情報を更新します
    if err := r.DB.Model(&model.TotalAsset{}).Where("id = ?", dto.ID).Updates(newAsset).Error; err != nil {
        return nil, err
    }

    // 更新された情報を取得します
    var asset model.TotalAsset
    if err := r.DB.Where("id = ?", dto.ID).Find(&asset).Error; err != nil {
        return nil, err
    }

    return &asset, nil
}

// 米国株式情報を作成します
func (r *DefaultTotalAssetRepository) CreateTodayTotalAsset(ctx context.Context, dto CreateTotalAssetDto) (*model.TotalAsset, error) {
    // 現在の日付を取得
    today := time.Now().UTC().Format("2006-01-02")

    // 同じ日付で同じユーザーのレコードが存在するか確認
    var existingAsset model.TotalAsset
    if err := r.DB.Where("user_id = ? AND DATE(created_at) = ?", dto.UserId, today).First(&existingAsset).Error; err == nil {
        // レコードが存在する場合、エラーを返す
        return nil, errors.New("既に資産が登録されています。新規追加ではなく更新を行ってください。")
    } else if !errors.Is(err, gorm.ErrRecordNotFound) {
        // 他のエラーが発生した場合
        return nil, err
    }


    newAsset := &model.TotalAsset{
        CashJpy:   dto.CashJpy,
        CashUsd:   dto.CashUsd,
        Stock:   dto.Stock,
        Crypto:   dto.Crypto,
        FixedIncomeAsset: dto.FixedIncomeAsset,
        UserId:   dto.UserId,
    }

    if err := r.DB.Create(&newAsset).Error; err != nil {
        return nil, err
    }

    return newAsset, nil
}