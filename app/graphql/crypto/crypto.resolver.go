package crypto

import (
	"context"
	"my-us-stock-backend/app/graphql/generated"
)

type Resolver struct {
    CryptoService CryptoService
}

func NewResolver(CryptoService CryptoService) *Resolver {
    return &Resolver{CryptoService: CryptoService}
}

func (r *Resolver) Cryptos(ctx context.Context) ([]*generated.Crypto, error) {
    return r.CryptoService.Cryptos(ctx)
}

func (r *Resolver) CreateCrypto(ctx context.Context, input generated.CreateCryptoInput) (*generated.Crypto, error) {
    newCrypto, err := r.CryptoService.CreateCrypto(ctx, input)
    if err != nil {
        return nil, err
    }

    return newCrypto, nil
}