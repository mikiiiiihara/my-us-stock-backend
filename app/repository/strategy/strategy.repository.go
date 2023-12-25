package strategy

import (
	"context"
	"my-us-stock-backend/app/database/model"

	"gorm.io/gorm"
)

// StrategyRepository インターフェースの定義
type StrategyRepository interface {
    FindStrategy(ctx context.Context, userId uint) (*model.Strategy, error)
    UpdateStrategy(ctx context.Context, dto UpdateStrategyDto) (*model.Strategy, error)
    CreateStrategy(ctx context.Context, dto CreateStrategyDto) (*model.Strategy, error)
}

// DefaultStrategyRepository 構造体の定義
type DefaultStrategyRepository struct {
    DB *gorm.DB
}

// NewStrategyRepository は DefaultStrategyRepository の新しいインスタンスを作成します
func NewStrategyRepository(db *gorm.DB) StrategyRepository {
    return &DefaultStrategyRepository{DB: db}
}

func (r *DefaultStrategyRepository) FindStrategy(ctx context.Context, userId uint) (*model.Strategy, error) {
    strategy := new(model.Strategy)  // ポインタとしてオブジェクトを作成
    err := r.DB.Where("user_id = ?", userId).First(&strategy).Error
    if err != nil {
        return nil, err
    }
    return strategy, nil
}

// 戦略メモを更新します
func (r *DefaultStrategyRepository) UpdateStrategy(ctx context.Context, dto UpdateStrategyDto) (*model.Strategy, error) {
    strategy := &model.Strategy{Text: dto.Text, UserId: dto.UserId}
    if err := r.DB.Model(&strategy).Where("id = ?", dto.ID).Updates(strategy).Error; err != nil {
        return nil, err
    }
    return strategy, nil
}

// 戦略メモをデータベースに保存します
func (r *DefaultStrategyRepository) CreateStrategy(ctx context.Context, dto CreateStrategyDto) (*model.Strategy, error) {
    strategy := &model.Strategy{Text: dto.Text, UserId: dto.UserId}
    if err := r.DB.Create(strategy).Error; err != nil {
        return nil, err
    }
    return strategy, nil
}