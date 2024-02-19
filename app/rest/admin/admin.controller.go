package admin

import (
	"net/http"

	"my-us-stock-backend/app/repository/market-price/fund"

	"github.com/gin-gonic/gin"
)

// FundPriceController holds the service for dealing with fund prices
type FundPriceController struct {
	Service FundPriceService
}

// NewFundPriceController creates a new controller for fund prices
func NewFundPriceController(service FundPriceService) *FundPriceController {
	return &FundPriceController{
		Service: service,
	}
}

// GetFundPrices handles GET requests to fetch fund prices
func (fpc *FundPriceController) GetFundPrices(c *gin.Context) {
	fundPrices, err := fpc.Service.FetchFundPrices(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fundPrices)
}

// UpdateFundPrice handles POST requests to update a fund price
func (fpc *FundPriceController) UpdateFundPrice(c *gin.Context) {
	var dto fund.UpdateFundPriceDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedFundPrice, err := fpc.Service.UpdateFundPrice(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedFundPrice)
}

// CreateFundPrice handles POST requests to create a new fund price
func (fpc *FundPriceController) CreateFundPrice(c *gin.Context) {
	var dto fund.CreateFundPriceDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdFundPrice, err := fpc.Service.CreateFundPrice(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdFundPrice)
}
