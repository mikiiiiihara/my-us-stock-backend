package fund

import (
	"context"
	"fmt"
	"my-us-stock-backend/app/database/model"

	"gorm.io/gorm"
)

// JapanFundRepository インターフェースの定義
type JapanFundRepository interface {
	FetchJapanFundListById(ctx context.Context, userId uint) ([]model.JapanFund, error)
    UpdateJapanFund(ctx context.Context, dto UpdateJapanFundDto) (*model.JapanFund, error)
	CreateJapanFund(ctx context.Context, dto CreateJapanFundDto) (*model.JapanFund, error)
	DeleteJapanFund(ctx context.Context, id uint) error
}

// DefaultJapanFundRepository 構造体の定義
type DefaultJapanFundRepository struct {
    DB *gorm.DB
}

// 共通フィールドを選択するためのヘルパー関数です。
func selectBaseQuery(db *gorm.DB) *gorm.DB {
    return db.Select("id", "get_price_total", "get_price", "code", "name", "user_id")
}

// NewJapanFundRepository は DefaultStrategyRepository の新しいインスタンスを作成します
func NewJapanFundRepository(db *gorm.DB) JapanFundRepository {
    return &DefaultJapanFundRepository{DB: db}
}

// 指定したuserIdのユーザーが保有する日本投資信託のリストを取得する
func (r *DefaultJapanFundRepository) FetchJapanFundListById(ctx context.Context, userId uint) ([]model.JapanFund, error) {
    var funds []model.JapanFund
    err := selectBaseQuery(r.DB).Where("user_id = ?", userId).Find(&funds).Error
    if err != nil {
        return nil, err
    }
    return funds, nil
}

// 日本投資信託情報を更新します
func (r *DefaultJapanFundRepository) UpdateJapanFund(ctx context.Context, dto UpdateJapanFundDto) (*model.JapanFund, error) {
    // 更新用のマップを作成します
    newFund := map[string]interface{}{}

    if dto.GetPrice != nil {
        newFund["get_price"] = dto.GetPrice
    }
    if dto.GetPriceTotal != nil {
        newFund["get_price_total"] = dto.GetPriceTotal
    }

    // 指定されたIDの株式情報を更新します
    if err := r.DB.Model(&model.JapanFund{}).Where("id = ?", dto.ID).Updates(newFund).Error; err != nil {
        return nil, err
    }

    // 更新された情報を取得します
    var JapanFund model.JapanFund
    if err := selectBaseQuery(r.DB).Where("id = ?", dto.ID).Find(&JapanFund).Error; err != nil {
        return nil, err
    }

    return &JapanFund, nil
}

// 日本投資信託情報を作成します
func (r *DefaultJapanFundRepository) CreateJapanFund(ctx context.Context, dto CreateJapanFundDto) (*model.JapanFund, error) {
    // 既に同じ銘柄が存在するかを確認
    var existingJapanFund model.JapanFund
    if err := selectBaseQuery(r.DB).Where("code = ? AND user_id = ? AND name = ?", dto.Code, dto.UserId, dto.Name).First(&existingJapanFund).Error; err == nil {
        return nil, fmt.Errorf("この銘柄は既に登録されています")
    }

    // 新しい米国株式情報を作成
    japanFund := &model.JapanFund{
        Name: dto.Name,
        Code:   dto.Code,
        GetPriceTotal: dto.GetPriceTotal,
        GetPrice: dto.GetPrice,
        UserId:   dto.UserId,
    }

    if err := r.DB.Create(&japanFund).Error; err != nil {
        return nil, err
    }

    return japanFund, nil
}

// 日本投資信託情報を削除します
func (r *DefaultJapanFundRepository) DeleteJapanFund(ctx context.Context, id uint) error {
    // 指定されたIDの株式情報を検索して削除
    if err := r.DB.Where("id = ?", id).Delete(&model.JapanFund{}).Error; err != nil {
        return err
    }
    return nil
}

