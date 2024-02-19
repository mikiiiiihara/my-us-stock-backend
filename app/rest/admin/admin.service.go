package admin

import (
	"context"
	"my-us-stock-backend/app/database/model"
	"my-us-stock-backend/app/repository/market-price/fund"
)

type FundPriceService interface {
	FetchFundPrices(ctx context.Context) ([]model.FundPrice, error)
	UpdateFundPrice(ctx context.Context, dto fund.UpdateFundPriceDto) (*model.FundPrice, error)
	CreateFundPrice(ctx context.Context, dto fund.CreateFundPriceDto) (*model.FundPrice, error)
}

// DefaultFundPriceService provides a struct to hold any dependencies
type DefaultFundPriceService struct {
	Repo fund.FundPriceRepository
}

// NewFundPriceService creates a new instance of the fund price service
func NewFundPriceService(repo fund.FundPriceRepository) FundPriceService {
	return &DefaultFundPriceService{Repo: repo}
}

// Implementation of the service methods
func (s *DefaultFundPriceService) FetchFundPrices(ctx context.Context) ([]model.FundPrice, error) {
	return s.Repo.FetchFundPriceList(ctx)
}

func (s *DefaultFundPriceService) UpdateFundPrice(ctx context.Context, dto fund.UpdateFundPriceDto) (*model.FundPrice, error) {
	return s.Repo.UpdateFundPrice(ctx, dto)
}

func (s *DefaultFundPriceService) CreateFundPrice(ctx context.Context, dto fund.CreateFundPriceDto) (*model.FundPrice, error) {
	return s.Repo.CreateFundPrice(ctx, dto)
}
