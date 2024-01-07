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
	FetchTotalAssetListById(ctx context.Context, userId uint, day int) ([]model.TotalAsset, error)
    FindTodayTotalAsset(ctx context.Context, userId uint) (*model.TotalAsset, error)
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
func (r *DefaultTotalAssetRepository) FetchTotalAssetListById(ctx context.Context, userId uint, day int) ([]model.TotalAsset, error) {
    var assets []model.TotalAsset

    // クエリビルダーの作成
    query := r.DB.Where("user_id = ?", userId).Order("created_at desc")

    // dayが0でなければLimitを設定
    if day != 0 {
        query = query.Limit(day)
    }

    // クエリの実行
    err := query.Find(&assets).Error
    if err != nil {
        return nil, err
    }

    return assets, nil
}


// 指定したuserIdのユーザーが保有する当日の資産総額を取得する
func (r *DefaultTotalAssetRepository) FindTodayTotalAsset(ctx context.Context, userId uint) (*model.TotalAsset, error) {
    // 現在の日付の始まりと終わりをUTCで取得
    nowUTC := time.Now().UTC()
    // ①ローカルタイムゾーンのオフセット(時)を取得
    _, offsetSec := time.Now().Zone()
    offsetHours := offsetSec / 3600

    // ②ローカルタイムゾーンでの現在の日付を計算
    localNow := nowUTC.Add(time.Duration(offsetHours) * time.Hour)
    // ③ローカルタイムゾーンでの現在の日付の範囲を取得
    //  ②で求めた日付での0:00~23:59を変数として作成
    todayStartLocal := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 0, 0, 0, 0, localNow.Location())
    todayEndLocal := todayStartLocal.Add(24 * time.Hour)
    // ④③で求めた２つの日付をUTCに戻して(JSTの場合、-9時間する)、それをtodayStart, todayEndとしてdbのselectの条件として用いる
    todayStartUTC := todayStartLocal.Add(time.Duration(-offsetHours) * time.Hour)
    todayEndUTC := todayEndLocal.Add(time.Duration(-offsetHours) * time.Hour)

    var asset model.TotalAsset
    err := r.DB.Where("user_id = ? AND created_at >= ? AND created_at < ?", userId, todayStartUTC, todayEndUTC).First(&asset).Error
    if err != nil {
        return nil, err
    }
    return &asset, nil
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
    if dto.Fund != nil {
        newAsset["fund"] = dto.Fund
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

// 当日の資産総額情報を作成します
func (r *DefaultTotalAssetRepository) CreateTodayTotalAsset(ctx context.Context, dto CreateTotalAssetDto) (*model.TotalAsset, error) {
    // FindTodayTotalAssetメソッドを使用して、同じ日付で同じユーザーのレコードが存在するか確認
    if existingAsset, err := r.FindTodayTotalAsset(ctx, dto.UserId); err == nil && existingAsset != nil {
        // レコードが存在する場合、エラーを返す
        return nil, errors.New("既に資産が登録されています。新規追加ではなく更新を行ってください。")
    } else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        // GORMのErrRecordNotFound以外のエラーが発生した場合
        return nil, err
    }

    newAsset := &model.TotalAsset{
        CashJpy:   dto.CashJpy,
        CashUsd:   dto.CashUsd,
        Stock:   dto.Stock,
        Fund: dto.Fund,
        Crypto:   dto.Crypto,
        FixedIncomeAsset: dto.FixedIncomeAsset,
        UserId:   dto.UserId,
    }

    if err := r.DB.Create(&newAsset).Error; err != nil {
        return nil, err
    }

    return newAsset, nil
}