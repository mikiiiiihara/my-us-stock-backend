package marketprice

import (
	"my-us-stock-backend/repository/market-price/dto"
	"my-us-stock-backend/repository/market-price/entity"
	"sort"
	"time"
)

// createDividendEntity は DividendEntity レスポンスを構築します。
func (repo *DefaultMarketPriceRepository) createDividendEntity(res *dto.DividendResponse) *entity.DividendEntity {
    dividends := filterDividends(res.Historical)
    totalCashAmount := 0.0
    for _, dividend := range dividends {
        totalCashAmount += dividend.Dividend
    }

    // 直近１年の配当総額を計算
    dividendTotal := 0.0
    if len(dividends) != 0 {
        dividendTotal = roundToThreeDecimals(totalCashAmount)
    }

    // 平均配当額を計算
    cashAmount := 0.0
    if len(dividends) != 0 {
        cashAmount = roundToThreeDecimals(totalCashAmount / float64(len(dividends)))
    }

    return &entity.DividendEntity{
        Ticker:           res.Symbol,
        DividendTime:     len(dividends),
        DividendMonth:    calculateDividendMonth(true, dividends),
        DividendFixedMonth: calculateDividendMonth(false, dividends),
        Dividend:         cashAmount,
        DividendTotal:    dividendTotal,
    }
}
// filterDividends は直近1年の配当記録をフィルタリングします。
func filterDividends(dividends []dto.Historical) []dto.Historical {
    oneYearAgo := time.Now().AddDate(-1, 0, 0)
    filteredDividends := make([]dto.Historical, 0)
    for _, dividend := range dividends {
        payDate, _ := time.Parse("2006-01-02", dividend.PaymentDate)
        if payDate.After(oneYearAgo) {
            filteredDividends = append(filteredDividends, dividend)
        }
    }
    return filteredDividends
}

// calculateDividendMonth は配当権利落月・支払い月を取得します。
func calculateDividendMonth(isPayment bool, dividends []dto.Historical) []int {
    monthSet := make(map[int]struct{})
    for _, dividend := range dividends {
        var month int
        if isPayment {
            month = parseMonth(dividend.PaymentDate)
        } else {
            month = parseMonth(dividend.Date)
        }
        monthSet[month] = struct{}{}
    }

    var months []int
    for month := range monthSet {
        months = append(months, month)
    }

    sort.Ints(months)
    return months
}