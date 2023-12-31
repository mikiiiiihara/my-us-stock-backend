package auth

import (
	"context"
	"errors"
	"fmt"
	"my-us-stock-backend/app/common/auth/logic"
	userModel "my-us-stock-backend/app/database/model"
	"my-us-stock-backend/app/graphql/utils"
	"my-us-stock-backend/app/repository/user"
	"net/http"
	"os"
	"time"

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
    FetchUserIdAccessToken(ctx context.Context) (uint, error)
}

// DefaultAuthService 構造体の定義
type DefaultAuthService struct {
	userRepository user.UserRepository
	userLogic logic.UserLogic
	responseLogic logic.ResponseLogic
	jwtLogic logic.JWTLogic
	authValidation AuthValidation
}

// NewAuthService は DefaultAuthService の新しいインスタンスを作成します
func NewAuthService(ur user.UserRepository, ul logic.UserLogic, rl logic.ResponseLogic, jl logic.JWTLogic, av AuthValidation) AuthService {
	return &DefaultAuthService{ur, ul, rl, jl, av}
}

// SignIn ログイン処理
func (as *DefaultAuthService) SignIn(ctx context.Context, c *gin.Context) (*userModel.User, error) {
    var signInRequestParam SignInRequest
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
	var signUpRequestParam SignUpRequest
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
	createDto := user.CreateUserDto{
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
    // 環境変数NODE_ENVを読み込む
    env := os.Getenv("ENV")

    // 開発環境かどうかを判定
    isDev := env == "development"
  secure := !isDev
  sameSiteValue := "Lax"
  if !isDev {
      sameSiteValue = "None"
  }

  // アクセストークンとリフレッシュトークンをクッキーとしてセット
  accessTokenCookie := &http.Cookie{
      Name:     "access_token",
      Value:    accessToken,
      Path:     "/",
      MaxAge:   3600,
      Secure:   secure,
      HttpOnly: true,
      SameSite: http.SameSiteLaxMode,
  }
  refreshTokenCookie := &http.Cookie{
      Name:     "refresh_token",
      Value:    refreshToken,
      Path:     "/",
      MaxAge:   3600*24*30,
      Secure:   secure,
      HttpOnly: true,
      SameSite: http.SameSiteLaxMode,
  }

  // Cookieをヘッダーに追加
  http.SetCookie(c.Writer, accessTokenCookie)
  http.SetCookie(c.Writer, refreshTokenCookie)

  // SameSite属性をヘッダーに追加
  c.Header("Set-Cookie", fmt.Sprintf("%s; SameSite=%s", accessTokenCookie.String(), sameSiteValue))
  c.Header("Set-Cookie", fmt.Sprintf("%s; SameSite=%s", refreshTokenCookie.String(), sameSiteValue))

        // UserModelをUserResponseに変換する
        userResponse := convertToUserResponse(user)

        // レスポンスにアクセストークンとリフレッシュトークンを含める
        response := AuthResponse{
            User:         userResponse,
        }
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
    // 環境変数NODE_ENVを読み込む
    env := os.Getenv("ENV")

    // 開発環境かどうかを判定
    isDev := env == "development"
    secure := !isDev
    sameSiteValue := "Lax"
    if !isDev {
        sameSiteValue = "None"
    }
    accessTokenCookie := &http.Cookie{
        Name:     "access_token",
        Value:    newAccessToken,
        Path:     "/",
        MaxAge:   3600,
        Secure:   secure,
        HttpOnly: true,
        SameSite: http.SameSiteLaxMode,
    }
      // Cookieをヘッダーに追加
    http.SetCookie(c.Writer, accessTokenCookie)

    // 新たなaccessTokenをcookieにセット
    c.Header("Set-Cookie", fmt.Sprintf("%s; SameSite=%s", accessTokenCookie.String(), sameSiteValue))
    return newAccessToken, nil
}


// convertToUserResponse はUserModelをUserResponseに変換します
func convertToUserResponse(user *userModel.User) UserResponse {
    return UserResponse{
        Name:  user.Name,
        Email: user.Email,
        // 必要に応じて他のフィールドもマッピング
    }
}

// validateAccessToken は与えられたアクセストークンが有効かどうかを検証します
func (as *DefaultAuthService) FetchUserIdAccessToken(ctx context.Context) (uint, error) {
    // cookieからアクセストークンを取得
    accessToken, _ := ctx.Value(utils.CookieKey).(string)
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
        return 0, utils.UnauthenticatedError("Error parsing token")
    }

    // トークンの有効性をチェック
    if !token.Valid {
        fmt.Println("Invalid token")
        return 0, utils.UnauthenticatedError("Invalid token")
    }

    // トークンのクレームを検証
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return 0, utils.UnauthenticatedError("Invalid token claims")
    }

    // 有効期限のチェック
    if exp, ok := claims["exp"].(float64); !ok || int64(exp) < time.Now().Unix() {
        return 0, utils.UnauthenticatedError("Token expired")
    }

    // ユーザーIDの取得とログ出力
    userId, ok := claims["id"].(float64)
    if !ok || userId == 0 {
        return 0, utils.UnauthenticatedError("Invalid user ID")
    }
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
