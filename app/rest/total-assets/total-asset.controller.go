package totalassets

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TotalAssetController struct {
	TotalAssetService TotalAssetService
}

func NewTotalAssetController(totalAssetService TotalAssetService) *TotalAssetController {
	return &TotalAssetController{TotalAssetService: totalAssetService}
}

// 資産新規登録
func (ac *TotalAssetController) CreateTodayTotalAsset(c *gin.Context) {
    ctx := c.Request.Context() // context.Context を取得
    response, err := ac.TotalAssetService.CreateTodayTotalAsset(ctx, c)
	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"result": response})
}