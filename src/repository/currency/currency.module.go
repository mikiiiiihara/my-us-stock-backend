package currency

import (
	"net/http"
	"os"
)

type CurrencyModule struct {
    Repository *CurrencyRepository
}

func NewCurrencyModule() *CurrencyModule {
    // 環境変数から URL を読み込む
    currencyURL := os.Getenv("CURRENCY_URL")
    if currencyURL == "" {
        // URLが設定されていない場合は、モジュールの作成を中止するか、デフォルト値を使用する
        panic("CURRENCY_URL is not set")
    }

    // nil を渡すと、CurrencyRepository 内で http.DefaultClient が使用されます
    repo := NewCurrencyRepository(http.DefaultClient, currencyURL)

    return &CurrencyModule{
        Repository: repo,
    }
}
