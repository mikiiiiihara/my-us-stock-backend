package fund

import (
	"context"
	"fmt"
	"my-us-stock-backend/app/database/model"

	"gorm.io/gorm"
)

// FundPriceRepository インターフェースの定義
type FundPriceRepository interface {
	FetchFundPriceList(ctx context.Context) ([]model.FundPrice, error)
    UpdateFundPrice(ctx context.Context, dto UpdateFundPriceDto) (*model.FundPrice, error)
	CreateFundPrice(ctx context.Context, dto CreateFundPriceDto) (*model.FundPrice, error)
}

// DefaultFundPriceRepository 構造体の定義
type DefaultFundPriceRepository struct {
    DB *gorm.DB
}

// 共通フィールドを選択するためのヘルパー関数です。
func selectBaseQuery(db *gorm.DB) *gorm.DB {
    return db.Select("id", "price","code", "name")
}

// NewFetchFundRepository は DefaultStrategyRepository の新しいインスタンスを作成します
func NewFetchFundRepository(db *gorm.DB) FundPriceRepository {
    return &DefaultFundPriceRepository{DB: db}
}

// 指定したuserIdのユーザーが保有する日本投資信託のリストを取得する
func (r *DefaultFundPriceRepository) FetchFundPriceList(ctx context.Context) ([]model.FundPrice, error) {
    var funds []model.FundPrice
    err := selectBaseQuery(r.DB).Find(&funds).Error
    if err != nil {
        return nil, err
    }
    return funds, nil
}

// 日本投資信託情報を更新します
func (r *DefaultFundPriceRepository) UpdateFundPrice(ctx context.Context, dto UpdateFundPriceDto) (*model.FundPrice, error) {
    // 更新用のマップを作成します
    newFund := map[string]interface{}{}

	newFund["price"] = dto.Price

    // 指定されたIDの株式情報を更新します
    if err := r.DB.Model(&model.FundPrice{}).Where("id = ?", dto.ID).Updates(newFund).Error; err != nil {
        return nil, err
    }

    // 更新された情報を取得します
    var FundPrice model.FundPrice
    if err := selectBaseQuery(r.DB).Where("id = ?", dto.ID).Find(&FundPrice).Error; err != nil {
        return nil, err
    }

    return &FundPrice, nil
}

// 日本投資信託情報を作成します
func (r *DefaultFundPriceRepository) CreateFundPrice(ctx context.Context, dto CreateFundPriceDto) (*model.FundPrice, error) {
    // 既に同じ銘柄が存在するかを確認
    var existingJapanFund model.FundPrice
    if err := selectBaseQuery(r.DB).Where("code = ?", dto.Code).First(&existingJapanFund).Error; err == nil {
        return nil, fmt.Errorf("この銘柄は既に登録されています")
    }

    // 新しい米国株式情報を作成
    fundPrice := &model.FundPrice{
        Name: dto.Name,
        Code:   dto.Code,
        Price: dto.Price,
    }

    if err := r.DB.Create(&fundPrice).Error; err != nil {
        return nil, err
    }

    return fundPrice, nil
}