package totalassets

type CreateTotalAssetDto struct {
    CashJpy float64 `json:"cashJpy"`
	CashUsd float64 `json:"cashUsd"`
	Stock float64 `json:"stock"`
	Fund float64 `json:"fund"`
	Crypto float64 `json:"crypto"`
	FixedIncomeAsset float64 `json:"fixedIncomeAsset"`
	UserId   uint  `json:"userId"`
}