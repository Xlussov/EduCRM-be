package domain

import "context"

type TxManager interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
