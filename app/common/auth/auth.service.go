package auth

import (
	"context"
	"errors"
	"fmt"
	"my-us-stock-backend/app/common/auth/logic"
	"my-us-stock-backend/app/common/auth/model"
	"my-us-stock-backend/app/common/auth/validation"
	userModel "my-us-stock-backend/app/database/model"
	"my-us-stock-backend/app/graphql/utils"
	"my-us-stock-backend/app/repository/user"
	"my-us-stock-backend/app/repository/user/dto"
	"net/http"
	"os"
	"time"

	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AuthService インターフェースの定義
type AuthService interface {
	SignIn(ctx context.Context, c *gin.Context) (*userModel.User, error)
	SignUp(ctx context.Context, c *gin.Context) (*userModel.User, error)
	SendAuthResponse(ctx context.Context, c *gin.Context, user *userModel.User, code int)
	RefreshAccessToken(c *gin.Context) (string, error)
    FetchUserIdAccessToken(accessToken string) (uint, error)
}

// DefaultAuthService 構造体の定義
type DefaultAuthService struct {
	userRepository user.UserRepository
	authLogic logic.AuthLogic
	userLogic logic.UserLogic
	responseLogic logic.ResponseLogic
	jwtLogic logic.JWTLogic
	authValidation validation.AuthValidation
}

// NewAuthService は DefaultAuthService の新しいインスタンスを作成します
func NewAuthService(ur user.UserRepository, al logic.AuthLogic, ul logic.UserLogic, rl logic.ResponseLogic, jl logic.JWTLogic, av validation.AuthValidation) AuthService {
	return &DefaultAuthService{ur, al, ul, rl, jl, av}
}

// SignIn ログイン処理
func (as *DefaultAuthService) SignIn(ctx context.Context, c *gin.Context) (*userModel.User, error) {
    var signInRequestParam model.SignInRequest
    if err := c.BindJSON(&signInRequestParam); err != nil {
        // JSONパースエラーが発生した場合、400 Bad Requestを返す
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return nil, err
    }

    if err := as.authValidation.SignInValidate(signInRequestParam); err != nil {
        // validationエラーの場合、400 Bad Requestを返す
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return nil, err
    }

    user, err := as.userRepository.GetUserByEmail(ctx, signInRequestParam.Email)
    if err != nil {
        // メールアドレス不一致の場合、400 Bad Requestを返す
        c.JSON(http.StatusBadRequest, gin.H{"error": "入力されたメールアドレスまたはパスワードが一致しません。"})
        return nil, err
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(signInRequestParam.Password)); err != nil {
        // パスワード不一致の場合、400 Bad Requestを返す
        c.JSON(http.StatusBadRequest, gin.H{"error": "入力されたメールアドレスまたはパスワードが一致しません。"})
        return nil, err
    }

    return user, nil
}


// SignUp 会員登録処理
func (as *DefaultAuthService) SignUp(ctx context.Context, c *gin.Context) (*userModel.User, error) {
	var signUpRequestParam model.SignUpRequest
    if err := c.BindJSON(&signUpRequestParam); err != nil {
        // JSONパースエラーが発生した場合、400 Bad Requestを返す
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return nil, err
    }
	if err := as.authValidation.SignUpValidate(signUpRequestParam); err != nil {
        // validationエラーの場合、400 Bad Requestを返す
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}

	users, err := as.userRepository.GetAllUserByEmail(ctx, signUpRequestParam.Email)
	if len(users) > 0 || err != nil {
        // すでに登録されているメールアドレスの場合、400 Bad Requestを返す
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(signUpRequestParam.Password), bcrypt.DefaultCost)
	if err != nil {
        // パスワード作成失敗の場合、400 Bad Requestを返す
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, err
	}
	createDto := dto.CreateUserDto{
		Name:     signUpRequestParam.Name,
		Email:    signUpRequestParam.Email,
		Password: string(hashPassword),
	}

	createdUser, err := as.userRepository.CreateUser(ctx, createDto)
	if err != nil {
        // ユーザー作成失敗の場合、500 Internal Server Errorを返す
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, err
	}

	return createdUser, nil
}


// SendAuthResponse レスポンス送信処理
func (as *DefaultAuthService) SendAuthResponse(ctx context.Context, c *gin.Context, user *userModel.User, code int) {
    // アクセストークンの生成
    accessToken, err := as.jwtLogic.CreateAccessToken(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create access token"})
        return
    }

    // リフレッシュトークンの生成
    refreshToken, err := as.jwtLogic.CreateRefreshToken(user)
    if err != nil {
		fmt.Println("Sending JSON response...4")
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create refresh token"})
        return
    }

    // UserModelをUserResponseに変換する
    userResponse := convertToUserResponse(user)

    // レスポンスにアクセストークンとリフレッシュトークンを含める
    response := model.AuthResponse{
        User:         userResponse,
    }
    // アクセストークンとリフレッシュトークンをクッキーとしてセット
    c.SetCookie("access_token", accessToken, 0, "/", "", false, true)
    c.SetCookie("refresh_token", refreshToken, 0, "/", "", false, true)
    c.JSON(code, response)
}

// RefreshAccessToken リフレッシュトークンを使用して新しいアクセストークンを生成
func (as *DefaultAuthService) RefreshAccessToken(c *gin.Context) (string, error) {
        // クッキーからrefresh_tokenを取得
        refreshToken, err := c.Cookie("refresh_token")
        if err != nil {
            return "", errors.New("access_token not found in cookie")
        }
    // refreshToken の検証ロジックを実装

    // ここで refreshToken の検証を行います
    // 例えば、データベース内のリフレッシュトークンと照合するなどの検証を行います

    userId,valid := validateRefreshTokenAndGetUserID(refreshToken)
    if !valid {
        // c.Status(http.StatusUnauthorized) // HTTPステータスコードを401に設定
        return "", errors.New("invalid refreshToken")
    }

    user, err := as.userRepository.FindUserByID(c, userId)
    if err != nil {
        return "", err
    }
    // refreshToken の検証が成功した場合、新しい accessToken を生成
    newAccessToken, err := as.jwtLogic.CreateAccessToken(user)
    if err != nil {
        return "", err
    }

    // 新たなaccessTokenをcookieにセット
    c.SetCookie("access_token", newAccessToken, 0, "/", "", false, true)
    return newAccessToken, nil
}


// convertToUserResponse はUserModelをUserResponseに変換します
func convertToUserResponse(user *userModel.User) model.UserResponse {
    return model.UserResponse{
        Name:  user.Name,
        Email: user.Email,
        // 必要に応じて他のフィールドもマッピング
    }
}

// validateAccessToken は与えられたアクセストークンが有効かどうかを検証します
func (as *DefaultAuthService) FetchUserIdAccessToken(accessToken string) (uint, error) {
    token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, &utils.GraphQLAuthError{Code: "UNAUTHENTICATED", Message: "Invalid signing method"}
        }
        // 署名キーを返す
        return []byte(os.Getenv("JWT_KEY")), nil
    })

    // トークン解析のエラーをチェック
    if err != nil {
        fmt.Println("Error parsing token:", err)
        return 0, &gqlerror.Error{
            Message: "Error parsing token",
            Extensions: map[string]interface{}{
                "code": "UNAUTHENTICATED",
            },
        }
    }

    // トークンの有効性をチェック
    if !token.Valid {
        fmt.Println("Invalid token")
        return 0, &gqlerror.Error{
            Message: "Invalid token",
            Extensions: map[string]interface{}{
                "code": "UNAUTHENTICATED",
            },
        }
    }

    // トークンのクレームを検証
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return 0, &gqlerror.Error{
            Message: "Invalid token claims",
            Extensions: map[string]interface{}{
                "code": "UNAUTHENTICATED",
            },
        }
    }

    // 有効期限のチェック
    if exp, ok := claims["exp"].(float64); !ok || int64(exp) < time.Now().Unix() {
        return 0, &gqlerror.Error{
            Message: "Token expired",
            Extensions: map[string]interface{}{
                "code": "UNAUTHENTICATED",
            },
        }
    }

    // ユーザーIDの取得とログ出力
    userId, ok := claims["id"].(float64)
    if !ok || userId == 0 {
        return 0, &gqlerror.Error{
            Message: "Invalid user ID",
            Extensions: map[string]interface{}{
                "code": "UNAUTHENTICATED",
            },
        }
    }
    fmt.Printf("Extracted UserID: %d, Type assertion successful: %t\n", uint(userId), ok)
    
    // すべての検証が成功した場合はユーザーIDを返す
    return uint(userId), nil
}

func validateRefreshTokenAndGetUserID(refreshToken string) (uint, bool) {
    // refreshToken の検証ロジックを実装
    token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("invalid refreshToken")
        }
        return []byte(os.Getenv("JWT_KEY")), nil
    })

    // エラーが発生した場合やトークンが無効な場合は false を返す
    if err != nil || !token.Valid {
        return 0, false
    }

    // 有効期限とユーザーIDを確認
    claims, claimOk := token.Claims.(jwt.MapClaims)
    if !claimOk {
        return 0, false
    }

    exp, expOk := claims["exp"].(float64)
    if !expOk || int64(exp) < time.Now().Unix() {
        return 0, false
    }

    // ユーザーIDの取得
    userId, userIdOk := claims["id"].(float64)
    if !userIdOk {
        return 0, false
    }

    // すべての検証が成功した場合はユーザーIDと true を返す
    return uint(userId), true
}
