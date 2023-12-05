package currency

// CurrencyModule は CurrencyRepository の依存関係を管理します。
type CurrencyModule struct {
    Repository *CurrencyRepository
}

// NewCurrencyModule は新しい CurrencyModule を作成します。
func NewCurrencyModule() *CurrencyModule {
    // nil を渡すと、CurrencyRepository 内で http.DefaultClient が使用されます。
    repo := NewCurrencyRepository(nil)

    return &CurrencyModule{
        Repository: repo,
    }
}