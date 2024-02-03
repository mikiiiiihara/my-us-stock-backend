package stock

import (
	"context"
	"fmt"
	"my-us-stock-backend/app/database/model"

	"gorm.io/gorm"
)

// UsStockRepository インターフェースの定義
type UsStockRepository interface {
	FetchUsStockListById(ctx context.Context, userId uint) ([]model.UsStock, error)
    UpdateUsStock(ctx context.Context, dto UpdateUsStockDto) (*model.UsStock, error)
	CreateUsStock(ctx context.Context, dto CreateUsStockDto) (*model.UsStock, error)
	DeleteUsStock(ctx context.Context, id uint) error
}

// DefaultUsStockRepository 構造体の定義
type DefaultUsStockRepository struct {
    DB *gorm.DB
}

// 共通フィールドを選択するためのヘルパー関数です。
func selectBaseQuery(db *gorm.DB) *gorm.DB {
    return db.Select("id", "sector", "quantity", "get_price", "code", "usd_jpy", "user_id")
}

// NewUsStockRepository は DefaultStrategyRepository の新しいインスタンスを作成します
func NewUsStockRepository(db *gorm.DB) UsStockRepository {
    return &DefaultUsStockRepository{DB: db}
}

// 指定したuserIdのユーザーが保有する米国株式のリストを取得する
func (r *DefaultUsStockRepository) FetchUsStockListById(ctx context.Context, userId uint) ([]model.UsStock, error) {
    var usStocks []model.UsStock
    err := selectBaseQuery(r.DB).Where("user_id = ?", userId).Find(&usStocks).Error
    if err != nil {
        return nil, err
    }
    return usStocks, nil
}

// 米国株式情報を更新します
func (r *DefaultUsStockRepository) UpdateUsStock(ctx context.Context, dto UpdateUsStockDto) (*model.UsStock, error) {
    // 更新用のマップを作成します
    newStock := map[string]interface{}{}

    if dto.GetPrice != nil {
        newStock["get_price"] = dto.GetPrice
    }
    if dto.Quantity != nil {
        newStock["quantity"] = dto.Quantity
    }
    if dto.UsdJpy != nil {
        newStock["usd_jpy"] = dto.UsdJpy
    }

    // 指定されたIDの株式情報を更新します
    if err := r.DB.Model(&model.UsStock{}).Where("id = ?", dto.ID).Updates(newStock).Error; err != nil {
        return nil, err
    }

    // 更新された情報を取得します
    var usStock model.UsStock
    if err := selectBaseQuery(r.DB).Where("id = ?", dto.ID).Find(&usStock).Error; err != nil {
        return nil, err
    }

    return &usStock, nil
}

// 米国株式情報を作成します
func (r *DefaultUsStockRepository) CreateUsStock(ctx context.Context, dto CreateUsStockDto) (*model.UsStock, error) {
    // 既に同じ銘柄が存在するかを確認
    var existingUsStock model.UsStock
    if err := selectBaseQuery(r.DB).Where("code = ? AND user_id = ?", dto.Code, dto.UserId).First(&existingUsStock).Error; err == nil {
        return nil, fmt.Errorf("この銘柄は既に登録されています")
    }

    // 新しい米国株式情報を作成
    usStock := &model.UsStock{
        Code:   dto.Code,
        GetPrice: dto.GetPrice,
        Quantity: dto.Quantity,
        UserId:   dto.UserId,
        Sector:   dto.Sector,
        UsdJpy:   dto.UsdJpy,
    }

    if err := r.DB.Create(&usStock).Error; err != nil {
        return nil, err
    }

    return usStock, nil
}

// 米国株式情報を削除します
func (r *DefaultUsStockRepository) DeleteUsStock(ctx context.Context, id uint) error {
    // 指定されたIDの株式情報を検索して削除
    if err := r.DB.Where("id = ?", id).Delete(&model.UsStock{}).Error; err != nil {
        return err
    }
    return nil
}

