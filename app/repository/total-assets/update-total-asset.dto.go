package totalassets

type UpdateTotalAssetDto struct {
    ID       uint     `json:"id"`
    CashJpy *float64 `json:"cashJpy"`
	CashUsd *float64 `json:"cashUsd"`
	Stock *float64 `json:"stock"`
	Crypto *float64 `json:"crypto"`
	FixedIncomeAsset *float64 `json:"fixedIncomeAsset"`
}
