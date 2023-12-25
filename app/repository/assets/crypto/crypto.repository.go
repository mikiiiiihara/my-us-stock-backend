package crypto

import (
	"context"
	"fmt"
	"my-us-stock-backend/app/database/model"

	"gorm.io/gorm"
)

// CryptoRepository インターフェースの定義
type CryptoRepository interface {
	FetchCryptoListById(ctx context.Context, userId uint) ([]model.Crypto, error)
    UpdateCrypto(ctx context.Context, dto UpdateCryptoDto) (*model.Crypto, error)
	CreateCrypto(ctx context.Context, dto CreateCryptDto) (*model.Crypto, error)
	DeleteCrypto(ctx context.Context, id uint) error
}

// DefaultCryptoRepository 構造体の定義
type DefaultCryptoRepository struct {
    DB *gorm.DB
}

// NewCryptoRepository は DefaultStrategyRepository の新しいインスタンスを作成します
func NewCryptoRepository(db *gorm.DB) CryptoRepository {
    return &DefaultCryptoRepository{DB: db}
}

// 指定したuserIdのユーザーが保有する米国株式のリストを取得する
func (r *DefaultCryptoRepository) FetchCryptoListById(ctx context.Context, userId uint) ([]model.Crypto, error) {
    var usStocks []model.Crypto
    err := r.DB.Where("user_id = ?", userId).Find(&usStocks).Error
    if err != nil {
        return nil, err
    }
    return usStocks, nil
}

// 米国株式情報を更新します
func (r *DefaultCryptoRepository) UpdateCrypto(ctx context.Context, dto UpdateCryptoDto) (*model.Crypto, error) {
    // 更新用のマップを作成します
    newStock := map[string]interface{}{}

    if dto.GetPrice != nil {
        newStock["get_price"] = dto.GetPrice
    }
    if dto.Quantity != nil {
        newStock["quantity"] = dto.Quantity
    }

    // 指定されたIDの株式情報を更新します
    if err := r.DB.Model(&model.Crypto{}).Where("id = ?", dto.ID).Updates(newStock).Error; err != nil {
        return nil, err
    }

    // 更新された情報を取得します
    var crypto model.Crypto
    if err := r.DB.Where("id = ?", dto.ID).Find(&crypto).Error; err != nil {
        return nil, err
    }

    return &crypto, nil
}

// 米国株式情報を作成します
func (r *DefaultCryptoRepository) CreateCrypto(ctx context.Context, dto CreateCryptDto) (*model.Crypto, error) {
    // 既に同じ銘柄が存在するかを確認
    var existingUsStock model.Crypto
    if err := r.DB.Where("code = ?", dto.Code).First(&existingUsStock).Error; err == nil {
        return nil, fmt.Errorf("この銘柄は既に登録されています")
    }

    // 新しい米国株式情報を作成
    crypto := &model.Crypto{
        Code:   dto.Code,
        GetPrice: dto.GetPrice,
        Quantity: dto.Quantity,
        UserId:   dto.UserId,
    }

    if err := r.DB.Create(&crypto).Error; err != nil {
        return nil, err
    }

    return crypto, nil
}

// 米国株式情報を削除します
func (r *DefaultCryptoRepository) DeleteCrypto(ctx context.Context, id uint) error {
    // 指定されたIDの株式情報を検索して削除
    if err := r.DB.Where("id = ?", id).Delete(&model.Crypto{}).Error; err != nil {
        return err
    }
    return nil
}

