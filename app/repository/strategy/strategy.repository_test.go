package strategy

import (
	"context"
	"my-us-stock-backend/app/database/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// テスト用のデータベース設定
func setupTestDB() *gorm.DB {
    db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    // テスト用のテーブルを準備
    db.AutoMigrate(&model.Strategy{})

    return db
}

func TestFindStrategy(t *testing.T) {
    db := setupTestDB()
    repo := NewStrategyRepository(db)

    // テスト用データを作成
    strategy := model.Strategy{Text: "Test Text", UserId: 99}
    db.Create(&strategy)

    // User IDで検索
    found, err := repo.FindStrategy(context.Background(), strategy.UserId)
    assert.NoError(t, err)
    assert.NotNil(t, found)
    assert.Equal(t, strategy.Text, found.Text)
    assert.Equal(t, strategy.UserId, found.UserId)
}

func TestUpdateStrategy(t *testing.T) {
    db := setupTestDB() // テスト用のデータベースセットアップ関数
    repo := NewStrategyRepository(db)

    // テスト用データを作成
    originalStrategy := model.Strategy{Text: "Original Text", UserId: 99}
    db.Create(&originalStrategy)

    // 更新用DTOの作成
    updateDto := UpdateStrategyDto{
        Text:   "Updated Text",
        ID:     originalStrategy.ID,
        UserId: originalStrategy.UserId,
    }

    // 戦略を更新
    updatedStrategy, err := repo.UpdateStrategy(context.Background(), updateDto)
    assert.NoError(t, err)
    assert.NotNil(t, updatedStrategy)
    assert.Equal(t, updateDto.Text, updatedStrategy.Text)
    assert.Equal(t, updateDto.UserId, updatedStrategy.UserId)

    // データベースから直接取得して検証
    var dbStrategy model.Strategy
    db.First(&dbStrategy, originalStrategy.ID)
    assert.Equal(t, updateDto.Text, dbStrategy.Text)
    assert.Equal(t, updateDto.UserId, dbStrategy.UserId)
}

func TestCreateStrategy(t *testing.T) {
    db := setupTestDB()
    repo := NewStrategyRepository(db)

    // 新しいデータを作成
	createDto := CreateStrategyDto{
        Text:   "New Strategy",
        UserId: 99,
    }
    created, err := repo.CreateStrategy(context.Background(), createDto)
    assert.NoError(t, err)
    assert.NotNil(t, created)
    assert.Equal(t, createDto.Text, created.Text)
    assert.Equal(t, createDto.UserId, created.UserId)

    // データベースでユーザーを確認
    var strategy model.Strategy
    db.First(&strategy, created.ID)
    assert.Equal(t, createDto.Text, created.Text)
    assert.Equal(t, createDto.UserId, created.UserId)
}
