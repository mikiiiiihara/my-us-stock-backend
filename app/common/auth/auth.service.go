package auth

import (
	"context"
	"my-us-stock-backend/app/common/auth/logic"
	"my-us-stock-backend/app/common/auth/model"
	"my-us-stock-backend/app/common/auth/validation"
	"my-us-stock-backend/app/repository/user"
	"my-us-stock-backend/app/repository/user/dto"
	userModel "my-us-stock-backend/app/repository/user/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AuthService インターフェースの定義
type AuthService interface {
	GetUserIdFromToken(w http.ResponseWriter, r *http.Request) (int, error)
	SignIn(ctx context.Context, c *gin.Context) (*userModel.User, error)
	SignUp(ctx context.Context, c *gin.Context) (*userModel.User, error)
	SendAuthResponse(ctx context.Context, c *gin.Context, user *userModel.User, code int)
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

// GetUserIdFromToken tokenよりuserIdを取得
func (as *DefaultAuthService) GetUserIdFromToken(w http.ResponseWriter, r *http.Request) (int, error) {
	// トークンからuserIdを取得
	userId, err := as.authLogic.GetUserIdFromContext(r)
	if err != nil {
		errMessage := "認証エラー"
		as.responseLogic.SendResponse(w, as.responseLogic.CreateErrorStringResponse(errMessage), http.StatusUnauthorized)
		return 0, err
	}

	return userId, nil
}

// SignIn ログイン処理
func (as *DefaultAuthService) SignIn(ctx context.Context, c *gin.Context) (*userModel.User, error) {
	var signInRequestParam model.SignInRequest
	if err := c.BindJSON(&signInRequestParam); err != nil {
		return nil, err
	}

	if err := as.authValidation.SignInValidate(signInRequestParam); err != nil {
		return nil, err
	}

	user, err := as.userRepository.GetUserByEmail(ctx, signInRequestParam.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(signInRequestParam.Password)); err != nil {
		return nil, err
	}

	return user, nil
}

// SignUp 会員登録処理
func (as *DefaultAuthService) SignUp(ctx context.Context, c *gin.Context) (*userModel.User, error) {
	var signUpRequestParam model.SignUpRequest
	if err := c.BindJSON(&signUpRequestParam); err != nil {
		return nil, err
	}

	if err := as.authValidation.SignUpValidate(signUpRequestParam); err != nil {
		return nil, err
	}

	users, err := as.userRepository.GetAllUserByEmail(ctx, signUpRequestParam.Email)
	if len(users) > 0 || err != nil {
		return nil, err
	}

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(signUpRequestParam.Password), bcrypt.DefaultCost)
    createDto := dto.CreateUserDto{
		Name:     signUpRequestParam.Name,
		Email:    signUpRequestParam.Email,
		Password: string(hashPassword),
	}

	createdUser, err := as.userRepository.CreateUser(ctx, createDto)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

// SendAuthResponse レスポンス送信処理
func (as *DefaultAuthService) SendAuthResponse(ctx context.Context, c *gin.Context, user *userModel.User, code int) {
	token, err := as.jwtLogic.CreateJwtToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	// UserModelをUserResponseに変換する
	userResponse := convertToUserResponse(user)

	response := model.AuthResponse{
		Token: token,
		User:  userResponse,
	}

	c.JSON(code, response)
}

// convertToUserResponse はUserModelをUserResponseに変換します
func convertToUserResponse(user *userModel.User) model.UserResponse {
	return model.UserResponse{
		Name:  user.Name,
		Email: user.Email,
		// 必要に応じて他のフィールドもマッピング
	}
}
