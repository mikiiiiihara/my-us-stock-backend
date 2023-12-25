package stock

// import (
// 	"context"
// 	"my-us-stock-backend/app/common/auth"
// 	"my-us-stock-backend/app/database/model"
// 	"my-us-stock-backend/app/graphql/generated"
// 	"my-us-stock-backend/app/graphql/utils"
// 	"my-us-stock-backend/app/repository/assets/stock"
// 	"strconv"
// )

// // UsStockService インターフェースの定義
// type UsStockService interface {
//     GetUserByID(ctx context.Context) (*generated.User, error)
// }

// // DefaultUsStockService 構造体の定義
// type DefaultUsStockService struct {
//     Repo stock.UsStockRepository // インターフェースを利用
//     Auth auth.AuthService    // 認証サービスのインターフェース
// }

// // NewUsStockService は DefaultUsStockService の新しいインスタンスを作成します
// func NewUsStockService(repo stock.UsStockRepository, auth auth.AuthService) UsStockService {
//     return &DefaultUsStockService{Repo: repo, Auth: auth}
// }

// // GetUserByID はユーザーをIDによって検索します
// func (s *DefaultUsStockService) UsStocks(ctx context.Context) (*generated.User, error) {
//     // アクセストークンの検証
//     userId, _ := s.Auth.FetchUserIdAccessToken(ctx)
//     if userId == 0 {
//         return nil, utils.UnauthenticatedError("Invalid user ID")
//     }
//     modelUser, err := s.Repo.FetchUsStockListById(ctx, userId)
//     if err != nil {
//         return nil, err
//     }
//     return convertModelUserToGeneratedUser(modelUser), nil
// }

// // convertModelUserToGeneratedUser は model.User を generated.User に変換します
// func convertModelUserToGeneratedUser(modelUser *model.User) *generated.User {
//     if modelUser == nil {
//         return nil
//     }
//     return &generated.User{
//         ID:    strconv.FormatUint(uint64(modelUser.ID), 10),
//         Name:  modelUser.Name,
//         Email: modelUser.Email,
//     }
// }
