package logic

import (
	"my-us-stock-backend/app/repository/user/model"
	"os"
	"strconv"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/form3tech-oss/jwt-go"
)

type JWTLogic interface {
	CreateAccessToken(user *model.User) (string, error)
	CreateRefreshToken(user *model.User) (string, error)
}

type jwtLogic struct{}

func NewJWTLogic() JWTLogic {
	return &jwtLogic{}
}

// CreateAccessToken jwtトークンの新規作成
func (jl *jwtLogic) CreateAccessToken(user *model.User) (string, error) {
	// headerのセット
	token := jwt.New(jwt.SigningMethodHS256)
	// claimsのセット
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = true
	claims["sub"] = strconv.Itoa(int(user.ID)) + user.Email + user.Name
	claims["id"] = user.ID
	claims["name"] = user.Name
	// latを取り除かないとミドルウェアで「Token used before issued」エラーになる
	// https://github.com/dgrijalva/jwt-go/issues/314#issuecomment-812775567
	// claims["iat"] = time.Now() // jwtの発行時間
	// 経過時間
	// 経過時間を過ぎたjetは処理しないようになる
	// ここでは15分の経過時間をリミットにしている
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	// 電子署名
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_KEY")))

	return tokenString, nil
}

func (jl *jwtLogic) CreateRefreshToken(user *model.User) (string, error) {
    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)

    // リフレッシュトークンにはユーザーIDのみ含める
    claims["id"] = user.ID

    // リフレッシュトークンの有効期限を7日間に設定
    claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

    // 電子署名
    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

// JwtMiddleware jwt認証のミドルウェア
var JwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})
